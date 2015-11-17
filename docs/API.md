# API
Client server communication will be done through TCP

## Virtual Messages
TCP is a connection based protocol with no support for messages. Therefor a method
of splitting the data from the connection into individual messages of different types
needs to be created.

### Prefixing messages with message length
Each individual message sent to and from the server needs to be prefixed by a 2 Byte (16bit) unsigned
integer which describes how many bytes the following message contains. The server and client will then know
how many of the following bytes to make the message. once that many bytes has been read and a message is composed :
the following 2 bytes will be the next messages prefix (and so on).

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
2 bytes     : uint      : Message type
Var bytes   : String    : The alias the user is requesting
```

---

### Approve Alias
#### Server --> Client
#### Message type: 2
#### Description
This message is sent to the client when the alias he picked has been approved.
#### Payload
```
2 bytes     : uint      : Message type
```

---

### Denied Alias
#### Server --> Client
#### Message type: 3
#### Description
This message is sent to the client when the alias he picked has been denied and he must pick another one.
#### Payload
```
2 bytes     : uint      : Message type
var bytes   : string    : Reason for disproving
```

---

### Request game list
#### Client --> Server
#### Message type: 4
#### Description
This message is sent to the server when the client wishes to update his open game list. The open game list is the list of opened pong games with one person in them waiting for someone else to join
#### Payload
```
2 bytes     : uint      : Message type
```

---

### Game list
#### Server --> Client
#### Message type: 5
#### Description
This message is sent to a client to give it the updated list of open games. This may happen regularly or immediately after a client requests the game list.
#### Payload
```
2 bytes     : uint      : Message type
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
2 bytes     : uint      : Message type
var bytes   : string    : Game name
```

---

### Create Game approved
#### Server --> Client
#### Message type: 7
#### Description
Sent to the client who requested to create a game. Once the game is created the client cannot
join any other games until he specifically leaves this game.
#### Payload
```
2 bytes     : uint      : Message type
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
2 bytes     : uint      : Message type
var bytes   : string    : Game name
var bytes   : string    : Reason for denying game creation
```

---
