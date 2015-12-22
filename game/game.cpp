#include <iostream>
#include <mutex>
#include <deque>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
#include <netinet/in.h>
#include <netdb.h>
#include <stdio.h>
#include <string.h>
#include <stdlib.h>
#include <unistd.h>
#include <errno.h>
#include <arpa/inet.h>
using namespace std;

typedef unsigned char BYTE

// Thread safe queuing
mutex mu;
deque<BYTE*> callQueue;

void add(BYTE* bytes) {
    lock_guard<mutex> locker(mu);
    callQueue.push_back(bytes);
}
deque<BYTE*> getAll() {
    lock_guard<mutex> locker(mu);
    deque<BYTE*> ret = callQueue; // deep copy
    callQueue.clear();
    return ret;
}

// Listener
void listener(string UDSaddr) {
    int sockfd = 0, size = 0;
    const int buffSize = 1400;
    BYTE recvBuff[buffSize];
    struct sockaddr_un serv_addr;

    memset(recvBuff, '0' ,sizeof(recvBuff));
    if ((sockfd = socket(AF_UNIX, SOCK_, 0))< 0) {
        cout << "Could not create socket\n";
        exit(1);
    }

    memcpy(serv_addr.sun_path, UDSaddr.c_str(), UDSaddr.length());
    serv_addr.sun_family = AF_UNIX;
    serv_addr.sun_path[UDSaddr.length()] = '\0';

    if (connect(sockfd, (struct sockaddr *)&serv_addr, sizeof(serv_addr))<0) {
        cout << "Connection Failed\n";
        exit(1);
    }
    while ((n = read(sockfd, recvBuff, buffSize-1) > 0) {
        BYTE* data = new BYTE[n];
        memcpy(data, recvBuff, n);
        add(data);
    }
    exit(1);
}

int main() {
    cout << "Hello World!" << endl;
}
