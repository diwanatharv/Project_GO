package main

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/gorilla/mux"
)

type Movie struct {
	Id       string    `json:"Id"`
	Isbn     string    `json:"Isbn"`
	Title    string    `json:"Title"`
	Director *Director `json:"Director"`
}
type Director struct {
	FirstName string `json:"FirstName"`
	LastName  string `json:"LastName"`
}

var movies []Movie

func getMovies(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(movies)
}
func deleteMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	var ind int
	for index, item := range movies {
		if item.Id == params["id"] {
			ind = index
		}
	}
	movies = append(movies[:ind], movies[ind+1:]...)
	//this line is equivalent to c.json
	json.NewEncoder(w).Encode(movies)
}
func getMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)
	for _, item := range movies {
		if item.Id == params["id"] {
			json.NewEncoder(w).Encode(item)
			return
		}
	}
}
func createMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var movie Movie
	//taking input from the reqbody
	_ = json.NewDecoder(r.Body).Decode(&movie)
	movies = append(movies, movie)

	json.NewEncoder(w).Encode(movies)
}
func updateMovie(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	params := mux.Vars(r)

	for index, item := range movies {
		if item.Id == params["id"] {
			movies = append(movies[:index], movies[index+1:]...)
			var movie Movie
			_ = json.NewDecoder(r.Body).Decode(&movie)
			movie.Id = params["id"]
			movies = append(movies, movie)
			json.NewEncoder(w).Encode(movies)
			return
		}
	}
	//this line is equivalent to c.json

}
func main() {
	r := mux.NewRouter()
	movies = append(movies, Movie{Id: "1", Isbn: "11", Title: "kashmiri files", Director: &Director{FirstName: "kamesh", LastName: "verma"}})
	movies = append(movies, Movie{Id: "2", Isbn: "12", Title: "kerla story", Director: &Director{FirstName: "ramesh", LastName: "pandey"}})
	r.HandleFunc("/movies", getMovies).Methods("GET")
	r.HandleFunc("/movie/{id}", getMovie).Methods("GET")
	r.HandleFunc("/movies", createMovie).Methods("POST")
	r.HandleFunc("/movie/{id}", updateMovie).Methods("PUT")
	r.HandleFunc("/movie/{id}", deleteMovie).Methods("DELETE")
	fmt.Printf("starting the server at 8000\n")
	log.Fatal(http.ListenAndServe(":8000", r))
}
