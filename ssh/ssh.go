package ssh

import (
	"context"
	"net"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	crypto_ssh "golang.org/x/crypto/ssh"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	sshTimeout = 10 * time.Second
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
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, string) (string, error)
}

type SshConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type ssh struct {
	cfg     *SshConfig
	client  *crypto_ssh.Client
	session *crypto_ssh.Session
}

func New(_ context.Context, cfg *SshConfig) Ssh {
	return &ssh{
		cfg: cfg,
	}
}

func DefaultConfig() *SshConfig {
	return &SshConfig{}
}

func (s *ssh) Init(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Init")

	return s.initSession(ctx)
}

func (s *ssh) Deinit(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: Deinit")

	return s.deinitSession(ctx)
}

func (s *ssh) Run(ctx context.Context, cmd string) (string, error) {
	s.cfg.Logger.Debug("ssh: Run")

	return s.runSession(ctx, cmd)
}

func (s *ssh) initSession(ctx context.Context) error {
	s.cfg.Logger.Debug("ssh: initSession")

	var cfg crypto_ssh.Config
	var err error

	auth, err := s.setAuth(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to set auth")
	}

	cfg.Ciphers = ciphers
	cfg.KeyExchanges = keyExchanges

	timeout, err := s.setTimeout(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to set timeout")
	}

	_config := &crypto_ssh.ClientConfig{
		User:    s.cfg.Config.Spec.SshConfig.User,
		Auth:    auth,
		Timeout: timeout,
		Config:  cfg,
		HostKeyCallback: func(hostname string, remote net.Addr, key crypto_ssh.PublicKey) error {
			return nil
		},
	}

	addr := s.cfg.Config.Spec.SshConfig.Host + ":" + strconv.FormatInt(s.cfg.Config.Spec.SshConfig.Port, 10)

	s.client, err = crypto_ssh.Dial("tcp", addr, _config)
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

func (s *ssh) runSession(_ context.Context, cmd string) (string, error) {
	s.cfg.Logger.Debug("ssh: runSession")

	if s.session == nil {
		return "", errors.New("invalid session")
	}

	out, err := s.session.CombinedOutput(cmd)
	if err != nil {
		return string(out), errors.Wrap(err, "failed to run cmd")
	}

	return string(out), nil
}

func (s *ssh) setAuth(_ context.Context) ([]crypto_ssh.AuthMethod, error) {
	s.cfg.Logger.Debug("ssh: setAuth")

	var err error
	var signer crypto_ssh.Signer

	auth := make([]crypto_ssh.AuthMethod, 0)

	if s.cfg.Config.Spec.SshConfig.Key != "" {
		if s.cfg.Config.Spec.SshConfig.Pass != "" {
			signer, err = crypto_ssh.ParsePrivateKeyWithPassphrase([]byte(s.cfg.Config.Spec.SshConfig.Key), []byte(s.cfg.Config.Spec.SshConfig.Pass))
		} else {
			signer, err = crypto_ssh.ParsePrivateKey([]byte(s.cfg.Config.Spec.SshConfig.Key))
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse private key")
		}
		auth = append(auth, crypto_ssh.PublicKeys(signer))
	} else {
		auth = append(auth, crypto_ssh.Password(s.cfg.Config.Spec.SshConfig.Pass))
	}

	return auth, nil
}

func (s *ssh) setTimeout(_ context.Context) (time.Duration, error) {
	s.cfg.Logger.Debug("ssh: setTimeout")

	var timeout time.Duration
	var err error

	if s.cfg.Config.Spec.SshConfig.Timeout != "" {
		timeout, err = time.ParseDuration(s.cfg.Config.Spec.SshConfig.Timeout)
		if err != nil {
			return 0, errors.Wrap(err, "failed to parse duration")
		}
	} else {
		timeout = sshTimeout
	}

	return timeout, nil
}
