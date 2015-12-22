#include <iostream>
#include <thread>
#include <mutex>
#include <deque>
#include <string>
#include <sys/socket.h>
#include <sys/types.h>
#include <sys/un.h>
#include <netdb.h>
#include <stdio.h>
#include <stdlib.h>
#include <unistd.h>
//#include <errno.h>
//#include <arpa/inet.h>
using namespace std;

typedef unsigned char BYTE;

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
    memset(recvBuff, '0' ,buffSize);

    while ((size = read(sockfd, recvBuff, buffSize-1)) > 0) {
        BYTE* data = new BYTE[size];
        memcpy(data, recvBuff, size);
        add(data);
    }
    exit(1);
}

// connect connects to the Unix domain socket and returns the sockedfd
int connect(string UDSaddr) {
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

// structure containing the game state
struct gameState {
    int round;
    int player1_position;
    int player2_position;
    int ballX;
    int ballY;
};

// structure containing a server message
struct serverMessage {
    int type;
    bool player1;
    int newPosition;
};

// union of server message struct and an array of bytes
union messageDecode {
    BYTE* bytes;
    serverMessage result;
};

// updatePlayerPositions looks at the call queue and
// updates player positions
gameState updatePlayerPositions(gameState gs) {
    auto queue = getAll();
    auto ugs = gs;
    while (!queue.empty()) {
        auto data = queue.front();
        messageDecode srvMsg;
        srvMsg.bytes = data;
        if (srvMsg.result.type != 1) {
            cout << "received " << srvMsg.result.type << " type message from server\n";
            exit(1);
        }
        if (srvMsg.result.newPosition <= (BOARD_HEIGHT - PLAYER_HEIGHT) && srvMsg.result.newPosition >= 0) {
            if (srvMsg.result.player1) {
                ugs.player1_position = srvMsg.result.newPosition;
            } else {
                ugs.player2_position = srvMsg.result.newPosition;
            }
        }
        delete[] data;
        queue.pop_front();
    }
    return ugs;
}

// gameLoop runs the game
int gameLoop(int sockfd, int tickRate) {
    if (tickRate < 0) {
        cout << "Tick rate must be above 0.\n";
        exit(1);
    }

    gameState gs;
    gs.round = 1;
    gs.player1_position = gs.player2_position =
        (BOARD_HEIGHT/2) - (PLAYER_HEIGHT/2);
    gs.ballX = (BOARD_LENGTH/2) - (BALL_LENGTH/2);
    gs.ballY = (BOARD_HEIGHT/2) - (BALL_LENGTH/2);

    // simple game loop for now
    // will optimise
    while(true) {
        gs = updatePlayerPositions(gs);
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
    int gameResult = gameLoop(sockfd, tickRate);
    auto seconds = int(double(clock() - begin) / CLOCKS_PER_SEC);

    // return results to server
    return 0;
}
