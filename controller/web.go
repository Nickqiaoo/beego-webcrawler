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
	"sync"
	"net/url"
)

var mutex sync.Mutex
var Url2 = "http://xk1.ahu.cn/CheckCode.aspx?"
var Ma map[string] []*http.Cookie

//  Home
func Home(w http.ResponseWriter, r *http.Request) {
	if Ma==nil{
		Ma=make( map[string] []*http.Cookie)
	}
	jar,_:=cookiejar.New(nil)
	u,_:=url.Parse(Url2)
 	c := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//  return http.ErrUseLastResponse
		// },
		Jar: jar,
	}
	fmt.Println("method", r.Method)
	t, err := template.ParseFiles("view/result.html", "view/footer.html", "view/header.html")
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
				req, err := c.Get(Url2)
				fmt.Println(c.Jar.Cookies(u))
				cookie := http.Cookie{Name: "testcookiename", Value: "testcookievalue", Path: "/", MaxAge: 86400}
				Ma[cookie.Name]=c.Jar.Cookies(u)
				http.SetCookie(w, &cookie)
				//c.Jar,_=cookiejar.New(nil)
				image, _ := ioutil.ReadAll(req.Body)
				file1.Write(image)
				//var imagecode string
				fmt.Println("请输入验证码")
				//fmt.Scanf("%s", &imagecode)
			}
			file := "static" + r.URL.Path[len("/static"):]
			http.ServeFile(w, r, file)
		}

		t, err := template.ParseFiles("view/login.html", "view/footer.html", "view/header.html")
		if err != (nil) {
			log.Fatal("template:", err)
		}
		t.ExecuteTemplate(w, "login", nil)
	}
	if r.Method == "POST" {
		for k,v:= range(Ma){
			fmt.Println(k,v)
		}
		c.Jar.SetCookies(u,Ma["testcookiename"])
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Println("yzm", r.Form["yzm"])
		result:=spider(r.Form["username"][0], r.Form["password"][0], r.Form["yzm"][0], &c )
		if err != (nil) {
			log.Fatal("template:", err)
		}
		t.ExecuteTemplate(w, "result", *result)
		//fmt.Fprintf(w, "Success")
	}
}
