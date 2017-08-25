package main

import(
	"server/controller"
	"net/http"
	"log"
)

func main() {
	http.HandleFunc("/",controller.Home)
	err := http.ListenAndServe(":9090", nil)
	if err != (nil) {
		log.Fatal("ListenAndServe:", err)
	}
}