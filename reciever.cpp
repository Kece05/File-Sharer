#include <iostream>
#include <fstream>
#include <string>
#include <sys/socket.h>
#include <netinet/in.h>
#include <unistd.h>

using namespace std;

int main() {
    //----Setting and defining the new socket----
    int server_fd, new_socket;
    struct sockaddr_in address;
    int opt = 1;
    int addrlen = sizeof(address);

    if ((server_fd = socket(AF_INET, SOCK_STREAM, 0)) == 0) {
        perror("socket failed");
        exit(EXIT_FAILURE);
    }

    if (setsockopt(server_fd, SOL_SOCKET, SO_REUSEADDR, &opt, sizeof(opt)) < 0) {
        perror("setsockopt");
        exit(EXIT_FAILURE);
    }

    address.sin_family = AF_INET;
    address.sin_addr.s_addr = INADDR_ANY;
    address.sin_port = htons(8081);
    //--------------------------------------------

    //------Binding and then listening on that socket-------
    if (bind(server_fd, (struct sockaddr *)&address, sizeof(address)) < 0) {
        perror("bind failed");
        exit(EXIT_FAILURE);
    }

    if (listen(server_fd, 3) < 0) {
        perror("listen");
        exit(EXIT_FAILURE);
    }
    //------------------------------------------------------

    //Continue to allow new connections
    while (true) {
        //Checking for connection
        if ((new_socket = accept(server_fd, (struct sockaddr *)&address, (socklen_t*)&addrlen)) < 0) {
            perror("accept");
            continue;
        }

        cout << "Connection accepted." << endl;

        //Recieving the file name
        char buffer[1024];
        int Fbytes_received;

        // Receive the file name first
        Fbytes_received = recv(new_socket, buffer, sizeof(buffer), 0);
        if (Fbytes_received > 0) {
            //Extracting the file name
            string dataRecieved(buffer, Fbytes_received);
            char delimiter = ' ';

            size_t pos = dataRecieved.find(delimiter);
            string file_name = (pos != string::npos) ? dataRecieved.substr(0, pos) : dataRecieved;

            //Opening the a new file
            ofstream outfile("Received/" + file_name, ios::binary);
            if (!outfile) {
                cerr << "Failed to open file for writing: " << file_name << endl;
                return -1;
            }

            //This writing is used if its a .txt file
            dataRecieved.erase(0, file_name.size()+1);
            outfile.write(dataRecieved.c_str(), dataRecieved.size());

            //This one is used for every other case
            int bytes_received;
            while ((bytes_received = recv(new_socket, buffer, sizeof(buffer), 0)) > 0) {
                //Writes the data to the file
                outfile.write(buffer, bytes_received); 
            }

            if (bytes_received < 0) {
                std::cerr << "Error receiving file data." << std::endl;
            } else {
                std::cout << "File received successfully!" << std::endl;
            }

            outfile.close();
        } else {
            std::cerr << "Error receiving file name." << std::endl;
        }


    }
    close(server_fd);
    return 0;
}
