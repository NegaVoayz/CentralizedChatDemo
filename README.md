# Centralized Chat Demo
A simplest chat room with literally no synchronization.

## Message Standard

### Package Design

A typical message package in the project is shown like below.

In this project, the server has id 0. The broadcast id is also 0 as well.

**Head:**

| Field | type | size |
| :--: | :--: | :--: |
| Timestamp | uint64 | 8B |
| Sender_id | uint32 | 4B |
| Receiver_id | uint32 | 4B |
| Message_type | uint32 | 4B |
| Message_size | uint32 | 4B |

**Body**:

|    Field     |  type  |         size         |
| :----------: | :----: | :------------------: |
| Message_body | []byte | **<Message_size>** B |

### Message Types

|      Name      |  Id  |                     Description                     |
| :------------: | :--: | :-------------------------------------------------: |
| Normal_Message |  0   |   Normal message type, contains a single message    |
|  Join_Message  |  1   | Sent when joining a chat room, contains id and name |
| Member_Message |  2   | Sent to inform members of other members' name & id  |
| Leave_Message  | 1023 |        Sent to say goodbye to other members         |

## Server & Client Communication Example

```mermaid
sequenceDiagram
	participant Server
	participant Ack
	participant Syn
	participant Fin
	Ack -->> Server: name:Ack,id:1061065
	Server ->> Ack: id:1061065
	Server ->> Ack: User Online: name:Ack,id:1061065
	Syn -->> Server: name:Syn,id:1354184
	Server ->> Syn: id:1354184
	Server ->> Syn: User Online: name:Ack,id:1061065
	Server ->> Syn: User Online: name:Syn,id:1354184
	Fin -->> Server: name:Fin,id:1142475
	Server ->> Fin: id:1142475
	Server ->> Fin: User Online: name:Ack,id:1061065
	Server ->> Fin: User Online: name:Syn,id:1354184
	Server ->> Fin: User Online: name:Fin,id:1142475
	Syn -->> Server: Broadcast: hi
	Server ->> Ack: Syn: hi
	Server ->> Fin: Syn: hi
	Ack -->> Server: Ack to Syn: hi syn
	Server ->> Syn: Ack: hi syn
	destroy Fin
	Server --x Fin: Found Fin(id: 1143475 ) Offline
	Server ->> Ack: User Offline: id:1142475
	Server ->> Syn: User Offline: id:1142475
	
```

## Build

```powershell
#server
go build srv.go
#client
go build clt.go
```
