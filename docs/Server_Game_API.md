# Server Game API v0.1
The Server Game API outlines the communication specification between the connection handling server (GO) and the game running binary(C++). There will be no sterilisation of data, instead it is a strict raw binary specification similar to that of the Client Server API.

## Technology
Interprocess communication will be hosted by a unix domain socket with the listener on the connection handling server in go. The go server will start each game process as needed with command line arguments.
* Socket Location - The path to the socket
* Tick Rate - The ticks per second the game is to run at
## Message types
Each message has message type. i.e a "move" message may have a type of 1 which tells the server
to read the following message as a "move" message. **The first byte** of each message will be an
unsigned 8bit integer holding the message type.

## Specification
### Ready
#### Game --> Server
#### Message type: 1
#### Description
This Message is sent to the server to inform it that all the game is ready for player input. This message means that the ball will start moving in 3 seconds
#### Payload
```
1 byte      : uint      : Message type
```
---

### Status
#### Game --> Server
#### Message type: 2
#### Description
This message contains the current status of the game. This message is to be passed to the players as it contains vital information about the position of the ball and the other player. The frequency of this message is reflected by the **Tick Rate** that the game has been started with. note that player 1 is the player on the left of the board
#### Payload
```
1 byte      : uint      : Message type
1 byte      : uint      : Current_Round
1 byte      : uint      : Player_1_Score
1 byte      : uint      : Player_2_Score
2 byte      : uint      : Player_1_Position
2 byte      : uint      : Player_2_Position
2 byte      : uint      : Ball_X_Coordinate
2 byte      : uint      : Ball_Y_Coordinate
```
---

### Finished
#### Game --> Server
#### Message type: 3
#### Description
This is the last message sent by the game to the server. It means that the game has finished and the process will now terminate
#### Payload
```
1 byte      : uint      : Message type
1 byte      : bool      : true = player_1 won, false = player_2 won,
1 byte      : uint      : Total Rounds Played
1 byte      : uint      : Player_1_Score
1 byte      : uint      : Player_2_Score
```
---

### Someone Disconnected
#### Server --> Game
#### Message type: 4
#### Description
This message tells the game that a player has left the game. This will cause the game to terminate after a **Finished** message.
#### Payload
```
1 byte      : uint      : Message type
1 byte      : bool      : true = player 1 left, false = player 2 ...
```
---
