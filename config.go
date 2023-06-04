package main

import (
	"os"
	"path/filepath"

	"github.com/go-yaml/yaml"
)

type SingleOrMulti struct {
	Values []string
}

type Config struct {
	Hosts []struct {
		Domain                       string        `yaml:"domain"`
		ShortName                    string        `yaml:"shortName"`
		Address                      SingleOrMulti `yaml:"address"`
		IsNetworkCheck               bool          `yaml:"isNetworkCheck"`
		PriorityNetworkCheckProtocol string        `yaml:"priorityNetworkCheckProtocol"`
	} `yaml:"hosts"`
	NetworkChecks []struct {
		Protocol string `yaml:"protocol"`
		List     []struct {
			Address    string `yaml:"address"`
			Port       int    `yaml:"port"`
			IsNowUp    bool
			IsBeforeUp bool
		} `yaml:"list"`
	} `yaml:"networkChecks"`
	DnsmasqArgs              []string `yaml:"dnsmasqArgs"`
	NetworkChecksAllListSize int
}

func (sm *SingleOrMulti) UnmarshalYAML(unmarshal func(interface{}) error) error {
	var multi []string
	err := unmarshal(&multi)
	if err != nil {
		var single string
		err := unmarshal(&single)
		if err != nil {
			mainLogger.Fatalln(err)
		}
		sm.Values = make([]string, 1)
		sm.Values[0] = single
	} else {
		sm.Values = multi
	}
	return nil
}

func readConfig() {
	filename, _ := filepath.Abs(opts.ConfigFile)
	yamlFile, err := os.ReadFile(filename)
	if err != nil {
		mainLogger.Fatalln(err)
	}
	err = yaml.Unmarshal(yamlFile, &config)
	if err != nil {
		mainLogger.Fatalln(err)
	}

	if len(config.NetworkChecks) >= 2 && config.NetworkChecks[0].Protocol != "ping" {
		config.NetworkChecks[0], config.NetworkChecks[1] = config.NetworkChecks[1], config.NetworkChecks[0]
	}
	config.NetworkChecksAllListSize = 0
	for _, lc := range config.NetworkChecks {
		config.NetworkChecksAllListSize += len(lc.List)
	}
}
