package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/BurntSushi/toml"
)

type tomlConfig struct {
	Login     string
	Password  string
	Host      string
	RemoteDir string `toml:"remote_dir"`
	LocalDir  string `toml:"local_dir"`
	Filters   []string
}

var (
	conf           tomlConfig
	configLocation string
)

func init() {
	flag.StringVar(&configLocation, "config", "config", "Configuration File")
	flag.Parse()

	if _, err := toml.DecodeFile(configLocation, &conf); err != nil {
		fmt.Printf("Could not parse configuration file: %v\n", err)
		os.Exit(1)
	}
}

func main() {
	fmt.Println(conf)
}
