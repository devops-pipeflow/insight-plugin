package review

import (
	"bufio"
	"bytes"
	"context"
	"encoding/base64"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	"github.com/reviewdog/reviewdog/diff"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	TypeError = "Error"
	TypeInfo  = "Info"
	TypeWarn  = "Warn"
)

const (
	queryLimit   = 1000
	urlChanges   = "/changes/"
	urlContent   = "/content"
	urlCurrent   = "current"
	urlDetail    = "/detail"
	urlDiff      = "/diff"
	urlFiles     = "/files/"
	urlNumber    = "&n="
	urlOption    = "&o="
	urlPatch     = "/patch"
	urlPrefix    = "/a"
	urlQuery     = "?q="
	urlReview    = "/review"
	urlRevisions = "/revisions/"
	urlStart     = "&start="
)

const (
	base64Content = ".base64"
	base64Message = "message.base64"
	commitMsg     = "/COMMIT_MSG"
	commitQuery   = "commit"
)

const (
	diffBin    = "Binary files differ"
	diffSep    = "diff --git"
	pathPrefix = "b/"
)

const (
	voteApproval    = "+1"
	voteDisapproval = "-1"
	voteLabel       = "Code-Review"
	voteMessage     = "Voting Code-Review by pipeflow insight"
)

var (
	queryOptions = []string{
		"CURRENT_FILES",
		"CURRENT_REVISION",
		"DETAILED_ACCOUNTS",
	}
)

type Review interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Clean(context.Context, string) error
	Diff(context.Context, int, string) (map[string]interface{}, error)
	Fetch(context.Context, string, string) (string, string, []string, error)
	Query(context.Context, string, int) ([]interface{}, error)
	Vote(context.Context, string, []Format) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type Format struct {
	File    string
	Line    int
	Type    string
	Details string
}

type review struct {
	cfg  *Config
	user string
	pass string
	url  string
}

func New(_ context.Context, cfg *Config) Review {
	return &review{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *review) Init(_ context.Context) error {
	r.cfg.Logger.Debug("review: Init")

	r.user = r.cfg.Config.Spec.ReviewConfig.User
	r.pass = r.cfg.Config.Spec.ReviewConfig.Pass
	r.url = r.cfg.Config.Spec.ReviewConfig.Url

	return nil
}

func (r *review) Deinit(_ context.Context) error {
	r.cfg.Logger.Debug("review: Deinit")

	return nil
}

func (r *review) Clean(_ context.Context, name string) error {
	r.cfg.Logger.Debug("review: Clean")

	if err := os.RemoveAll(name); err != nil {
		return errors.Wrap(err, "failed to clean")
	}

	return nil
}

func (r *review) Diff(ctx context.Context, change int, file string) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("review: Diff")

	buf, err := r.get(ctx, r.urlDiff(change, file))
	if err != nil {
		return nil, nil
	}

	ret, err := r.unmarshal(buf)
	if err != nil {
		return nil, nil
	}

	return ret, nil
}

// nolint:funlen,gocyclo
func (r *review) Fetch(ctx context.Context, root, commit string) (path, name string, files []string, err error) {
	r.cfg.Logger.Debug("review: Fetch")

	filterFiles := func(data map[string]interface{}) map[string]interface{} {
		buf := make(map[string]interface{})
		for key, val := range data {
			if v, ok := val.(map[string]interface{})["status"]; ok {
				if v.(string) == "D" || v.(string) == "R" {
					continue
				}
			}
			buf[key] = val
		}
		return buf
	}

	// Query commit
	buf, err := r.get(ctx, r.urlQuery(commitQuery+":"+commit, []string{"CURRENT_REVISION"}, 0))
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query")
	}

	queryRet, err := r.unmarshalList(buf)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to unmarshalList")
	}

	revisions := queryRet[0].(map[string]interface{})["revisions"].(map[string]interface{})
	current := revisions[queryRet[0].(map[string]interface{})["current_revision"].(string)].(map[string]interface{})

	changeNum := int(queryRet[0].(map[string]interface{})["_number"].(float64))
	revisionNum := int(current["_number"].(float64))

	path = filepath.Join(root, strconv.Itoa(changeNum), queryRet[0].(map[string]interface{})["current_revision"].(string))

	// Get files
	buf, err = r.get(ctx, r.urlFiles(changeNum, revisionNum))
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to files")
	}

	fs, err := r.unmarshal(buf)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to unmarshal")
	}

	// Match files
	fs = filterFiles(fs)

	// Get content
	for key := range fs {
		buf, err = r.get(ctx, r.urlContent(changeNum, revisionNum, key))
		if err != nil {
			return "", "", nil, errors.Wrap(err, "failed to content")
		}

		file := filepath.Base(key) + base64Content
		if key == commitMsg {
			file = base64Message
		}

		err = r.write(filepath.Join(path, filepath.Dir(key)), file, string(buf))
		if err != nil {
			return "", "", nil, errors.Wrap(err, "failed to fetch")
		}
	}

	// Return files
	for key := range fs {
		if key == commitMsg {
			files = append(files, base64Message)
		} else {
			files = append(files, filepath.Join(filepath.Dir(key), filepath.Base(key)+base64Content))
		}
	}

	return path, queryRet[0].(map[string]interface{})["project"].(string), files, nil
}

