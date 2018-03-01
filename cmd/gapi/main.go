package main

import (
	"crypto/tls"
	"encoding/json"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"

	gapi "github.com/AutogrowSystems/go-grafana-api"
	"github.com/ghodss/yaml"
)

var defaultConfigPath = fmt.Sprintf("%s/.grafana/api.yml", os.Getenv("HOME"))

type config struct {
	Auth       string `json:"auth"`
	URL        string `json:"url"`
	SkipVerify bool   `json:"skip_verify"`
}

func loadConfig(path string) (config, error) {
	cfg := config{}
	b, err := ioutil.ReadFile(path)
	if err != nil {
		return cfg, err
	}
	err = yaml.Unmarshal(b, &cfg)
	if err != nil {
		return cfg, err
	}

	return cfg, err
}

func main() {

	var cfgFile string
	var create, read, find, list bool //, read, update, delete bool
	var datasource, dashboard, org bool
	var thingName string
	var thingID int64

	flag.StringVar(&cfgFile, "c", defaultConfigPath, "config file")

	flag.BoolVar(&create, "create", false, "create a thing")
	flag.BoolVar(&read, "read", false, "read a thing")
	flag.BoolVar(&find, "find", false, "find a thing")
	flag.BoolVar(&list, "list", false, "list a thing")

	flag.BoolVar(&datasource, "datasource", false, "work with datasource")
	flag.BoolVar(&dashboard, "dashboard", false, "work with dashboard")
	flag.BoolVar(&org, "org", false, "work with org")

	flag.StringVar(&thingName, "name", "", "name to find by")
	flag.Int64Var(&thingID, "id", 0, "id to find by")

	flag.Parse()

	log.SetOutput(os.Stderr)

	cfg, err := loadConfig(cfgFile)
	panicIf(err)

	client, err := gapi.New(cfg.Auth, cfg.URL)
	panicIf(err)

	http.DefaultTransport.(*http.Transport).TLSClientConfig = &tls.Config{InsecureSkipVerify: cfg.SkipVerify}

	switch {
	case datasource:
		switch {
		case create:
			data, err := ioutil.ReadAll(os.Stdin)
			panicIf(err)
			ds := gapi.DataSource{}
			err = json.Unmarshal(data, &ds)
			panicIf(err)
			id, err := client.NewDataSource(&ds)
			panicIf(err)
			log.Println("created datasource with ID", id)
			fmt.Println(id)
		}

	case org:
		switch {
		case create:
			org, err := client.NewOrg(thingName)
			panicIf(err)
			log.Println("created new org with ID", org.Id)
			fmt.Println(org.Id)

		case list:
			orgs, err := client.Orgs()
			panicIf(err)
			for _, o := range orgs {
				fmt.Printf("%-6d %s\n", o.Id, o.Name)
			}

		case find:
			if thingName == "" && thingID == 0 {
				log.Println("ERROR: must specify name or id to find by")
				os.Exit(1)
			}

			orgs, err := client.Orgs()
			panicIf(err)

			var data []byte
			for _, o := range orgs {
				if o.Name == thingName || o.Id == thingID {
					data, err = json.MarshalIndent(o, "", "  ")
				}
			}

			panicIf(err)
			if len(data) == 0 {
				log.Println("ERROR: not found")
				os.Exit(1)
			}

			fmt.Println(string(data))
		}

	default:
		log.Println("ERROR: not implemented")
		os.Exit(1)
	}
}

func panicIf(err error) {
	if err != nil {
		panic(err)
	}
}
