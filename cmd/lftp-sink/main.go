package main

import (
	"flag"
	"fmt"
	"os"
	"os/exec"
	"regexp"
	"strings"

	"github.com/BurntSushi/toml"
	sink "github.com/adityaxdiwakar/lftp-sinking"
)

type tomlConfig struct {
	Login     string
	Password  string
	Host      string
	RemoteDir string `toml:"remote_dir"`
	LocalDir  string `toml:"local_dir"`
	Threads   int
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

func buildMirror(name string) string {
	name = strings.ReplaceAll(name, " ", "\\ ")

	return fmt.Sprintf("mirror --delete --use-pget-n=%d sftp://%s:%s@%s%s%s",
		conf.Threads, conf.Login, conf.Password, conf.Host, conf.RemoteDir, name)
}

func buildPget(name string) string {
	name = strings.ReplaceAll(name, " ", "\\ ")

	return fmt.Sprintf("pget -n %d -c sftp://%s:%s@%s%s%s",
		conf.Threads, conf.Login, conf.Password, conf.Host, conf.RemoteDir, name)
}

func main() {
	folders, files, err := sink.GetDirectoryListing(conf.Host, conf.Login, conf.Password, conf.RemoteDir)
	if err != nil {
		fmt.Println(err)
		os.Exit(1)
	}

	filters := []*regexp.Regexp{}
	for _, filter := range conf.Filters {
		filters = append(filters, regexp.MustCompile(filter))
	}

	foldersFiltered := []string{}
	filesFiltered := []string{}
	for _, folder := range *folders {
		for _, filter := range filters {
			if filter.MatchString(folder) {
				foldersFiltered = append(foldersFiltered, folder)
				break
			}
		}
	}

	for _, file := range *files {
		for _, filter := range filters {
			if filter.MatchString(file) {
				filesFiltered = append(filesFiltered, file)
				break
			}
		}
	}

	mirrors := []string{}
	for _, folder := range foldersFiltered {
		mirrors = append(mirrors, buildMirror(folder))
	}

	pgets := []string{}
	for _, file := range filesFiltered {
		pgets = append(pgets, buildPget(file))
	}

	lftpCommands := ""
	if len(mirrors) != 0 {
		lftpCommands = strings.Join(mirrors, "; ")
	}

	if len(pgets) != 0 {
		lftpCommands = strings.Join(
			[]string{lftpCommands, strings.Join(pgets, "; ")}, "; ")
	}

	cmd := exec.Command("lftp", "-c", lftpCommands)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	if err := cmd.Run(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
