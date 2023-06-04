package main

import (
	"net"
	"strconv"
	"sync"
	"time"

	"github.com/go-ping/ping"
)

func doNetworkChecks() {
	wg := new(sync.WaitGroup)
	isFin := make(chan bool, config.NetworkChecksAllListSize)
	for protocol_i, nc := range config.NetworkChecks {
		for list_j := range nc.List {
			wg.Add(1)
			go networkCheckProcess(protocol_i, list_j, isFin, wg)
		}
	}
	wg.Wait()
	close(isFin)
}

func networkCheckProcess(i int, j int, isFin chan bool, wg *sync.WaitGroup) {
	defer wg.Done()

	config.NetworkChecks[i].List[j].IsBeforeUp = config.NetworkChecks[i].List[j].IsNowUp

	if config.NetworkChecks[i].Protocol == "ping" {
		pinger, err := ping.NewPinger(config.NetworkChecks[i].List[j].Address)
		pinger.SetPrivileged(true)
		if err != nil {
			mainLogger.Panicln(err)
		}
		pinger.Count = 2
		pinger.Timeout = time.Duration(time.Duration(5) * time.Second)
		pinger.Run()
		stats := pinger.Statistics()
		if stats.PacketsRecv >= 1 {
			config.NetworkChecks[i].List[j].IsNowUp = true
		} else {
			config.NetworkChecks[i].List[j].IsNowUp = false
		}
	}
	if config.NetworkChecks[i].Protocol == "tcp" {
		conn, err := net.DialTimeout("tcp", config.NetworkChecks[i].List[j].Address+":"+strconv.Itoa(config.NetworkChecks[i].List[j].Port), time.Duration(5)*time.Second)
		if err != nil {
			config.NetworkChecks[i].List[j].IsNowUp = false
			isFin <- true
			return
		}
		defer conn.Close()
		config.NetworkChecks[i].List[j].IsNowUp = true
	}
	isFin <- true
}
