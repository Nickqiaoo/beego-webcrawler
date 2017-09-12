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
/*
package main

import(
	//"server/controller"
	"os"
	"log"
	//"github.com/axgle/mahonia"
	"github.com/PuerkitoBio/goquery"
	"fmt"

)

func main() {
	file3, err := os.OpenFile("spider.txt", os.O_RDWR, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file3.Close()
	//decoder := mahonia.NewEncoder("gbk")
	result, _ := goquery.NewDocumentFromReader(file3)
	s ,_:=result.Find("#__EVENTVALIDATION").Attr("value")
	fmt.Println(s)
}*/