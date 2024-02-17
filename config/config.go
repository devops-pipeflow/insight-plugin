package config

type Config struct {
	ApiVersion string   `yaml:"apiVersion"`
	Kind       string   `yaml:"kind"`
	MetaData   MetaData `yaml:"metadata"`
	Spec       Spec     `yaml:"spec"`
}

type MetaData struct {
	Name string `yaml:"name"`
}

type Spec struct {
	EnvVariables   []EnvVariable  `yaml:"envVariables"`
	BuildConfig    BuildConfig    `yaml:"buildConfig"`
	CodeConfig     CodeConfig     `yaml:"codeConfig"`
	NodeConfig     NodeConfig     `yaml:"nodeConfig"`
	ArtifactConfig ArtifactConfig `yaml:"artifactConfig"`
	GptConfig      GptConfig      `yaml:"gptConfig"`
	RepoConfig     RepoConfig     `yaml:"repoConfig"`
	ReviewConfig   ReviewConfig   `yaml:"reviewConfig"`
	SshConfig      SshConfig      `yaml:"sshConfig"`
}

type EnvVariable struct {
	Name  string `yaml:"name"`
	Value string `yaml:"value"`
}

type BuildConfig struct {
	LoggingConfig LoggingConfig `yaml:"loggingConfig"`
}

type CodeConfig struct {
}

type NodeConfig struct {
	Duration string `yaml:"duration"`
}

type ArtifactConfig struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type GptConfig struct {
}

type RepoConfig struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type ReviewConfig struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type SshConfig struct {
	Host    string `yaml:"host"`
	Port    int64  `yaml:"port"`
	User    string `yaml:"user"`
	Pass    string `yaml:"pass"`
	Key     string `yaml:"key"`
	Timeout string `yaml:"timeout"`
}

type LoggingConfig struct {
	Start int64 `yaml:"start"`
	Len   int64 `yaml:"len"`
	Count int64 `yaml:"count"`
}

func New() *Config {
	return &Config{}
}
