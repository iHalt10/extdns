package main

import (
	"fmt"
	"os"
)

func writeHosts() {
	file, err := os.Create(opts.HostsFile)
	if err != nil {
		mainLogger.Fatalln(err)
	}
	defer file.Close()

	for _, host := range config.Hosts {
		for _, addr := range host.Address.Values {
			isWrite := true
			if host.IsNetworkCheck {
				for _, i := range getCheckSequence(host.PriorityNetworkCheckProtocol) {
					for j := range config.NetworkChecks[i].List {
						if addr == config.NetworkChecks[i].List[j].Address {
							if !config.NetworkChecks[i].List[j].IsNowUp {
								isWrite = false
							}
							goto LoopExit
						}
					}
				}
			LoopExit:
			}
			if isWrite {
				output := fmt.Sprintf("%s %s", addr, host.Domain)
				if host.ShortName != "" {
					output = fmt.Sprintf("%s %s", output, host.ShortName)
				}
				output += "\n"
				file.Write(([]byte)(output))
			}
		}
	}
}

func getCheckSequence(priorityNetworkCheckProtocol string) []int {
	sequence := []int{0, 1}
	if len(config.NetworkChecks) < 1 {
		return []int{}
	}
	if len(config.NetworkChecks) == 1 {
		return []int{0}
	}
	if priorityNetworkCheckProtocol == "tcp" {
		sequence = []int{1, 0}
	}
	return sequence
}
