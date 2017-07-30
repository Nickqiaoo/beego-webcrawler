package controller

import (
	"fmt"
	"html/template"
	"log"
	"net/http"
	"strings"
)

func Home(w http.ResponseWriter, r *http.Request) {
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
func Login(w http.ResponseWriter, r *http.Request) {
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		spider(r.Form["username"][0],r.Form["password"][0])
		fmt.Fprintf(w, "Success")
	}
}

