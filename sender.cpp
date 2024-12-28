#include <iostream>
#include <fstream>
#include <string>
#include <cstdio>
#include <sys/socket.h>
#include <arpa/inet.h>
#include <unistd.h>

using namespace std;

int main(int argc, char* argv[]) {
    //Getting the ip
    const string desiredIp = argv[1]; 
    const string fileName = argv[2];


    //-----Setting up the socket-----
    int sock = 0;
    struct sockaddr_in serv_addr;

    if ((sock = socket(AF_INET, SOCK_STREAM, 0)) < 0) {
        perror("Socket creation error");
        return -1;
    }
    //-------------------------------

    //------Setting up a Connection to target IP------
    serv_addr.sin_family = AF_INET;
    serv_addr.sin_port = htons(8081);

    if (inet_pton(AF_INET, desiredIp.c_str(), &serv_addr.sin_addr) <= 0) {
        perror("Invalid address/ Address not supported");
        return -1;
    }

    if (connect(sock, (struct sockaddr *)&serv_addr, sizeof(serv_addr)) < 0) {
        perror("Connection failed");
        return -1;
    }
    //------Setting up a Connection to target IP------

    //Sending the file name
    string update = fileName + " ";
    send(sock, update.c_str(), fileName.size()+1, 0);

    //Pulling what is wanted to be sent
    const string file = "Upload/" + fileName;
    ifstream infile(file, ios::binary);

    //Sending file
    char buffer[1024];
    while (infile.read(buffer, sizeof(buffer))) {
        send(sock, buffer, infile.gcount(), 0);
    }
    //Validates to make sure that all the data is sent over
    if (infile.gcount() > 0) {
        send(sock, buffer, infile.gcount(), 0);
    }

    //Closing and deleting file and socket
    infile.close();
    
    remove(file.c_str());

    return 0;
}
