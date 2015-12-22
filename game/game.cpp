#include <iostream>
#include <thread>
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
#include <time>
#include <arpa/inet.h>
using namespace std;

typedef unsigned char BYTE

const int BOARD_LENGTH     = 750;
const int BOARD_HEIGHT     = 500;
const int PLAYER_HEIGHT    = 125;
const int PLAYER1_X        = 30;
const int PLAYER2_X        = 720;
const int BALL_LENGTH      = 25;
const int NUMBER_OF_ROUNDS = 10;

// Thread safe queuing
mutex queMu;
deque<BYTE*> callQueue;

void add(BYTE* bytes) {
    lock_guard<mutex> locker(queMu);
    callQueue.push_back(bytes);
}
deque<BYTE*> getAll() {
    lock_guard<mutex> locker(queMu);
    deque<BYTE*> ret = callQueue; // deep copy
    callQueue.clear();
    return ret;
}

// listener listens for messages from the main server
void listener(int sockfd) {
    int size = 0;
    const int buffSize = 1400;
    BYTE recvBuff[buffSize];
    memset(recvBuff, '0' ,sizeof(recvBuff));

    while ((n = read(sockfd, recvBuff, buffSize-1) > 0) {
        BYTE* data = new BYTE[n];
        memcpy(data, recvBuff, n);
        add(data);
    }
    exit(1);
}

// connect connects to the Unix domain socket and returns the sockedfd
int connect(string USDAddr) {
    int sockfd = 0;
    struct sockaddr_un serv_addr;

    if ((sockfd = socket(AF_UNIX, SOCK_STREAM, 0))< 0) {
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

    return sockfd;
}


// gameLoop runs the game
int gameLoop(int sockfd, int tickRate) {
    if (tickRate < 0) {
        cout << "Tick rate must be above 0.\n";
        exit(1);
    }

    // simple game loop for now
    // will optimise
    while(true) {
        // TODO: update player positions
        // TODO: update ball position
        // TODO: send gamestate to server
    }
    return 1;
}

int main(int argc, char const *argv[]) {
    if (argc != 3) {
        cout << "Usage: " << argv[0] << " <UDS Address> <Tick Rate>\n";
        return 1;
    }
    string USDAddr = argv[1];
    int tickRate   = atoi(argv[2]);
    int sockfd     = connect(USDAddr);

    thread lis(listener, sockfd);
    lis.detach();

    auto begin = clock();
    int gameResult = gameLoop(sockfd);
    auto seconds = int(double(clock() - begin) / CLOCKS_PER_SEC);

    // return results to server
    return 0;
}
