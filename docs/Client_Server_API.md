# Client Server API v0.2
Client server communication will be done through TCP

## Virtual Messages
TCP is a connection based protocol with no support for messages. Therefor a method
of splitting the data from the connection into individual messages of different types
needs to be created.

##### Prefixing messages with message length
Each individual message sent to and from the server needs to be prefixed by a 2 Byte (16bit) unsigned
integer which describes how many bytes the following message contains. The server and client will then know
how many of the following bytes to make the message. once that many bytes has been read and a message is composed,
the following 2 bytes in the stream will be the next messages prefix (and so on).

## Message types
Each message has message type. i.e a "move" message may have a type of 1 which tells the server
to read the following message as a "move" message. **The first byte** of each message will be an
unsigned 8bit integer holding the message type.

## Strings in messages
In order for the server to be able to tell when a string ends each string placed in the payload
of a message needs to be null terminated i.e it's last byte needs to be **\0**

---

## Pre - Gameplay Messages
pre gameplay messages are all the messages that can be made before an actual game starts

### Request Alias
#### Client --> Server
#### Message type: 1
#### Description
This message is sent from a client to the server upon connecting. The purpose of the
message is to give the server a unique string (username/alias) representing what that client
wishes to call himself.
#### Payload
```
1 byte      : uint      : Message type
Var bytes   : String    : The alias the user is requesting
```

---

### Alias Approved
#### Server --> Client
#### Message type: 2
#### Description
This message is sent to the client when the alias he picked has been approved.
#### Payload
```
1 byte      : uint      : Message type
```

---

### Alias Denied
#### Server --> Client
#### Message type: 3
#### Description
This message is sent to the client when the alias he picked has been denied and he must pick another one.
#### Payload
```
1 byte      : uint      : Message type
var bytes   : string    : Reason for disproving
```

---

### Request Game List
#### Client --> Server
#### Message type: 4
#### Description
This message is sent to the server when the client wishes to update his open game list. The open game list is the list of opened pong games with one person in them waiting for someone else to join
#### Payload
```
1 byte      : uint      : Message type
```

---

### Game List
#### Server --> Client
#### Message type: 5
#### Description
This message is sent to a client to give it the updated list of open games. This may happen regularly or immediately after a client requests the game list.
#### Payload
```
1 byte      : uint      : Message type
2 bytes     : uint      : Number of games

// the following represents a game. it will be repeated
// depending on the number of open games.
4 bytes     : int       : Unix timestamp of when the game was created
var bytes   : string    : Game id, a unique identifier for a game.
var bytes   : string    : Alias of the user who created the game.
var bytes   : string    : Game name (game creator picks this)
```

---

### Create Game
#### Client --> Server
#### Message type: 6
#### Description
This message is sent when a client wishes to create a game lobby.
#### Payload
```
1 byte      : uint      : Message type
var bytes   : string    : Game name
```

---

### Create Game Approved
#### Server --> Client
#### Message type: 7
#### Description
Sent to the client who requested to create a game. Once the game is created the client cannot
join any other games until he specifically leaves this game.
#### Payload
```
1 byte      : uint      : Message type
var bytes   : string    : Game id
var bytes   : string    : Game name
```

---

### Create Game Denied
#### Server --> Client
#### Message type: 8
#### Description
Sent to the client who wanted to create a game.
#### Payload
```
1 byte      : uint      : Message type
var bytes   : string    : Game name
var bytes   : string    : Reason for denying game creation
```

---

### Join Game
#### Client --> Server
#### Message type: 9
#### Description
Sent to the server when a client wishes to join a game.
#### Payload
```
1 byte      : uint      : Message type
var bytes   : string    : Game id
```

---

### Leave Game
#### Client --> Server
#### Message type: 10
#### Description
Sent to the Server when a client wishes to leave a game.
#### Payload
```
1 byte      : uint      : Message type
```

---

## Gameplay Messages
A game is composed of 2 players and a ball. The game board is 500 units in height and 750 units in length. Each player has a height of 125 units and no width and moves only vertically. Each player is placed 30 units before the edge of their side of the game board. The ball is 25x25 units. As the ball moves if the player moves into it's way the ball bounces off the player and starts heading towards the other players side. if the player does not get in it's way and the ball moves off through their side of the board, the other player wins that round. Each game has 10 rounds.

player positions are relative to the bottom unit of the player. i.e if the player position is 0 :
he is at the bottom of the game board. and if the player position is 375 he is at the top of the game
board (375 + 125 = 500)

ball position Y is in the same manner as the player position.
ball position X is from the left of the screen i.e. ball position 0 is leftmost position and 725 is rightmost position.

---

### Start Game
#### Server --> Client
#### Message type: 11
#### Description
When two clients are in a lobby a game can start. The Server sends each client in the game
a start game message.
#### Payload
```
1 byte      : uint      : Message type
1 byte      : boolean   : your side of game board (1 = right, 0 = left)
2 bytes     : uint      : your position (relative to the bottom unit of the player)
2 bytes     : uint      : other player position (relative to the bottom unit of the player)
2 bytes     : uint      : ball x position (relative to left of ball)
2 bytes     : uint      : ball y position (relative to bottom of ball)
var bytes   : string    : Other player alias
var bytes   : string    : Game id
var bytes   : string    : Game name
```

---

### State Update
#### Server --> Client
#### Message type: 12
#### Description
Message sent from the server to the clients in a game 60 times a second telling them the status of the objects in the game.
#### Payload
```
1 byte      : uint      : Message type
2 bytes     : uint      : other player position (relative to the bottom unit of the player)
2 bytes     : uint      : ball x position (relative to left of ball)
2 bytes     : uint      : ball y position (relative to bottom of ball)
```

---

### Round Update
#### Server --> Client
#### Message type: 13
#### Description
Sent to the players when the round is over. The ball position is given in this one because the ball resets when the round is over. There may be a second of waiting time
before the ball begins moving again to give players time to recover.
#### Payload
```
1 byte      : uint      : Message type
1 bytes     : uint      : your score
1 bytes     : uint      : other player score
2 bytes     : uint      : ball x position (relative to left of ball)
2 bytes     : uint      : ball y position (relative to bottom of ball)
```

---

### Game Over
#### Server --> Client
#### Message type: 14
#### Description
Sent to the players when the game has finished.
#### Payload
```
1 byte      : uint      : Message type
1 bytes     : uint      : your score
1 bytes     : uint      : other player score
1 bytes     : uint      : game state (0 = your victory, 1 = other player's victory, 2 = draw)
```

---

### Move
#### Client --> Server
#### Message type: 15
#### Description
When a player is changing his position this message is called.
#### Payload
```
1 byte      : uint      : Message type
2 bytes     : uint      : The players new position
```

---
