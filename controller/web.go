package controller

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"os"
	"strings"
)

var url2 = "http://xk1.ahu.cn/CheckCode.aspx?"
var jar, _ = cookiejar.New(nil)
var c = &http.Client{
	//CheckRedirect: func(req *http.Request, via []*http.Request) error {
	//  return http.ErrUseLastResponse
	// },
	Jar: jar,
}

//  Home
func Home(w http.ResponseWriter, r *http.Request) {
	fmt.Println("method", r.Method)
	if r.Method == "GET" {
		if strings.HasPrefix(r.URL.Path, "/static") {
			fmt.Println(r.URL.Path)
			if strings.HasPrefix(r.URL.Path, "/static/im") {
				//GE验证码
				file1, err := os.OpenFile("./static/image.jpg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
				defer file1.Close()
				req, err := c.Get(url2)
				image, _ := ioutil.ReadAll(req.Body)
				file1.Write(image)
				//var imagecode string
				fmt.Println("请输入验证码")
				//fmt.Scanf("%s", &imagecode)
			}
			file := "static" + r.URL.Path[len("/static"):]
			http.ServeFile(w, r, file)
			return
		}

		t, err := template.ParseFiles("view/login.html", "view/footer.html", "view/header.html")
		if err != (nil) {
			log.Fatal("template:", err)
		}
		t.ExecuteTemplate(w, "login", nil)
	}
	if r.Method == "POST" {
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Println("yzm", r.Form["yzm"])
		spider(r.Form["username"][0], r.Form["password"][0], r.Form["yzm"][0], c)
		fmt.Fprintf(w, "Success")
	}
}
