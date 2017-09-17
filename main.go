package main

import (
	"log"
	"net/http"
	"server/controller"
)

func main() {
	http.HandleFunc("/credit", controller.Querycredit)
	http.HandleFunc("/login/", controller.Welcome)
	http.HandleFunc("/static/", controller.Welcome)
	http.HandleFunc("/grade", controller.Querygrade)
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
/*
package main

import "fmt"

func main() {
	twoD := make([][]int, 3)
	row:=make([]int,3)
	row[0] = 1
    row[1] = 2
    row[2] = 3
    for i := 0; i < 3; i++ {
		if i==2 {
			row[0] = 4
			row[1] = 5
			row[2] = 6
		}
        twoD[i] = make([]int,3)
        twoD[i]=row
    }
    fmt.Println("2d: ", twoD)
}*/
