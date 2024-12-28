package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"os/exec"
	"regexp"
	"strings"
	"time"
)

//#include "fileTransfer.cpp"

func main() {
	ipList := findingActive()
	for ip, name := range ipList {
		fmt.Printf("IP: %s, Active Name: %s\n", ip, name)
	}
}

// Func is used to get base ip: xxx.xxx.xxx.
func getBaseIp() string {
	//Establishing a remote connection to get local ip
	conn, err := net.Dial("udp", "8.8.8.8:80")
	if err != nil {
		log.Fatal(err)
	}
	defer conn.Close()

	//Pulling strictly the IP and not the port
	localAddr := conn.LocalAddr().(*net.UDPAddr).IP.String()
	parts := strings.Split(localAddr, ".")

	//Getting the base IP
	baseIp := fmt.Sprintf("%s.%s.%s.", parts[0], parts[1], parts[2])

	return baseIp
}

// Finding and adding active ips to a list
func findingActive() map[string]string {
	ip := getBaseIp()
	ipList := make(map[string]string)

	for i := 0; i <= 255; i++ {
		//Setting the url up: xxx.xxx.xxx.xxx:/8080/Active
		url := fmt.Sprintf("%s%d:8080/Active", ip, i)

		//Setting a response time or else timeout
		ctx, cancel := context.WithTimeout(context.Background(), 250*time.Millisecond)
		defer cancel()

		//Running and collecting the output
		cmd := exec.CommandContext(ctx, "curl", url)
		output, _ := cmd.CombinedOutput()

		name := getActiveName(string(output))

		//Checking to make sure there's a name
		if len(name) > 0 {
			localIP := fmt.Sprintf("%s%d", ip, i)
			ipList[localIP] = name
		}
	}
	return ipList
}

// Return only the "Active : <name>" section of the output
func getActiveName(output string) string {
	//Setting a regex pattern for "Active : <name>"
	re := regexp.MustCompile(`Active\s*:\s*(\w+)`)
	match := re.FindStringSubmatch(output)

	//Make sure it was the right output
	if len(match) > 1 {
		return match[1] //Returning the name
	}

	//Return an empty string if no match was found
	return ""

}
