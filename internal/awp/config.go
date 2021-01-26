package awpconf

import (
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/gravitational/configure"
)

type Config struct {
	Address  string `env:"ADDRESS" yaml:"address"`
	AsGroup  bool   `env:"AS_GROUP" yaml:"as_group"`
	_asGroup string
	Group    int `env:"GROUP" yaml:"group"`
	_group   string
	Name     string `env:"NAME" yaml:"name"`
	Token    string `env:"TOKEN" yaml:"token"`
	Log      string `env:"LOG_LEVEL" yaml:"log"`
	BodySize int    `env:"BODY_SIZE" yaml:"body_size"`
	Stats    struct {
		Health   string  `env:"HEALTH" yaml:"health"`
		Adequacy float64 `env:"ADEQUACY" yaml:"adequacy"`
	} `yaml:"stats"`
}

func (c *Config) GetAsGroup() *string {
	return &c._asGroup
}

func (c *Config) GetGroup() *string {
	return &c._group
}
func (c *Config) setGroup(g string) {
	c._group = g
}
func (c *Config) setAsGroup(a string) {
	c._asGroup = a
}

var Cfg Config

func Parse() error {
	defaultConf()
	lenArgs := len(os.Args)
	if lenArgs < 2 {
		fmt.Printf("Not enough arguments.\nUsage: %s path/to/config.yaml\n", os.Args[0])
		os.Exit(1)
	} else if os.Args[1] == "-h" || os.Args[1] == "--help" {
		fmt.Printf("Usage: %s path/to/config.yaml\n", os.Args[0])
		os.Exit(0)
	}
	yaml, err := ioutil.ReadFile(os.Args[1])
	if err != nil {
		return err
	}
	err = configure.ParseYAML(yaml, &Cfg)
	if err != nil {
		return err
	}
	err = configure.ParseEnv(&Cfg)
	if err != nil {
		return err
	}
	normalize()
	return nil
}

func defaultConf() {
	Cfg.Address = ":8080"
	Cfg.Name = "Anon Wall Poster"
	Cfg.Log = "info"
	Cfg.Group = -1
	Cfg.AsGroup = true
}

func normalize() {
	if Cfg.AsGroup {
		Cfg.setAsGroup("1")
	} else {
		Cfg.setAsGroup("0")
	}
	Cfg.setGroup(fmt.Sprint(Cfg.Group))
	Cfg.Log = strings.ToLower(Cfg.Log)
}
