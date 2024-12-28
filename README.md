# File Management and Network Application

## Overview
This application facilitates file transfer and management using a combination of Go and C++. It includes functionalities to upload and receive files, ensuring proper folder structures and network communication. The application also includes a web-based interface for user interaction.

---

## Features
1. **Folder Management**
   - Automatically ensures the existence of `Upload` and `Received` folders at startup. This feature was implemented with a Go function that creates these folders if they don't exist.
   - Creates missing folders as needed.

2. **File Transfer**
   - Sends files to specified peers using a C++ backend.
   - Receives files from peers and saves them in the `Received` folder.

3. **Web Interface**
   - Users can upload files via a web page.
   - View and download received files directly from the browser.

4. **Network Communication**
   - Handles file transfer requests and responses.
   - Allows checking the active status of peers.

---

## Setup Instructions

### Prerequisites
1. Install Go.
2. Install a C++ compiler (e.g., `g++`).

### Installation
1. Clone the repository:
   ```bash
   git clone <repository_url>
   cd <repository_folder>
   ```
2. Run the Go application:
   ```bash
   go run main.go
   ```

---

## Usage

### Folder Initialization
At startup, the application checks for the `Upload` and `Received` folders:
- If these folders do not exist, they are created automatically.

### Web Interface
1. Start the server:
   ```bash
   go run main.go
   ```
2. Open a browser and navigate to:
   ```
   http://localhost:8080/
   ```

3. Upload files through the interface or view/download received files.

### File Transfer
- Upload a file to the `Upload` folder via the web interface or directly.
- Use the `sender` binary to send files to a specified IP.

### File Reception
- Received files are stored in the `Received` folder and are accessible via the web interface.

---

## API Endpoints

### `GET /recieved`
Handles file recieved and displays file from the transfer.

### `POST /upload`
Handles file uploads and initiates file transfer.
- **Parameters:**
  - `file`: The file to upload.
  - `desiredIP`: The IP address of the recipient.

### `GET /scanPeers`
Scans for active peers.

### `GET /Active`
Returns the active status of the server.

---

## Error Handling
- Properly logs and handles errors during file transfer.
- Displays meaningful error messages for common issues (e.g., missing folders, network failures).

---

## Folder Permissions
- Ensure that the `Upload` and `Received` folders have appropriate read/write permissions for the application to function correctly.

## Contributions And Acknowledgements 
This project was developed with contributions across languages. The HTML file was primarily created to serve as the web interface, with significant input from ChatGPT to streamline its design and ensure integration with the backend. Additionally, the system includes C++ components for efficient file transfer, whose adjustments and implementations were also guided with assistance from ChatGPT. 



