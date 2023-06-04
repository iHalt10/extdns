package main

import (
	"log"
	"os"
	"os/exec"
	"time"

	flags "github.com/jessevdk/go-flags"
)

type Options struct {
	ConfigFile string `long:"config-file" default:"./config.yaml" description:"config file"`
	HostsFile  string `long:"hosts-file" default:"./hosts" description:"hosts file"`
}

var (
	opts       Options
	config     Config
	mainLogger *log.Logger
	// dnsmasqLogger *log.Logger
	dnsmasq = exec.Command("dnsmasq", "-d")
)

func main() {
	mainLogger = log.New(os.Stdout, "[main]: ", log.LstdFlags)
	_, err := flags.Parse(&opts)
	if err != nil {
		os.Exit(0)
	}

	initialize()
	mainLogger.Println("Start")
	mainLogger.Printf("config: %+v\n", config)
	run()
}

func initialize() {
	readConfig()
	if len(config.DnsmasqArgs) >= 0 {
		dnsmasq = exec.Command("dnsmasq", config.DnsmasqArgs...)
	}
	dnsmasq.Stdout = os.Stdout
	doNetworkChecks()
	writeHosts()
}

func run() {
	dnsmasq.Start()

	for {
		time.Sleep(3 * time.Second)
		doNetworkChecks()

		isFix := false
		for i, nc := range config.NetworkChecks {
			for j := range nc.List {
				if config.NetworkChecks[i].List[j].IsBeforeUp != config.NetworkChecks[i].List[j].IsNowUp {
					isFix = true
				}
			}
		}
		if isFix {
			writeHosts()
			mainLogger.Printf("Fix hosts: %v", config.NetworkChecks)

			dnsmasq.Process.Kill()
			dnsmasq = exec.Command("dnsmasq", config.DnsmasqArgs...)
			dnsmasq.Stdout = os.Stdout
			dnsmasq.Start()
		}
	}
}
