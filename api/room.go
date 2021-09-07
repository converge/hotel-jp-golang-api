package api

import (
	"encoding/json"
	"log"
	"net/http"

	"github.com/gorilla/mux"

	"github.com/google/uuid"
)

type Room struct {
	ID           uuid.UUID `json:"id"`
	Availability string    `json:"availability"`
}

type SuccessResponse struct {
	SuccessResponse string `json:"success"`
}

type Problem struct {
	Description string `json:"description"`
}

// this is temporary, and will be replaced when there is a DB
var room Room

type Availability int

const (
	free Availability = iota
	reserved
	inuse
)

var availability = [...]string{"free", "reserved", "inuse"}

func (av Availability) String() string {
	return availability[av]
}

func (r Room) IsValid() bool {
	for _, item := range availability {
		if item == r.Availability {
			return true
		}
	}
	return false
}

func CreateRoom(w http.ResponseWriter, r *http.Request) {
	var receivedRoom Room
	err := json.NewDecoder(r.Body).Decode(&receivedRoom)
	if err != nil {
		log.Println(err)
	}

	if !receivedRoom.IsValid() {
		var problem = Problem{
			Description: "Invalid availability value",
		}
		w.WriteHeader(http.StatusBadRequest)
		w.Header().Set("Content-Type", "application/json")
		err := json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}

	// create new room
	room = Room{
		ID:           uuid.New(),
		Availability: receivedRoom.Availability,
	}
	success := SuccessResponse{
		SuccessResponse: "Room created!",
	}
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(success)
}

func GetRoom(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	if room.ID.ID() == 0 {
		var problem = Problem{Description: "Room not found!"}
		w.WriteHeader(http.StatusNotFound)
		err := json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}
	err := json.NewEncoder(w).Encode(room)
	if err != nil {
		log.Println(err)
	}
}

func DeleteRoom(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	id := params["id"]

	if id != room.ID.String() {
		var problem = Problem{
			Description: "Unable to delete, room not found!",
		}
		w.Header().Set("Content-Type", "application/json")
		w.WriteHeader(404)
		err := json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}
	// reset Room
	room = Room{}
	w.WriteHeader(http.StatusNoContent)
}

func PatchRoom(w http.ResponseWriter, r *http.Request) {

	var patchRoom Room
	err := json.NewDecoder(r.Body).Decode(&patchRoom)
	if err != nil {
		log.Println(err)
	}

	params := mux.Vars(r)
	id := params["id"]
	uId, err := uuid.Parse(id)
	if err != nil {
		log.Println(err)
	}
	patchRoom.ID = uId

	w.Header().Set("Content-Type", "application/json")
	if patchRoom.ID == room.ID && patchRoom.Availability == room.Availability {
		w.WriteHeader(http.StatusNoContent)
		return
	} else if patchRoom.ID != room.ID {
		problem := Problem{
			Description: "Room not found!",
		}
		w.WriteHeader(404)
		err = json.NewEncoder(w).Encode(problem)
		if err != nil {
			log.Println(err)
		}
		return
	}
	room.Availability = patchRoom.Availability

	success := SuccessResponse{
		SuccessResponse: "Entity updated!",
	}
	w.WriteHeader(200)
	err = json.NewEncoder(w).Encode(success)
	if err != nil {
		log.Println(err)
	}
}
