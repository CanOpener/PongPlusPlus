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
#include <chrono>
using namespace std;

typedef unsigned char BYTE;

// constants
const int BOARD_LENGTH     = 750;
const int BOARD_HEIGHT     = 500;
const int PLAYER_HEIGHT    = 125;
const int PLAYER1_X        = 30;
const int PLAYER2_X        = 720;
const int BALL_LENGTH      = 25;
const int NUMBER_OF_ROUNDS = 10;

// structure containing the game state
struct gameState {
    uint8_t messageType;
    uint8_t current_Round;
    uint8_t player_1_score;
    uint8_t player_2_score;
    uint16_t player_1_position;
    uint16_t player_2_position;
    uint16_t ballX;
    uint16_t ballY;
};

// structure containing a server message
struct serverMessage {
    uint8_t type;
    bool player1;
    uint16_t newPosition;
};

// structure containing a finished message
// sent to the server when the game is about to terminate
struct finishedMessage {
    uint8_t messageType;
    bool p1Won;
    uint8_t rounds;
    uint8_t player_1_Score;
    uint8_t player_2_Score;
};

// union of server message struct and an array of bytes
union messageDecode {
    BYTE* bytes;
    serverMessage result;
};

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

    // send ready message
    uint8_t ready = 1;
    write(sockfd, &ready, sizeof(uint8_t));

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

bool ballPlayerAligned(int playerY, int ballY) {
    return (ballY > (playerY-BALL_LENGTH) && ballY < (playerY + PLAYER_HEIGHT));
}

bool ballPastPlayer(int ballx) {
    return (ballx < PLAYER1_X || ballx+BALL_LENGTH > PLAYER2_X);
}

// absolutely basic. will be upgraded when
// integration with main server is figured out
gameState updateBallPosition(gameState gs, int sockfd) {
    static auto ballXvelocity = 0;
    static auto before        = clock();
    static auto waitTime      = 5;

    auto now = clock();
    auto ugs = gs;

    if ((now - before) < waitTime * CLOCKS_PER_SEC) {
        return ugs;
    }
    before   = now;
    waitTime = 0;

    // new round
    if (ballXvelocity == 0) {
        ((ugs.player_1_score > ugs.player_1_score) ? ballXvelocity = -2 : ballXvelocity = 2);
    }

    if (ballPastPlayer(ugs.ballX + ballXvelocity)) {
        if (ballXvelocity < 0) {
            if (ballPlayerAligned(ugs.player_1_position, ugs.ballY)) {
                ballXvelocity *= -1;
            } else {
                ugs.current_Round++;
                ugs.player_2_score++;
                ugs.ballX = (BOARD_LENGTH/2) - (BALL_LENGTH/2);
                ballXvelocity = 0;
                waitTime = 3;
            }
        } else {
            if (ballPlayerAligned(ugs.player_2_position, ugs.ballY)) {
                ballXvelocity *= -1;
            } else {
                ugs.current_Round++;
                ugs.player_1_score++;
                ugs.ballX = (BOARD_LENGTH/2) - (BALL_LENGTH/2);
                ballXvelocity = 0;
                waitTime = 3;
            }
        }
    }

    ugs.ballX += ballXvelocity;

    if (ugs.current_Round >= 10) {
        finishedMessage fin;
        fin.messageType = 3;
        fin.rounds = ugs.current_Round;
        fin.p1Won = (ugs.player_1_score > ugs.player_2_score);
        fin.player_1_Score = ugs.player_1_score;
        fin.player_2_Score = ugs.player_2_score;
        write(sockfd, &fin, sizeof(finishedMessage));
        exit(1);
    }

    return ugs;
}

// updatePlayerPositions looks at the call queue and
// updates player positions
gameState updatePlayerPositions(gameState gs, int sockfd) {
    auto queue = getAll();
    auto ugs = gs;
    while (!queue.empty()) {
        auto data = queue.front();
        messageDecode srvMsg;
        srvMsg.bytes = data;
        if (srvMsg.result.type != 1) {
            // someone disconnected. automatically disqualified.
            finishedMessage fin;
            fin.messageType = 3;
            fin.p1Won = !srvMsg.result.player1; // if p1 left, p1 lost
            fin.rounds = ugs.current_Round;
            fin.player_1_Score = ugs.player_1_score;
            fin.player_2_Score = ugs.player_2_score;
            delete[] data;
            write(sockfd, &fin, sizeof(finishedMessage));
            exit(1);
        }
        if (srvMsg.result.newPosition <= (BOARD_HEIGHT - PLAYER_HEIGHT) && srvMsg.result.newPosition >= 0) {
            if (srvMsg.result.player1) {
                ugs.player_1_position = srvMsg.result.newPosition;
            } else {
                ugs.player_2_position = srvMsg.result.newPosition;
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
    gs.messageType       = 12;
    gs.current_Round     = 1;
    gs.player_1_score    = gs.player_2_score = 0;
    gs.player_1_position = (BOARD_HEIGHT/2) - (PLAYER_HEIGHT/2);
    gs.player_2_position = (BOARD_HEIGHT/2) - (PLAYER_HEIGHT/2);
    gs.ballX             = (BOARD_LENGTH/2) - (BALL_LENGTH/2);
    gs.ballY             = (BOARD_HEIGHT/2) - (BALL_LENGTH/2);

    // simple game loop for now
    // will optimise later
    auto begin = clock();
    while (true) {
        begin = clock();
        gs = updatePlayerPositions(gs, sockfd);
        gs = updateBallPosition(gs, sockfd);
        write(sockfd, &gs, sizeof(gameState));
        this_thread::sleep_for(chrono::milliseconds((clock() - begin) - (CLOCKS_PER_SEC / 1000)));
    }
    return 1;
}

int main(int argc, char const *argv[]) {
    if (argc != 3) {
        cout << "Usage: " << argv[0] << " <UDS Address> <Tick Rate>\n";
        return 1;
    }
    string UDSaddr = argv[1];
    int tickRate   = atoi(argv[2]);
    int sockfd     = connect(UDSaddr);

    thread lis(listener, sockfd);
    lis.detach();

    gameLoop(sockfd, tickRate);

    return 0;
}
