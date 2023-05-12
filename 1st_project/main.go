package main

import (
	"fmt"
	"log"
	"net/http"
)

func formhandler(w http.ResponseWriter, r *http.Request) {
	err := r.ParseForm()
	if err != nil {
		http.Error(w, "parse form ", 404)
		return
	}
	fmt.Fprintf(w, "post request successfull")
	name := r.FormValue("name")
	address := r.FormValue("address")
	fmt.Fprintf(w, "name %s\n", name)
	fmt.Fprintf(w, "adress is %s\n", address)

}

func hellohandler(w http.ResponseWriter, r *http.Request) {
	if r.URL.Path != "/hello" {
		http.Error(w, "not found", 404)
		return
	}
	if r.Method != "GET" {
		http.Error(w, "not found", 404)
		return
	}
	// this writes to the response
	fmt.Fprintf(w, "hello")
}
func main() {
	fileserver := http.FileServer(http.Dir("./static"))
	http.Handle("/", fileserver)
	http.HandleFunc("/form", formhandler)
	http.HandleFunc("/hello", hellohandler)
	fmt.Println("starting the server at 8080/n")
	err := http.ListenAndServe(":8080", nil)
	if err != nil {
		log.Panic("error in starting the server")
	}
}
