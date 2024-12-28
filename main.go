package main

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"os"
	"os/exec"
	"runtime"
	"strings"
)

var name string

func main() {
	//Setting up the localhost
	initalization()

	fmt.Println("----- Starting Server At http://localhost:8080/ ----")

	//Running the starting page
	http.Handle("/", http.FileServer(http.Dir("./")))

	//When requested it will run /Backend/peerScan.go
	http.HandleFunc("/scanPeers", sendRequest)
	http.HandleFunc("/Active", response)
	http.HandleFunc("/upload", transferFile)
	http.HandleFunc("/recieved", recievedFiles)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

// This function initalizes all the required steps for the main to run
func initalization() {
	os := runtime.GOOS

	fmt.Print("Enter Name to be Use for other Users: ")
	fmt.Scan(&name)

	//Deteecting type of device
	if os == "windows" {
		killPortWindows("8080")
		killPortWindows("8081")
	} else {
		killPortMac("8080")
		killPortMac("8081")
	}

	//Setting Up The CPP files so it can ran later
	sendCommand := exec.Command("g++", "-o", "sender", "sender.cpp")
	_, errSend := sendCommand.CombinedOutput()
	if errSend != nil {
		fmt.Println("Error Starting File(sender): " + errSend.Error())
	}

	recevieCommand := exec.Command("g++", "-o", "reciever", "reciever.cpp")
	_, errRecieve := recevieCommand.CombinedOutput()
	if errRecieve != nil {
		fmt.Println("Error Starting File(Send): " + errRecieve.Error())
	}

	ensureFolders()

	//Run receiver.cpp in the background
	runReceiverInBackground()
}

// Checks for folders and if not there then adds them
func ensureFolders() {
	folders := []string{"Upload", "Received"}

	for _, folder := range folders {
		if _, err := os.Stat(folder); os.IsNotExist(err) {
			//Checks for file, if not there create it
			err := os.Mkdir(folder, 0755)

			if err != nil {
				fmt.Printf("Failed to create folder '%s': %s\n", folder, err)
			} else {
				fmt.Printf("Folder '%s' created successfully.\n", folder)
			}

		}
	}
}

// Sends a response to any computer that sends a request
func response(w http.ResponseWriter, r *http.Request) {
	//Sending Response
	w.Write([]byte("Active : " + name))
}

// Killing a specified on Mac
func killPortMac(port string) {

	//Killing port from previous use
	cmd := exec.Command("lsof", "-i", ":"+port)
	output, err := cmd.CombinedOutput()

	//Will only run if there is an active port
	if err == nil {
		//Pulling the PID
		outputStr := string(output)
		lines := strings.Split(outputStr, "\n")

		var line int

		//There was a constant error happing with lsof
		//so this is to make sure it always will have the
		//right line number
		for i := 0; i < len(lines); i++ {
			fields := strings.Fields(lines[i])

			if len(fields) > 0 && fields[0] == "COMMAND" {
				line = i + 1
				break
			}
		}
		pid := strings.Fields(lines[line])[1]
		//Killing the port with the PID
		kill := exec.Command("kill", "-9", pid)
		_, killErr := kill.CombinedOutput()

		//Return wether sucessful or not
		if killErr != nil {
			fmt.Printf("Failed to kill process on port %s: %s\n", port, killErr)
		} else {
			fmt.Printf("Successfully killed process on port %s (PID: %s).\n", port, pid)
		}
	}
}

// Killing a specified port on Windows
func killPortWindows(port string) {
	//Pulling all active ports
	cmdNetstat := exec.Command("netstat", "-ano")
	output, _ := cmdNetstat.CombinedOutput()
	fullOutput := string(output)

	//Finding only port 8080
	cmdFindstr := exec.Command("findstr", ":"+port)
	cmdFindstr.Stdin = strings.NewReader(fullOutput)
	finalOutput, err := cmdFindstr.CombinedOutput()

	//Will only run if there is an active port
	if err == nil {
		//Pulling the PID
		outputStr := string(finalOutput)
		lines := strings.Split(outputStr, "\n")
		fields := strings.Fields(lines[1])
		pid := fields[4]

		fmt.Println(pid)

		//Killing the port with the PID
		kill := exec.Command("taskkill", "/PID", pid, "/F")
		_, killErr := kill.CombinedOutput()

		//Return wether sucessful or not
		if killErr != nil {
			fmt.Printf("Failed to kill process on port %s: %s\n", port, killErr)
		} else {
			fmt.Printf("Successfully killed process on port %s (PID: %s).\n", port, pid)
		}
	}
}

func sendRequest(w http.ResponseWriter, r *http.Request) {
	cmd := exec.Command("go", "run", "peerScan.go")
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error executing peerScan.go:", err)
		return
	}

	w.Header().Set("Content-Type", "text/plain")
	w.Write(output)
}

// Running the receiver file in the background
func runReceiverInBackground() {
	cmd := exec.Command("./reciever")

	//Runs it in the background
	cmd.Start()

	//Displays recievers output
	output, err := cmd.CombinedOutput()
	if err != nil {
		fmt.Println("Error running receiver:", err)
	} else {
		fmt.Println("Receiver output:", string(output))
	}
}

// This function pulls the data send by the POST request
// From the HTML file, and then saves the data, which
// the cpp file will transfer it over to the desired ip
// -Side note: I know I could directly transfer the file from
// This function but I wanted to learn how to network with cpp
func transferFile(w http.ResponseWriter, r *http.Request) {
	fmt.Printf("Request received: Method=%s, URL=%s\n", r.Method, r.URL.Path)
	w.Header().Set("Access-Control-Allow-Origin", "*")
	if r.Method != http.MethodPost {
		http.Error(w, "Invalid request method", http.StatusMethodNotAllowed)
		return
	}

	//Parse multipart form
	err := r.ParseMultipartForm(1024 * 1024 * 1024) //1024 MB limit
	if err != nil {
		http.Error(w, "Error parsing form", http.StatusInternalServerError)
		return
	}

	//Pulling the form data
	file, header, err := r.FormFile("file")
	if err != nil {
		http.Error(w, "Error retrieving the file", http.StatusBadRequest)
		return
	}
	defer file.Close()

	//------Save the file locally------
	dst, err := os.Create("Upload/" + header.Filename)
	if err != nil {
		http.Error(w, "Error saving the file", http.StatusInternalServerError)
		return
	}
	defer dst.Close()

	_, err = io.Copy(dst, file)
	if err != nil {
		http.Error(w, "Error writing the file", http.StatusInternalServerError)
		return
	}
	//----------------------------------

	//Getting the desired ip
	desiredIP := r.FormValue("desiredIP")
	fmt.Printf("File uploaded: %s, Desired IP: %s\n", header.Filename, desiredIP)

	cmdCompile := exec.Command("./sender", desiredIP, header.Filename)
	output, _ := cmdCompile.CombinedOutput()

	//Responding to the client
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(output))
}

// This function pulls any files from the Received folder
func recievedFiles(w http.ResponseWriter, r *http.Request) {
	folder := "Received/"

	//Opening the directory
	files, err := os.ReadDir(folder)
	if err != nil {
		http.Error(w, "Could not read the directory", http.StatusInternalServerError)
		return
	}

	//-----Adding each file to the fileNames array--------
	var fileNames []string
	for _, file := range files {
		fileNames = append(fileNames, file.Name())
	}
	//----------------------------------------------------

	//Converting and then send the list of the files back to the html file
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(fileNames)
}
