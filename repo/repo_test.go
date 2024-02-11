//go:build repo_test

// go test -cover -covermode=atomic -parallel 2 -tags=repo_test -v github.com/devops-pipeflow/insight-plugin/repo

package repo

import (
	"io"
	"os"
	"testing"

	"github.com/hashicorp/go-hclog"
	"github.com/stretchr/testify/assert"
	"gopkg.in/yaml.v3"

	"github.com/devops-pipeflow/insight-plugin/config"
)

func initRepo() repo {
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

	return repo{
		cfg: &Config{
			Config: *c,
			Logger: hclog.New(&hclog.LoggerOptions{
				Name:  "insight",
				Level: hclog.LevelFromString("DEBUG"),
			}),
		},
		user: c.Spec.RepoConfig.User,
		pass: c.Spec.RepoConfig.Pass,
		url:  c.Spec.RepoConfig.Url,
	}
}

func TestGet(t *testing.T) {
	r := initRepo()

	_, err := r.Get("platform/build/soong", "branch:master")
	assert.Equal(t, nil, err)

	_, err = r.Get("platform/build/soong", "commit:42ada5cff3fca011b5a0d017955f14dc63898807")
	assert.Equal(t, nil, err)
}

func TestQuery(t *testing.T) {
	r := initRepo()

	_, err := r.Query("platform/build/soong", "branch:master")
	assert.Equal(t, nil, err)

	_, err = r.Query("platform/build/soong", "branch:main commit:42ada5cff3fca011b5a0d017955f14dc63898807")
	assert.Equal(t, nil, err)

	_, err = r.Query("platform/build/soong", "tag:android-vts-10.0_r4")
	assert.Equal(t, nil, err)

	_, err = r.Query("platform/build/soong", "tag:android-vts-10.0_r4 commit:42ada5cff3fca011b5a0d017955f14dc63898807")
	assert.Equal(t, nil, err)
}

func TestRequest(t *testing.T) {
	assert.Equal(t, nil, nil)
}
