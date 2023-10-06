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
	Sights []Sight `yaml:"sights"`
	Repo   Repo    `yaml:"repo"`
	Review Review  `yaml:"review"`
	Gpt    Gpt     `yaml:"gpt"`
}

type Sight struct {
	Name   string `yaml:"name"`
	Enable bool   `yaml:"enable"`
}

type Repo struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Review struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

type Gpt struct {
	Url  string `yaml:"url"`
	User string `yaml:"user"`
	Pass string `yaml:"pass"`
}

var (
	Build   string
	Version string
)

func New() *Config {
	return &Config{}
}