func (r *review) Query(ctx context.Context, search string, start int) ([]interface{}, error) {
	helper := func(search string, start int) []interface{} {
		buf, err := r.get(ctx, r.urlQuery(search, queryOptions, start))
		if err != nil {
			return nil
		}
		ret, err := r.unmarshalList(buf)
		if err != nil {
			return nil
		}
		return ret
	}

	buf := helper(search, start)
	if len(buf) == 0 {
		return []interface{}{}, nil
	}

	more, ok := buf[len(buf)-1].(map[string]interface{})["_more_changes"].(bool)
	if !ok {
		more = false
	}

	if !more {
		return buf, nil
	}

	if b, err := r.Query(ctx, search, start+len(buf)); err == nil {
		buf = append(buf, b...)
	}

	return buf, nil
}

// nolint:funlen,gocyclo
func (r *review) Vote(ctx context.Context, commit string, data []Format) error {
	match := func(data Format, diffs []*diff.FileDiff) bool {
		for _, d := range diffs {
			if strings.Replace(d.PathNew, pathPrefix, "", 1) != data.File {
				continue
			}
			if data.Line <= 0 {
				return true
			}
			for _, h := range d.Hunks {
				for _, l := range h.Lines {
					if l.Type == diff.LineAdded && l.LnumNew == data.Line {
						return true
					}
				}
			}
		}
		return false
	}

	build := func(data []Format, diffs []*diff.FileDiff) (map[string]interface{}, map[string]interface{}, string) {
		if len(data) == 0 {
			return nil, map[string]interface{}{voteLabel: voteApproval}, voteMessage
		}
		c := map[string]interface{}{}
		for _, item := range data {
			if item.Details == "" || (item.File != commitMsg && !match(item, diffs)) {
				continue
			}
			l := item.Line
			if l <= 0 {
				l = 1
			}
			b := map[string]interface{}{"line": l, "message": item.Details}
			if _, ok := c[item.File]; !ok {
				c[item.File] = []map[string]interface{}{b}
			} else {
				c[item.File] = append(c[item.File].([]map[string]interface{}), b)
			}
		}
		if len(c) == 0 {
			return nil, map[string]interface{}{voteLabel: voteApproval}, voteMessage
		} else {
			return c, map[string]interface{}{voteLabel: voteDisapproval}, voteMessage
		}
	}

	// Query commit
	ret, err := r.get(ctx, r.urlQuery(commitQuery+":"+commit, []string{"CURRENT_REVISION"}, 0))
	if err != nil {
		return errors.Wrap(err, "failed to query")
	}

	c, err := r.unmarshalList(ret)
	if err != nil {
		return errors.Wrap(err, "failed to unmarshalList")
	}

	revisions := c[0].(map[string]interface{})["revisions"].(map[string]interface{})
	current := revisions[c[0].(map[string]interface{})["current_revision"].(string)].(map[string]interface{})

	// Get patch
	ret, err = r.get(ctx, r.urlPatch(int(c[0].(map[string]interface{})["_number"].(float64)), int(current["_number"].(float64))))
	if err != nil {
		return errors.Wrap(err, "failed to patch")
	}

	// Parse diff
	dec := make([]byte, base64.StdEncoding.DecodedLen(len(ret)))
	if _, err = base64.StdEncoding.Decode(dec, ret); err != nil {
		return errors.Wrap(err, "failed to decode")
	}

	index := bytes.Index(dec, []byte(diffSep))
	if index < 0 {
		return errors.New("failed to index")
	}

	var b []byte

	for _, item := range bytes.SplitAfter(dec[index:], []byte(diffSep)) {
		if !bytes.Contains(item, []byte(diffBin)) {
			b = bytes.Join([][]byte{b, item}, []byte(""))
		}
	}

	diffs, err := diff.ParseMultiFile(bytes.NewReader(b))
	if err != nil {
		return errors.Wrap(err, "failed to parse")
	}

	// Review commit
	comments, labels, message := build(data, diffs)
	buf := map[string]interface{}{"comments": comments, "labels": labels, "message": message}
	if err := r.post(ctx, r.urlReview(int(c[0].(map[string]interface{})["_number"].(float64)),
		int(current["_number"].(float64))), buf); err != nil {
		return errors.Wrap(err, "failed to review")
	}

	return nil
}

