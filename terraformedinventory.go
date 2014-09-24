package main

import (
	"encoding/json"
	"flag"
	"fmt"
	"github.com/hashicorp/terraform/terraform"
	"os"
	"path/filepath"
)

var list = flag.Bool("list", false, "list mode")
var host = flag.String("host", "", "host mode")

type Host struct {
	Name string `json:"name"`
	IP   string `json:"ip"`
	vars HostVars
}

type HostVars map[string]interface{}

type Group struct {
	name  string
	hosts []*Host
}

func main() {
	flag.Parse()
	file := flag.Arg(0)

	if file == "" {
		fmt.Printf("Usage: %s [options] path\n", os.Args[0])
		os.Exit(1)
	}

	if !*list && *host == "" {
		fmt.Println("Either --host or --list must be specified")
		os.Exit(1)
	}

	path, err := filepath.Abs(file)
	if err != nil {
		fmt.Printf("Invalid file: %s\n", err)
		os.Exit(1)
	}

	f, err := os.Open(path)
	if err != nil {
		fmt.Printf("Error opening tfstate file: %s\n", err)
		os.Exit(1)
	}

	state, err := terraform.ReadState(f)
	if err != nil {
		fmt.Printf("Error reading state: %s\n", err)
		os.Exit(1)
	}

	if *list {
		hosts := mustGetList(state)
		os.Stdout.Write(mustMarshal(hosts))

	} else if *host != "" {
		host := mustGetHost(state, *host)
		os.Stdout.Write(mustMarshal(host.vars))
	}
}

func mustGetList(state *terraform.State) interface{} {
	hosts, err := getList(state)

	if err != nil {
		fmt.Printf("Error getting host list: %s\n", err)
		os.Exit(1)
	}

	return hosts
}

func mustGetHost(state *terraform.State, hostname string) *Host {
	host, err := getHost(state, hostname)

	if err != nil {
		fmt.Printf("Error getting host variables: %s\n", err)
		os.Exit(1)
	}

	return host
}

func mustMarshal(whatever interface{}) []byte {
	json, err := json.Marshal(whatever)

	if err != nil {
		fmt.Printf("Error encoding JSON: %s\n", err)
		os.Exit(1)
	}

	return json
}

func getList(state *terraform.State) (interface{}, error) {
	hostnames := make([]string, 0)
	for _, h := range readHosts(state) {
		hostnames = append(hostnames, h.IP)
	}

	groups := map[string][]string{"production": hostnames}
	return groups, nil
}

func getHost(state *terraform.State, hostname string) (*Host, error) {
	for _, h := range readHosts(state) {
		if h.IP == hostname {
			return h, nil
		}
	}

	return nil, fmt.Errorf("No such host: %s", hostname)
}

func readHosts(state *terraform.State) []*Host {
	hosts := make([]*Host, 0)

	for _, resource := range state.Resources {
		if resource.Type == "aws_instance" {
			host := &Host{
				Name: resource.Attributes["id"],
				IP:   resource.Attributes["public_ip"],
				vars: HostVars{},
			}

			hosts = append(hosts, host)
		}
	}

	return hosts
}
