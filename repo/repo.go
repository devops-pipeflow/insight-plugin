package repo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"net/url"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	urlConcat = "/+/"
	urlJSON   = "format=JSON"
	urlHeads  = "refs/heads/"
	urlLog    = "/+log/"
	urlSearch = "/?s="
	urlTags   = "refs/tags/"
	urlText   = "format=TEXT"
)

const (
	opBranch    = "branch:"
	opCommit    = "commit:"
	opDelimiter = " "
	opGroups    = 2
	opTag       = "tag:"
)

type Repo interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Fetch(context.Context, string, string, string) ([]byte, error)
	Get(context.Context, string, string) (map[string]interface{}, error)
	Query(context.Context, string, string) (map[string]interface{}, error)
}

type Config struct {
	Config config.Config
	Logger hclog.Logger
}

type repo struct {
	cfg  *Config
	user string
	pass string
	url  string
}

func New(_ context.Context, cfg *Config) Repo {
	return &repo{
		cfg: cfg,
	}
}

func DefaultConfig() *Config {
	return &Config{}
}

func (r *repo) Init(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Init")

	r.user = r.cfg.Config.Spec.RepoConfig.User
	r.pass = r.cfg.Config.Spec.RepoConfig.Pass
	r.url = r.cfg.Config.Spec.RepoConfig.Url

	return nil
}

func (r *repo) Deinit(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Deinit")

	return nil
}

// Fetch
//
// Example:
//
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+/refs/heads/main/README.md
//
// commit:COMMIT: https://android.googlesource.com/platform/build/soong/+/25900543331a1508110da4926ca45557b4c236da/README.md
//
// tag:TAG: https://android.googlesource.com/platform/build/soong/+/refs/heads/android14-release/README.md
func (r *repo) Fetch(_ context.Context, project, file, operator string) ([]byte, error) {
	r.cfg.Logger.Debug("repo: Fetch")

	var buf []byte
	var err error

	if project == "" || file == "" || operator == "" || len(strings.Split(operator, opDelimiter)) >= opGroups {
		return nil, errors.New("parameter invalid")
	}

	if strings.HasPrefix(operator, opBranch) {
		branch := strings.TrimPrefix(operator, opBranch)
		buf, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+urlHeads+branch+"/"+file+"?"+urlText, r.user, r.pass)
	} else if strings.HasPrefix(operator, opCommit) {
		commit := strings.TrimPrefix(operator, opCommit)
		buf, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+commit+"/"+file+"?"+urlText, r.user, r.pass)
	} else if strings.HasPrefix(operator, opTag) {
		tag := strings.TrimPrefix(operator, opTag)
		buf, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+urlHeads+tag+"/"+file+"?"+urlText, r.user, r.pass)
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

// Get
//
// Example:
//
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+/refs/heads/main?format=JSON
//
// commit:COMMIT: https://android.googlesource.com/platform/build/soong/+/42ada5cff3fca011b5a0d017955f14dc63898807?format=JSON
//
// tag:TAG: https://android.googlesource.com/platform/build/soong/+/refs/tags/android-vts-10.0_r4?format=JSON
//
// nolint: lll
func (r *repo) Get(_ context.Context, project, operator string) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("repo: Get")

	var body []byte
	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" || len(strings.Split(operator, opDelimiter)) >= opGroups {
		return nil, errors.New("parameter invalid")
	}

	if strings.HasPrefix(operator, opBranch) {
		branch := strings.TrimPrefix(operator, opBranch)
		body, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+urlHeads+branch+"?"+urlJSON, r.user, r.pass)
	} else if strings.HasPrefix(operator, opCommit) {
		commit := strings.TrimPrefix(operator, opCommit)
		body, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+commit+"?"+urlJSON, r.user, r.pass)
	} else if strings.HasPrefix(operator, opTag) {
		tag := strings.TrimPrefix(operator, opTag)
		body, err = r.get(r.url+"/"+url.PathEscape(project)+urlConcat+urlTags+tag+"?"+urlJSON, r.user, r.pass)
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	body = []byte(strings.ReplaceAll(string(body), ")]}'", ""))

	if err := json.Unmarshal(body, &buf); err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	return buf, nil
}

// Query
//
// Example:
//
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+log/refs/heads/main?format=JSON
//
// branch:BRANCH commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/heads/main/?s=42ada5cff3fca011b5a0d017955f14dc63898807&format=JSON
//
// tag:TAG: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4?format=JSON
//
// tag:TAG commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4/?s=9863d53618714a36c3f254d949497a7eb2d11863&format=JSON
//
// nolint: gocyclo,lll
func (r *repo) Query(_ context.Context, project, operator string) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("repo: Query")

	parser := func(op string) (string, string, string, error) {
		var branch, commit, tag string

		buf := strings.Split(op, opDelimiter)
		if len(buf) > opGroups {
			return "", "", "", errors.New("operator invalid")
		}

		for _, val := range buf {
			if strings.HasPrefix(val, opBranch) {
				branch = strings.TrimPrefix(val, opBranch)
			} else if strings.HasPrefix(val, opCommit) {
				commit = strings.TrimPrefix(val, opCommit)
			} else if strings.HasPrefix(val, opTag) {
				tag = strings.TrimPrefix(val, opTag)
			} else {
				continue
			}
		}

		if branch != "" && tag != "" {
			return "", "", "", errors.New("operator invalid")
		}

		if len(buf) == 1 && commit != "" {
			return "", "", "", errors.New("operator invalid")
		}

		return branch, commit, tag, nil
	}

	var body []byte
	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" {
		return nil, errors.New("parameter invalid")
	}

	branch, commit, tag, err := parser(operator)
	if err != nil {
		return nil, err
	}

	if branch != "" {
		if commit != "" {
			body, err = r.get(r.url+"/"+url.PathEscape(project)+urlLog+urlHeads+branch+urlSearch+commit+"&"+urlJSON, r.user, r.pass)
		} else {
			body, err = r.get(r.url+"/"+url.PathEscape(project)+urlLog+urlHeads+branch+"?"+urlJSON, r.user, r.pass)
		}
	} else if tag != "" {
		if commit != "" {
			body, err = r.get(r.url+"/"+url.PathEscape(project)+urlLog+urlTags+tag+urlSearch+commit+"&"+urlJSON, r.user, r.pass)
		} else {
			body, err = r.get(r.url+"/"+url.PathEscape(project)+urlLog+urlTags+tag+"?"+urlJSON, r.user, r.pass)
		}
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	body = []byte(strings.ReplaceAll(string(body), ")]}'", ""))

	if err := json.Unmarshal(body, &buf); err != nil {
		return nil, errors.Wrap(err, "unmarshal failed")
	}

	return buf, nil
}

func (r *repo) get(_url, user, pass string) ([]byte, error) {
	r.cfg.Logger.Debug("repo: get")

	req, err := http.NewRequest(http.MethodGet, _url, http.NoBody)
	if err != nil {
		return nil, errors.Wrap(err, "request failed")
	}

	if user != "" && pass != "" {
		req.SetBasicAuth(user, pass)
	}

	resp, err := http.DefaultClient.Do(req)
	if err != nil {
		return nil, errors.Wrap(err, "client failed")
	}

	defer func() {
		_ = resp.Body.Close()
	}()

	if resp.StatusCode != http.StatusOK {
		return nil, errors.New("client failed")
	}

	data, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read failed")
	}

	return data, nil
}
