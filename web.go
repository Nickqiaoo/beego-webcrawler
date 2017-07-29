package main

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method", r.Method)
	if r.Method == "GET" {
		if strings.HasPrefix(r.URL.Path, "/static") {
			file := "static" + r.URL.Path[len("/static"):]
			http.ServeFile(w, r, file)
			return
		}
		t, err := template.ParseFiles("view/login.html","view/footer.html","view/header.html")
		if err != (nil) {
			log.Fatal("template:", err)
		}
		t.ExecuteTemplate (w,"login",nil)
	}
}
func login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Fprintf(w, "Success")
	}
}
func main() {
	http.HandleFunc("/", home)
	http.HandleFunc("/login", login)
	err := http.ListenAndServe(":9090", nil)
	if err != (nil) {
		log.Fatal("ListenAndServe:", err)
	}
}