func (r *review) write(dir, file, data string) error {
	r.cfg.Logger.Debug("review: write")

	_ = os.MkdirAll(dir, os.ModePerm)

	f, err := os.Create(filepath.Join(dir, file))
	if err != nil {
		return errors.Wrap(err, "failed to create")
	}
	defer func() { _ = f.Close() }()

	w := bufio.NewWriter(f)
	if _, err := w.WriteString(data); err != nil {
		return errors.Wrap(err, "failed to write")
	}
	defer func() { _ = w.Flush() }()

	return nil
}

func (r *review) unmarshal(data []byte) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("review: unmarshal")

	buf := map[string]interface{}{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	return buf, nil
}

func (r *review) unmarshalList(data []byte) ([]interface{}, error) {
	r.cfg.Logger.Debug("review: unmarshalList")

	var buf []interface{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf, nil
}

func (r *review) urlContent(change, revision int, name string) string {
	r.cfg.Logger.Debug("review: urlContent")

	buf := r.url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles + url.QueryEscape(name) + urlContent

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlFiles + url.QueryEscape(name) + urlContent
	}

	return buf
}

func (r *review) urlDetail(change int) string {
	r.cfg.Logger.Debug("review: urlDetail")

	buf := r.url + urlChanges + strconv.Itoa(change) + urlDetail

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) + urlDetail
	}

	return buf
}

func (r *review) urlDiff(change int, file string) string {
	r.cfg.Logger.Debug("review: urlDiff")

	buf := r.url + urlChanges + strconv.Itoa(change) + urlRevisions + urlCurrent + urlFiles + url.PathEscape(file) + urlDiff

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) + urlRevisions + urlCurrent + urlFiles + url.PathEscape(file) + urlDiff
	}

	return buf
}

func (r *review) urlFiles(change, revision int) string {
	r.cfg.Logger.Debug("review: urlFiles")

	buf := r.url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlFiles
	}

	return buf
}

func (r *review) urlPatch(change, revision int) string {
	r.cfg.Logger.Debug("review: urlPatch")

	buf := r.url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlPatch

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlPatch
	}

	return buf
}

func (r *review) urlQuery(search string, option []string, start int) string {
	r.cfg.Logger.Debug("review: urlQuery")

	query := urlQuery + url.PathEscape(search) +
		urlOption + strings.Join(option, urlOption) +
		urlStart + strconv.Itoa(start) +
		urlNumber + strconv.Itoa(queryLimit)

	buf := r.url + urlChanges + query
	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + query
	}

	return buf
}

func (r *review) urlReview(change, revision int) string {
	r.cfg.Logger.Debug("review: urlReview")

	buf := r.url + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlReview

	if r.user != "" && r.pass != "" {
		buf = r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
			urlRevisions + strconv.Itoa(revision) + urlReview
	}

	return buf
}

func (r *review) get(_ context.Context, _url string) ([]byte, error) {
	r.cfg.Logger.Debug("review: get")

	req, err := http.NewRequest(http.MethodGet, _url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "failed to request")
	}

	if r.user != "" && r.pass != "" {
		req.SetBasicAuth(r.user, r.pass)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return nil, errors.New("invalid status")
	}

	data, err := io.ReadAll(rsp.Body)
	if err != nil {
		return nil, errors.Wrap(err, "failed to read")
	}

	return data, nil
}

func (r *review) post(_ context.Context, _url string, data map[string]interface{}) error {
	r.cfg.Logger.Debug("review: post")

	buf, err := json.Marshal(data)
	if err != nil {
		return errors.Wrap(err, "failed to marshal")
	}

	req, err := http.NewRequest(http.MethodPost, _url, bytes.NewBuffer(buf))
	if err != nil {
		return errors.Wrap(err, "failed to request")
	}

	req.Header.Set("Content-Type", "application/json;charset=utf-8")

	if r.user != "" && r.pass != "" {
		req.SetBasicAuth(r.user, r.pass)
	}

	rsp, err := http.DefaultClient.Do(req)
	if err != nil {
		return errors.Wrap(err, "failed to do")
	}

	defer func() {
		_ = rsp.Body.Close()
	}()

	if rsp.StatusCode != http.StatusOK {
		return errors.New("invalid status")
	}

	_, err = io.ReadAll(rsp.Body)
	if err != nil {
		return errors.Wrap(err, "failed to read")
	}

	return nil
}
