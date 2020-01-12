package main

import (
	"encoding/json"
	"io/ioutil"
	"log"
	"os/user"
	"strings"

	aw "github.com/deanishe/awgo"
)

var (
	iconAvailable = &aw.Icon{Value: "./icon.png"}
	wf            *aw.Workflow // Our Workflow struct
)

const (
	ProjectConfFile = "/Library/ApplicationSupport/Code/User/projects.json"
)

type ProjectConf struct {
	Name     string   `json:"name"`
	RootPath string   `json:"rootPath"`
	Paths    []string `json:"paths"`
	Group    string   `json:"group"`
	Enabled  bool     `json:"enabled"`
}

func init() {
	// Create a new Workflow using default settings.
	// Critical settings are provided by Alfred via environment variables,
	// so this *will* die in flames if not run in an Alfred-like environment.
	wf = aw.New()
}

func main() {
	var query string

	if args := wf.Args(); len(args) > 0 {
		//	// If you're using "{query}" or "$1" (with quotes) in your
		//	// Script Fdilter, $1 will always be set, even if to an empty
		//	// string.
		//	// This guard serves mostly to prevent errors when run on
		//	// the command line.
		query = args[0]
	}
	log.Printf("query=%s", query)

	cUser, err := user.Current()
	if nil != err {
		log.Fatal(err)
	}
	ConfBytes, err := ioutil.ReadFile(cUser.HomeDir + ProjectConfFile)
	if err != nil {
		log.Fatal(err)
	}
	ConfString := strings.ReplaceAll(string(ConfBytes), "$home", cUser.HomeDir)
	pc := make([]*ProjectConf, 0)
	err = json.Unmarshal([]byte(ConfString), &pc)
	if err != nil {
		panic(err)
	}
	for _, entry := range pc {
		log.Printf("---- entry is %v \n", entry)
		if query != "" {
			if strings.Contains(entry.Name, query) {
				wf.NewItem(entry.Name).Icon(iconAvailable).Valid(true).Arg(entry.RootPath).Subtitle(entry.RootPath)
			}
		} else {
			wf.NewItem(entry.Name).Icon(iconAvailable).Valid(true).Arg(entry.RootPath).Subtitle(entry.RootPath)
		}
		log.Printf("---- entry is %v \n", entry)
	}

	wf.WarnEmpty("No matching repo found", "Try a different query?")

	wf.SendFeedback()
}
