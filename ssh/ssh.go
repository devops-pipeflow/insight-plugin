package ssh

import (
	"bytes"
	"context"
	"net"
	"os"
	"strconv"
	"time"

	"github.com/hashicorp/go-hclog"
	"github.com/pkg/errors"
	crypto_ssh "golang.org/x/crypto/ssh"

	"github.com/devops-pipeflow/insight-plugin/config"
)

const (
	connDuration = 30 * time.Second
	connInternal = 1

	ptyHeight = 40
	ptyWidth  = 80
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

type Ssh interface {
	Init(context.Context) error
	Deinit(context.Context) error
	Run(context.Context, []string) error
}

type SshConfig struct {
	Config config.Config
	Logger hclog.Logger
}

type ssh struct {
	cfg     *SshConfig
	session *crypto_ssh.Session
	host    string
	port    int64
	user    string
	pass    string
	key     string
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

	session, err := s.initSession(ctx)
	if err != nil {
		return errors.Wrap(err, "failed to init session")
	}

	s.session = session

	return nil
}

func (s *ssh) Deinit(_ context.Context) error {
	s.cfg.Logger.Debug("ssh: Deinit")

	return s.session.Close()
}

func (s *ssh) Run(ctx context.Context, cmds []string) error {
	s.cfg.Logger.Debug("ssh: Run")

	return s.runSession(ctx, cmds)
}

func (s *ssh) initSession(_ context.Context) (*crypto_ssh.Session, error) {
	s.cfg.Logger.Debug("ssh: initSession")

	var cfg crypto_ssh.Config
	var signer crypto_ssh.Signer

	auth := make([]crypto_ssh.AuthMethod, 0)

	if s.key != "" {
		pem, err := os.ReadFile(s.key)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read key")
		}
		if s.pass != "" {
			signer, err = crypto_ssh.ParsePrivateKeyWithPassphrase(pem, []byte(s.pass))
		} else {
			signer, err = crypto_ssh.ParsePrivateKey(pem)
		}
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse private key")
		}
		auth = append(auth, crypto_ssh.PublicKeys(signer))
	} else {
		auth = append(auth, crypto_ssh.Password(s.pass))
	}

	cfg.Ciphers = []string{
		cipher3DesCbc, cipherAes128Cbc, cipherAes128Ctr, cipherAes128Gcm,
		cipherAes192Cbc, cipherAes192Ctr, cipherAes256Cbc, cipherAes256Ctr,
		cipherArcFour128, cipherArcFour256,
	}

	cfg.KeyExchanges = []string{
		keyExchangeDiffieHellmanGroup1Sha1,
		keyExchangeDiffieHellmanGroupExchangeSha1,
		keyExchangeDiffieHellmanGroupExchangeSha256,
	}

	timeout := connDuration
	if s.cfg.Config.Spec.NodeConfig.Duration != "" {
		t, err := strconv.ParseInt(s.cfg.Config.Spec.NodeConfig.Duration, 10, 64)
		if err != nil {
			return nil, errors.Wrap(err, "failed to parse int")
		}
		timeout = time.Duration(t)
	}

	_config := &crypto_ssh.ClientConfig{
		User:    s.user,
		Auth:    auth,
		Timeout: timeout,
		Config:  cfg,
		HostKeyCallback: func(hostname string, remote net.Addr, key crypto_ssh.PublicKey) error {
			return nil
		},
	}

	conn, err := crypto_ssh.Dial("tcp", s.host+":"+strconv.FormatInt(s.port, 10), _config)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ssh client")
	}

	defer func(conn *crypto_ssh.Client) {
		_ = conn.Close()
	}(conn)

	session, err := conn.NewSession()
	if err != nil {
		return nil, errors.Wrap(err, "failed to create ssh session")
	}

	defer func(session *crypto_ssh.Session) {
		_ = session.Close()
	}(session)

	modes := crypto_ssh.TerminalModes{
		crypto_ssh.ECHO:          0,     // disable echoing
		crypto_ssh.TTY_OP_ISPEED: 14400, // input speed = 14.4kbaud
		crypto_ssh.TTY_OP_OSPEED: 14400, // output speed = 14.4kbaud
	}

	if err := session.RequestPty("xterm", ptyHeight, ptyWidth, modes); err != nil {
		return nil, errors.Wrap(err, "failed to request ssh pty")
	}

	return session, nil
}

func (s *ssh) runSession(_ context.Context, cmds []string) error {
	s.cfg.Logger.Debug("ssh: runSession")

	var outBuf, errBuf bytes.Buffer

	stdinBuf, _ := s.session.StdinPipe()
	s.session.Stdout = &outBuf
	s.session.Stderr = &errBuf

	err := s.session.Shell()
	if err != nil {
		return errors.Wrap(err, "failed to run shell")
	}

	cmds = append(cmds, "exit 0")
	for _, item := range cmds {
		item += "\n"
		_, err = stdinBuf.Write([]byte(item))
		if err != nil {
			return errors.Wrap(err, "failed to write buffer")
		}
	}

	err = s.session.Wait()
	if err != nil {
		return errors.Wrap(err, "failed to wait session")
	}

	if errBuf.String() != "" {
		return errors.New(errBuf.String())
	}

	return nil
}
