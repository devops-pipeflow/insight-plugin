package review

import (
	"bufio"
	"bytes"
	"context"
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

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	base64Content = ".base64"
	base64Message = "message.base64"
	commitMsg     = "/COMMIT_MSG"
)

type Review interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Clean(string) error
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
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

func (r *review) Init(ctx context.Context) error {
	r.cfg.Logger.Debug("review: Init")

	r.user = r.cfg.Config.Spec.Review.User
	r.pass = r.cfg.Config.Spec.Review.Pass
	r.url = r.cfg.Config.Spec.Review.Url

	return nil
}

func (r *review) Deinit(ctx context.Context) error {
	r.cfg.Logger.Debug("review: Deinit")

	return nil
}

func (r *review) Clean(name string) error {
	r.cfg.Logger.Debug("review: Clean")

	if err := os.RemoveAll(name); err != nil {
		return errors.Wrap(err, "failed to clean")
	}

	return nil
}

// nolint:funlen,gocyclo
func (r *review) Fetch(ctx context.Context, root, commit string) (dname, rname string, flist []string, emsg error) {
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
	buf, err := r.get(ctx, r.urlQuery("commit:"+commit, []string{"CURRENT_REVISION"}, 0))
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to query")
	}

	queryRet, err := r.unmarshalList(buf)
	if err != nil {
		return "", "", nil, errors.Wrap(err, "failed to unmarshalList")
	}

	revisions := queryRet["revisions"].(map[string]interface{})
	current := revisions[queryRet["current_revision"].(string)].(map[string]interface{})

	changeNum := int(queryRet["_number"].(float64))
	revisionNum := int(current["_number"].(float64))

	path := filepath.Join(root, strconv.Itoa(changeNum), queryRet["current_revision"].(string))

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
	var files []string

	for key := range fs {
		if key == commitMsg {
			files = append(files, base64Message)
		} else {
			files = append(files, filepath.Join(filepath.Dir(key), filepath.Base(key)+base64Content))
		}
	}

	return path, queryRet["project"].(string), files, nil
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

func (r *review) unmarshalList(data []byte) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("review: unmarshalList")

	var buf []map[string]interface{}

	if err := json.Unmarshal(data[4:], &buf); err != nil {
		return nil, errors.Wrap(err, "failed to unmarshal")
	}

	if len(buf) == 0 {
		return nil, errors.New("failed to match")
	}

	return buf[0], nil
}

func (r *review) urlContent(change, revision int, name string) string {
	r.cfg.Logger.Debug("review: urlContent")

	buf := r.url + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/files/" + url.QueryEscape(name) + "/content"

	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/files/" + url.QueryEscape(name) + "/content"
	}

	return buf
}

func (r *review) urlDetail(change int) string {
	r.cfg.Logger.Debug("review: urlDetail")

	buf := r.url + "/changes/" + strconv.Itoa(change) + "/detail"

	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + strconv.Itoa(change) + "/detail"
	}

	return buf
}

func (r *review) urlFiles(change, revision int) string {
	r.cfg.Logger.Debug("review: urlFiles")

	buf := r.url + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/files/"

	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/files/"
	}

	return buf
}

func (r *review) urlPatch(change, revision int) string {
	r.cfg.Logger.Debug("review: urlPatch")

	buf := r.url + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/patch"

	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/patch"
	}

	return buf
}

func (r *review) urlQuery(search string, option []string, start int) string {
	r.cfg.Logger.Debug("review: urlQuery")

	query := "?q=" + search + "&o=" + strings.Join(option, "&o=") + "&n=" + strconv.Itoa(start)

	buf := r.url + "/changes/" + query
	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + query
	}

	return buf
}

func (r *review) urlReview(change, revision int) string {
	r.cfg.Logger.Debug("review: urlReview")

	buf := r.url + "/changes/" + strconv.Itoa(change) +
		"/revisions/" + strconv.Itoa(revision) + "/review"

	if r.user != "" && r.pass != "" {
		buf = r.url + "/a/changes/" + strconv.Itoa(change) +
			"/revisions/" + strconv.Itoa(revision) + "/review"
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
