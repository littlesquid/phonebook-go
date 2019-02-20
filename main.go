package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	_ "github.com/go-sql-driver/mysql"
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
	router.HandleFunc("/people", GetPeople).Methods("GET")
	router.HandleFunc("/people/{id}", GetPerson).Methods("GET")
	router.HandleFunc("/people", CreatePerson).Methods("POST")
	router.HandleFunc("/people/{id}", DeletePerson).Methods("DELETE")

	fmt.Println("Connected to port 8000")
	log.Fatal(http.ListenAndServe(":8000", router))
}

func clearData() {
	people = people[:0]
}

func GetPeople(w http.ResponseWriter, r *http.Request) {
	conn := ConnectDatabase()
	defer conn.Close()

	results, err := conn.Query("SELECT id, firstname, lastname FROM phone_book")
	if err != nil {
		panic(err.Error())
	}

	for results.Next() {
		var person Person
		err = results.Scan(&person.ID, &person.Firstname, &person.Lastname)
		if err != nil {
			panic(err.Error())
		}
		people = append(people, person)
	}
	json.NewEncoder(w).Encode(people)
}

func GetPerson(w http.ResponseWriter, r *http.Request) {
	conn := ConnectDatabase()
	defer conn.Close()

	var person Person
	params := mux.Vars(r)
	err := conn.QueryRow("SELECT id, firstname, lastname FROM phone_book where id = ?", params["id"]).Scan(&person.ID, &person.Firstname, &person.Lastname)
	if err != nil {
		panic(err.Error())
	}

	json.NewEncoder(w).Encode(person)

	/*for _, item := range people {
		fmt.Println(item.ID)
		if item.ID == params["id"] {
			json.NewEncoder(w).Encode(item)
		}
	}*/
}

func CreatePerson(w http.ResponseWriter, r *http.Request) {
	defer fmt.Println("complete create person") //run after the surrounding function returns.
	conn := ConnectDatabase()
	defer conn.Close()
	var person Person
	_ = json.NewDecoder(r.Body).Decode(&person)
	person.ID = uuid.New().String()
	people = append(people, person)
	fmt.Println(person.ID)
	insert, err := conn.Query("INSERT INTO `phone_book` (id, firstname, lastname) values('" + person.ID + "', '" + person.Firstname + "', '" + person.Lastname + "')")
	if err != nil {
		panic(err.Error())
	}
	defer insert.Close()

	json.NewEncoder(w).Encode(people)
}

func DeletePerson(w http.ResponseWriter, r *http.Request) {
	conn := ConnectDatabase()
	defer conn.Close()

	params := mux.Vars(r)
	stmt, err := conn.Prepare("DELETE FROM phone_book where id = ?")
	if err != nil {
		panic(err.Error())
	}
	defer stmt.Close()

	_, err = stmt.Exec(params["id"])
	if err != nil {
		panic(err)
	}
	clearData()
	GetPeople(w, r)
	// params := mux.Vars(r)
	// for index, item := range people {
	// 	if item.ID == params["id"] {
	// 		people = append(people[:index], people[index+1:]...)
	// 		break
	// 	}
	// 	json.NewEncoder(w).Encode(people)
	// }
}
