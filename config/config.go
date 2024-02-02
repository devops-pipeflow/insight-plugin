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
	EnvVariables []EnvVariable `yaml:"envVariables"`
	BuildConfig  BuildConfig   `yaml:"buildConfig"`
	CodeConfig   CodeConfig    `yaml:"codeConfig"`
	GptConfig    GptConfig     `yaml:"gptConfig"`
	NodeConfig   NodeConfig    `yaml:"nodeConfig"`
	RepoConfig   RepoConfig    `yaml:"repoConfig"`
	ReviewConfig ReviewConfig  `yaml:"reviewConfig"`
}

type EnvVariable struct {
	Key   string `yaml:"key"`
	Value string `yaml:"value"`
}

type BuildConfig struct {
	LoggingConfig LoggingConfig `yaml:"loggingConfig"`
}

type CodeConfig struct {
}

type GptConfig struct {
}

type NodeConfig struct {
	Duration string `yaml:"duration"`
	Interval string `yaml:"interval"`
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

type LoggingConfig struct {
	Start int64 `yaml:"start"`
	Len   int64 `yaml:"len"`
	Count int64 `yaml:"count"`
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
