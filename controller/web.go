package controller

import (
	"fmt"
	"html/template"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"os"
	"strings"
	"time"
)

var Url2 = "http://xk1.ahu.cn/CheckCode.aspx?"

// Welcome 登录界面
func Welcome(w http.ResponseWriter, r *http.Request) {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(Url2)
	c := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//  return http.ErrUseLastResponse
		// },
		Jar: jar,
	}
	fmt.Println("method", r.Method)
	if r.Method == "GET" {
		fmt.Println(r.URL.Path)
		if strings.HasPrefix(r.URL.Path, "/static") {
			if strings.HasPrefix(r.URL.Path, "/static/im") {
				//GE验证码
				file1, err := os.OpenFile("./static/image.jpg", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
				if err != nil {
					log.Fatal(err)
				}
				defer file1.Close()

				req, err := c.Get(Url2)
				cook := c.Jar.Cookies(u)
				fmt.Println(cook[0].Name, cook[0].Value)
				cookie := http.Cookie{Name: cook[0].Name, Value: cook[0].Value, Path: "/", MaxAge: 800}
				http.SetCookie(w, &cookie)
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
		var err error
		cookie := make([]*http.Cookie, 1)
		cookie[0], err = r.Cookie("ASP.NET_SessionId")
		if err != nil {
			fault(&w)
			return
		}
		fmt.Println(cookie[0].Value)
		c.Jar.SetCookies(u, cookie)
		r.ParseForm()
		fmt.Println("username:", r.Form["username"])
		fmt.Println("password:", r.Form["password"])
		fmt.Println("yzm", r.Form["yzm"])
		info := spider(r.Form["username"][0], r.Form["password"][0], r.Form["yzm"][0], &c)
		//delete(Ma, ccookie.Value)
		if info == nil {
			fault(&w)
			return
		}
		t, err := template.ParseFiles("view/welcome.html", "view/footer.html", "view/header.html")
		checkErr(err)
		err = t.ExecuteTemplate(w, "welcome", *info)
		checkErr(err)
	}
}

func fault(w *http.ResponseWriter) {
	t, err := template.ParseFiles("view/fault.html", "view/footer.html", "view/header.html")
	if err != (nil) {
		log.Fatal("template:", err)
	}
	t.ExecuteTemplate(*w, "fault", nil)
}

func GetString(le int) string {
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < le; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)
}
