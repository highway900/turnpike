package turnpike

import (
    "fmt"
    "container/list"
)

const (
    NO_CLIENT = "None"
)


type Room struct {
    masterToken string
    masterID string
    // members["token"] = clientID
    tokensMap map[string]string
    authMap map[string]bool
}


// Check if clientID is the master
//
func (r *Room) IsMaster(clientID string) bool {
    fmt.Printf("IsMaster Master: %s Client: %s\n", (*r).masterID, clientID)
    if r.masterID == clientID {
        return true
    }
    return false
}

// Checks if the room has a client
//
func (r *Room) HasClient(clientID string) bool {
    ok, _ := r.authMap[clientID]
    return ok
}


// Adds a security token for the member
// this will primarily come from a master client
//
func (r *Room) AddToken(token string) {
    r.tokensMap[token] = NO_CLIENT
}


// Connect a clientID to the security token
// If the token position is already filled return false
//
func (r *Room) AddClient(token string, clientID string) bool {
    // Check if the user is the master
    if token == r.masterToken {
        r.masterID = clientID
        fmt.Printf("AddClient Master clientID: %s\n", r.masterID)
        return true
    } else if r.tokensMap[token] == NO_CLIENT {
        r.tokensMap[token] = clientID
        r.authMap[clientID] = true
        fmt.Printf("AddClient clientID: %s\n", clientID)
        return true
    } else {
        return false
    }
}

func (r *Room) RemoveClientByID(clientID string) {
    for k, v := range r.tokensMap {
        if v == clientID {
            r.RemoveClient(k)
        }
    }
}

// Removes a client from the room
//
func (r *Room) RemoveClient(token string) bool {  
    id, ok := r.tokensMap[token]
    if !ok {
        return false
    }
    r.authMap[id] = false
    r.tokensMap[token] = NO_CLIENT
    return true
}

// Check if the client ID has been authenticated
//
//
func (r *Room) Auth(clientID string) bool {
    if r.IsMaster(clientID) {
        return true
    }
    return r.authMap[clientID]
}


// Creates a empty room
// with no tokensMap requires masterToken to be set.
//
func newRoom(token string) *Room {
    return &Room {
        masterToken: token,
        tokensMap: make(map[string]string),
        authMap: make(map[string]bool)}
}


type RoomsMap map[string]*Room


// Adds a room to the map
// requires a masterToken to authenticate who is in charge
func (rm RoomsMap) Add(id string, masterToken string) *Room {
    r := newRoom(masterToken)
    rm[id] = r
    return r
}


// Check if a room is in the rooms map
//
func (rm RoomsMap) Contains(roomID string) (*Room, bool) {
    room, ok := rm[roomID]
    return room, ok
}


// Remove a room from the rooms map
//
func (rm RoomsMap) Remove(id string) {
    delete(rm, id)
}


// Finds the rooms a client exists in
// Takes a clientID and returns a list container
// with all the rooms the client exists in.
//
func (rm RoomsMap) FindClient(id string) *list.List {
    l := list.New()
    for _, v := range rm {
        if v.HasClient(id) {
            l.PushBack(v)
        }
    }
    return l
}