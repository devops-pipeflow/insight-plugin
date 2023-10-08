package repo

import (
	"context"
	"encoding/json"
	"io"
	"net/http"
	"strings"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"

	"github.com/devops-pipeflow/insight-plugin/config"
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
	Get(string, string) (map[string]interface{}, error)
	Query(string, string) (map[string]interface{}, error)
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

	r.user = r.cfg.Config.Spec.Review.User
	r.pass = r.cfg.Config.Spec.Review.Pass
	r.url = r.cfg.Config.Spec.Review.Url

	return nil
}

func (r *repo) Deinit(ctx context.Context) error {
	r.cfg.Logger.Debug("repo: Deinit")

	return nil
}

// nolint: lll
// Example:
// branch:BRANCH: https://android.googlesource.com/platform/build/soong/+/refs/heads/master?format=JSON
// commit:COMMIT: https://android.googlesource.com/platform/build/soong/+/42ada5cff3fca011b5a0d017955f14dc63898807?format=JSON
//
//	tag:TAG: https://android.googlesource.com/platform/build/soong/+/refs/tags/android-vts-10.0_r4
func (r *repo) Get(project, operator string) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("repo: Get")

	var buf map[string]interface{}
	var err error

	if project == "" || operator == "" || len(strings.Split(operator, opDelimiter)) >= opGroups {
		return nil, errors.New("parameter invalid")
	}

	if strings.HasPrefix(operator, opBranch) {
		branch := strings.TrimPrefix(operator, opBranch)
		buf, err = r.request(r.url+"/"+project+"/+/"+"refs/heads/"+branch+"?format=JSON", r.user, r.pass)
	} else if strings.HasPrefix(operator, opCommit) {
		commit := strings.TrimPrefix(operator, opCommit)
		buf, err = r.request(r.url+"/"+project+"/+/"+commit+"?format=JSON", r.user, r.pass)
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

// nolint: gocyclo,lll
// Example:
//
//	branch:BRANCH: https://android.googlesource.com/platform/build/soong/+log/refs/heads/master?format=JSON
//
// branch:BRANCH commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/heads/master/?s=42ada5cff3fca011b5a0d017955f14dc63898807&format=JSON
//
//	              tag:TAG: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4?format=JSON
//	tag:TAG commit:COMMIT: https://android.googlesource.com/platform/build/soong/+log/refs/tags/android-vts-10.0_r4/?s=9863d53618714a36c3f254d949497a7eb2d11863&format=JSON
func (r *repo) Query(project, operator string) (map[string]interface{}, error) {
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
			buf, err = r.request(r.url+"/"+project+"/+log/"+"refs/heads/"+branch+"/?s="+commit+"&format=JSON", r.user, r.pass)
		} else {
			buf, err = r.request(r.url+"/"+project+"/+log/"+"refs/heads/"+branch+"?format=JSON", r.user, r.pass)
		}
	} else if tag != "" {
		if commit != "" {
			buf, err = r.request(r.url+"/"+project+"/+log/"+"refs/tags/"+tag+"/?s="+commit+"&format=JSON", r.user, r.pass)
		} else {
			buf, err = r.request(r.url+"/"+project+"/+log/"+"refs/tags/"+tag+"?format=JSON", r.user, r.pass)
		}
	} else {
		err = errors.New("operator invalid")
	}

	if err != nil {
		return nil, err
	}

	return buf, nil
}

func (r *repo) request(url, user, pass string) (map[string]interface{}, error) {
	r.cfg.Logger.Debug("repo: request")

	var buf map[string]interface{}

	req, err := http.NewRequest(http.MethodGet, url, http.NoBody)
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

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		return nil, errors.New("read failed")
	}

	body = []byte(strings.ReplaceAll(string(body), ")]}'", ""))

	if err := json.Unmarshal(body, &buf); err != nil {
		return nil, errors.New("unmarshal failed")
	}

	return buf, nil
}
