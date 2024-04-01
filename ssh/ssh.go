package ssh

import (
	"context"
	"net"
	"strconv"
	"strings"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	cryptossh "golang.org/x/crypto/ssh"

	"github.com/devops-pipeflow/insight-plugin/config"
	proto "github.com/devops-pipeflow/server/plugins/insight"
)

const (
	connTimeout = 10 * time.Second
	operatorAnd = "&&"
)

const (
	cipher3DesCbc    = "3des-cbc"
	cipherAes128Cbc  = "aes128-cbc"
	cipherAes128Ctr  = "aes128-ctr"
	cipherAes128Gcm  = "aes128-gcm@openssh.com"
	cipherAes192Cbc  = "aes192-cbc"
	cipherAes192Ctr  = "aes192-ctr"
	cipherAes256Cbc  = "aes256-cbc"
	cipherAes256Ctr  = "aes256-ctr"
	cipherArcFour128 = "arcfour128"
	cipherArcFour256 = "arcfour256"

	keyExchangeDiffieHellmanGroup1Sha1          = "diffie-hellman-group1-sha1"
	keyExchangeDiffieHellmanGroupExchangeSha1   = "diffie-hellman-group-exchange-sha1"
	keyExchangeDiffieHellmanGroupExchangeSha256 = "diffie-hellman-group-exchange-sha256"
)

var (
	ciphers = []string{
		cipher3DesCbc, cipherAes128Cbc, cipherAes128Ctr, cipherAes128Gcm,
		cipherAes192Cbc, cipherAes192Ctr, cipherAes256Cbc, cipherAes256Ctr,
		cipherArcFour128, cipherArcFour256,
	}

	keyExchanges = []string{
		keyExchangeDiffieHellmanGroup1Sha1,
		keyExchangeDiffieHellmanGroupExchangeSha1,
		keyExchangeDiffieHellmanGroupExchangeSha256,
	}
)

type Ssh interface {
	Init(context.Context, *proto.SshConfig) error
	Deinit(context.Context) error
	Run(context.Context, []string) (string, error)
}

type SshConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type ssh struct {
	cfg     *SshConfig
	client  *cryptossh.Client
	session *cryptossh.Session
}

func New(_ context.Context, cfg *SshConfig) Ssh {
	return &ssh{
		cfg: cfg,
	}
}

func DefaultConfig() *SshConfig {
	return &SshConfig{}
}

func (s *ssh) Init(ctx context.Context, cfg *proto.SshConfig) error {
	s.cfg.Logger.Debug("ssh: Init")

	if cfg.Host != "" {
		s.cfg.Config.Spec.SshConfig = config.SshConfig(*cfg)
	}

	return nil
}

func (s *ssh) Deinit(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Deinit")

	return nil
}

func (s *ssh) Run(ctx context.Context, cmds []string) (string, error) {
	s.cfg.Logger.Debug("ssh: Run")

	c := s.cfg.Config.Spec.SshConfig

	if err := s.initSession(ctx, c.Host, c.Port, c.User, c.Pass, c.Key, c.Timeout); err != nil {
		return "", errors.Wrap(err, "failed to init session")
	}

	defer func(s *ssh, ctx context.Context) {
		_ = s.deinitSession(ctx)
	}(s, ctx)

	return s.runSession(ctx, cmds)
}

func (s *ssh) initSession(ctx context.Context, host string, port int64, user, pass, key, timeout string) error {
	s.cfg.Logger.Debug("ssh: initSession")

	var cfg cryptossh.Config
	var err error

	auth, err := s.setAuth(ctx, pass, key)
	if err != nil {
		return errors.Wrap(err, "failed to set auth")
	}

	cfg.Ciphers = ciphers
	cfg.KeyExchanges = keyExchanges

	t, err := s.setTimeout(ctx, timeout)
	if err != nil {
		return errors.Wrap(err, "failed to set timeout")
	}

	_config := &cryptossh.ClientConfig{
		User:    user,
		Auth:    auth,
		Timeout: t,
		Config:  cfg,
		HostKeyCallback: func(hostname string, remote net.Addr, key cryptossh.PublicKey) error {
			return nil
		},
	}

	addr := host + ":" + strconv.FormatInt(port, 10)

	s.client, err = cryptossh.Dial("tcp", addr, _config)
	if err != nil {
		return errors.Wrap(err, "failed to create ssh client")
	}

	s.session, err = s.client.NewSession()
	if err != nil {
		_ = s.client.Close()
		return errors.Wrap(err, "failed to create ssh session")
	}

	return nil
}

func (s *ssh) deinitSession(_ context.Context) error {
	s.cfg.Logger.Debug("ssh: deinitSession")

	if s.session != nil {
		_ = s.session.Close()
	}

	if s.client != nil {
		_ = s.client.Close()
	}

	return nil
}

func (s *ssh) runSession(_ context.Context, cmds []string) (string, error) {
	s.cfg.Logger.Debug("ssh: runSession")

	if s.session == nil {
		return "", errors.New("invalid session")
	}

	out, err := s.session.CombinedOutput(strings.Join(cmds, operatorAnd))
	if err != nil {
		return string(out), errors.Wrap(err, "failed to run cmd")
	}

	return string(out), nil
}

func (s *ssh) setAuth(_ context.Context, pass, key string) ([]cryptossh.AuthMethod, error) {
	s.cfg.Logger.Debug("ssh: setAuth")

	var err error
	var signer cryptossh.Signer

	auth := make([]cryptossh.AuthMethod, 0)

	if key != "" {
		if pass != "" {
			signer, err = cryptossh.ParsePrivateKeyWithPassphrase([]byte(key), []byte(pass))
		} else {
			signer, err = cryptossh.ParsePrivateKey([]byte(key))
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse private key")
		}
		auth = append(auth, cryptossh.PublicKeys(signer))
	} else {
		auth = append(auth, cryptossh.Password(pass))
	}

	return auth, nil
}

func (s *ssh) setTimeout(_ context.Context, timeout string) (time.Duration, error) {
	s.cfg.Logger.Debug("ssh: setTimeout")

	var t time.Duration
	var err error

	if timeout != "" {
		t, err = time.ParseDuration(timeout)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse duration")
		}
	} else {
		t = connTimeout
	}

	return t, nil
}
