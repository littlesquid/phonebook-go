package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/google/uuid"
	"github.com/gorilla/mux"
)

type Person struct {
	ID        string   `json:"id,omitempty"`
	Firstname string   `json:"firstname, omitempty"`
	Lastname  string   `json:"lastname, omitempty"`
	Address   *Address `json:"address, omitempty"`
}

type Address struct {
	City  string `json:"city, omitempty"`
	State string `json:"state, omitempty"`
}

var people []Person

func main() {
	router := mux.NewRouter()
	addData()
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func addData() {
	people = append(people, Person{ID: uuid.New().String(), Firstname: "John", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
	people = append(people, Person{ID: uuid.New().String(), Firstname: "Jane", Lastname: "Doe", Address: &Address{City: "City X", State: "State X"}})
}

func clearData() {
	people = people[:0]
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	fmt.Println(params)
	for _, item := range people {
		fmt.Println(item.ID)
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = uuid.New().String()
	people = append(people, person)
	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	params := mux.Vars(r)
	for index, item := range people {
		if item.ID == params["id"] {
			people = append(people[:index], people[index+1:]...)
			break
		}
		json.NewEncoder(w).Encode(people)
	}
}
