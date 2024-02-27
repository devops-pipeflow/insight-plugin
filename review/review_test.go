//go:build review_test

// go test -cover -covermode=atomic -parallel 2 -tags=review_test -v github.com/devops-pipeflow/insight-plugin/review

package review

import (
	"context"
	"encoding/base64"
	"io"
	"net/url"
	"os"
	"path/filepath"
	"strconv"
	"strings"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	changeGerrit   = 883543
	commitGerrit   = "5907d4189ff8e798a9914186c91e4bf7b3166973"
	revisionGerrit = 17
)

func initReview() review {
	helper := func(name string) *config.Config {
		c := config.New()
		fi, _ := os.Open(name)
		defer func() {
			_ = fi.Close()
		}()
		buf, _ := io.ReadAll(fi)
		_ = yaml.Unmarshal(buf, c)
		return c
	}

	c := helper("../test/config/config.yml")

	return review{
		cfg: &Config{
			Config: *c,
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:  "insight",
				Level: hclog.LevelFromString("DEBUG"),
			}),
		},
		user: c.Spec.ReviewConfig.User,
		pass: c.Spec.ReviewConfig.Pass,
		url:  c.Spec.ReviewConfig.Url,
	}
}

func TestClean(t *testing.T) {
	d, _ := os.Getwd()
	root := filepath.Join(d, "gerrit-test-clean")

	_ = os.Mkdir(root, os.ModePerm)

	r := initReview()
	err := r.Clean(context.Background(), root)
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestFetch(t *testing.T) {
	d, _ := os.Getwd()
	root := filepath.Join(d, "gerrit-test-fetch")

	ctx := context.Background()
	r := initReview()

	_, _, _, err := r.Fetch(ctx, root, commitGerrit)
	assert.Equal(t, nil, err)

	err = r.Clean(ctx, root)
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestQuery(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.Query(ctx, "change:"+strconv.Itoa(changeGerrit))
	assert.Equal(t, nil, err)
}

// nolint: dogsled
func TestVote(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	buf := make([]Format, 0)

	err := r.Vote(ctx, "", buf)
	assert.NotEqual(t, nil, err)

	err = r.Vote(ctx, commitGerrit, buf)
	assert.Equal(t, nil, err)

	buf = make([]Format, 1)
	buf[0] = Format{
		Details: "Disapproved",
		File:    "Android.mk",
		Line:    1,
		Type:    TypeError,
	}

	err = r.Vote(ctx, commitGerrit, buf)
	assert.Equal(t, nil, err)
}

func TestWrite(t *testing.T) {
	r := initReview()

	d, _ := os.Getwd()
	err := r.write(d, "gerrit-test-write", "Hello World!")
	assert.Equal(t, nil, err)

	_ = os.RemoveAll(filepath.Join(d, "gerrit-test-write"))
}

func TestUnmarshal(t *testing.T) {
	r := initReview()
	data := ")]}'{\"project\": \"myProject\",\"branch\": \"master\"}"

	buf, err := r.unmarshal([]byte(data))
	assert.Equal(t, nil, err)
	assert.Equal(t, "myProject", buf["project"])
	assert.Equal(t, "master", buf["branch"])
}

func TestUnmarshalList(t *testing.T) {
	r := initReview()
	data := ")]}'[{\"project\":\"demo1\",\"branch\":\"master1\"},{\"project\":\"demo2\",\"branch\":\"master2\"}]"

	buf, err := r.unmarshalList([]byte(data))
	assert.Equal(t, nil, err)
	assert.Equal(t, "demo1", buf["project"])
	assert.Equal(t, "master1", buf["branch"])
}

func TestUrlContent(t *testing.T) {
	r := initReview()

	change := 1
	revision := 2
	name := "foo.txt"
	_url := r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles + url.QueryEscape(name) + urlContent

	buf := r.urlContent(change, revision, name)
	assert.Equal(t, _url, buf)
}

func TestUrlDetail(t *testing.T) {
	r := initReview()

	change := 1
	_url := r.url + urlPrefix + urlChanges + strconv.Itoa(change) + urlDetail

	buf := r.urlDetail(change)
	assert.Equal(t, _url, buf)
}

func TestUrlFiles(t *testing.T) {
	r := initReview()

	change := 1
	revision := 2
	_url := r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlFiles

	buf := r.urlFiles(change, revision)
	assert.Equal(t, _url, buf)
}

func TestUrlPatch(t *testing.T) {
	r := initReview()

	change := 1
	revision := 2
	_url := r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlPatch

	buf := r.urlPatch(change, revision)
	assert.Equal(t, _url, buf)
}

func TestUrlQuery(t *testing.T) {
	r := initReview()

	search := "is:open"
	option := []string{"LABELS"}
	start := 1
	_url := r.url + urlPrefix + urlChanges + urlQuery + search + urlOption + strings.Join(option, urlOption) + urlNumber + strconv.Itoa(start)

	buf := r.urlQuery(search, option, start)
	assert.Equal(t, _url, buf)
}

func TestUrlReview(t *testing.T) {
	r := initReview()

	change := 1
	revision := 2
	_url := r.url + urlPrefix + urlChanges + strconv.Itoa(change) +
		urlRevisions + strconv.Itoa(revision) + urlReview

	buf := r.urlReview(change, revision)
	assert.Equal(t, _url, buf)
}

func TestGetContent(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.get(ctx, r.urlContent(-1, -1, ""))
	assert.NotEqual(t, nil, err)

	buf, err := r.get(ctx, r.urlContent(changeGerrit, revisionGerrit, url.PathEscape("Android.mk")))
	assert.Equal(t, nil, err)

	dst := make([]byte, len(buf))
	n, _ := base64.StdEncoding.Decode(dst, buf)
	assert.NotEqual(t, 0, n)
}

func TestGetDetail(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.get(ctx, r.urlDetail(-1))
	assert.NotEqual(t, nil, err)

	buf, err := r.get(ctx, r.urlDetail(changeGerrit))
	assert.Equal(t, nil, err)

	_, err = r.unmarshal(buf)
	assert.Equal(t, nil, err)
}

func TestGetFiles(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.get(ctx, r.urlFiles(-1, -1))
	assert.NotEqual(t, nil, err)

	buf, err := r.get(ctx, r.urlFiles(changeGerrit, revisionGerrit))
	assert.Equal(t, nil, err)

	_, err = r.unmarshal(buf)
	assert.Equal(t, nil, err)
}

func TestGetPatch(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.get(ctx, r.urlPatch(-1, -1))
	assert.NotEqual(t, nil, err)

	buf, err := r.get(ctx, r.urlPatch(changeGerrit, revisionGerrit))
	assert.Equal(t, nil, err)

	dst := make([]byte, len(buf))
	n, _ := base64.StdEncoding.Decode(dst, buf)
	assert.NotEqual(t, 0, n)
}

func TestGetQuery(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	_, err := r.get(ctx, r.urlQuery("commit:-1", []string{"CURRENT_REVISION"}, 0))
	assert.NotEqual(t, nil, err)

	buf, err := r.get(ctx, r.urlQuery("commit:"+commitGerrit, []string{"CURRENT_REVISION"}, 0))
	assert.Equal(t, nil, err)

	_, err = r.unmarshalList(buf)
	assert.Equal(t, nil, err)
}

func TestPostReview(t *testing.T) {
	ctx := context.Background()
	r := initReview()

	err := r.post(ctx, r.urlReview(-1, -1), nil)
	assert.NotEqual(t, nil, err)

	buf := map[string]interface{}{
		"comments": map[string]interface{}{
			"Android.mk": []map[string]interface{}{
				{
					"line":    1,
					"message": "Commented by pipeflow insight",
				},
			},
		},
		"labels": map[string]interface{}{
			"Code-Review": -1,
		},
		"message": "Voting Code-Review by pipeflow insight",
	}

	err = r.post(ctx, r.urlReview(changeGerrit, revisionGerrit), buf)
	assert.Equal(t, nil, err)
}
