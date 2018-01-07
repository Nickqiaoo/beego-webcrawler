package controllers

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/astaxie/beego"
	"github.com/axgle/mahonia"
)

var Checkcodeurl = "http://xk1.ahu.cn/CheckCode.aspx?"

type MainController struct {
	beego.Controller
}

func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

func (c *MainController) Login() {
	c.TplName = "login.html"
}

func (c *MainController) Checkcode() {
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		//CheckRedirect: func(req *http.Request, via []*http.Request) error {
		//  return http.ErrUseLastResponse
		// },
		Jar: jar,
	}
	req, _ := client.Get(Checkcodeurl)
	cook := client.Jar.Cookies(checkcodeurl)
	//fmt.Println(cook[0].Name, cook[0].Value)
	c.Ctx.Output.Cookie(cook[0].Name, cook[0].Value)
	imagecode, _ := ioutil.ReadAll(req.Body)
	c.Ctx.Output.Body(imagecode)
}

func (c *MainController) Craw() {
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
	}
	var err error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err = c.Ctx.Request.Cookie("ASP.NET_SessionId")
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	client.Jar.SetCookies(checkcodeurl, cookie)
	c.Ctx.Request.ParseForm()
	//fmt.Println("username:", c.Ctx.Request.Form["username"])
	//fmt.Println("password:", c.Ctx.Request.Form["password"])
	//fmt.Println("yzm", c.Ctx.Request.Form["yzm"])
	if len(c.Ctx.Request.Form["username"])==0{
		c.TplName = "fault.html"
		return
	}


	info := spider(c.Ctx.Request.Form["username"][0], c.Ctx.Request.Form["password"][0], c.Ctx.Request.Form["yzm"][0], &client)
	if info == nil {
		c.TplName = "fault.html"
		return
	}
	c.Data["Name"] = info.Name
	c.Data["Num"] = info.Num
	c.TplName = "welcome.html"
}

func (c *MainController) Querycredit() {
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
	}
	var err error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err = c.Ctx.Request.Cookie("ASP.NET_SessionId")
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	client.Jar.SetCookies(checkcodeurl, cookie)
	c.Ctx.Request.ParseForm()
	//fmt.Println("name:", c.Ctx.Request.Form["name"])
	//fmt.Println("num:", c.Ctx.Request.Form["num"])

	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	cname := encoder.ConvertString(c.Ctx.Request.Form["name"][0])
	resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + c.Ctx.Request.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
	//fmt.Println(resulturl)
	req, _ := http.NewRequest("GET", resulturl, nil)

	req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+c.Ctx.Request.Form["num"][0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	//checkErr(err)
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	fmt.Println("查询页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}

	doc := decoder.NewReader(response.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	view, _ := result.Find("#__VIEWSTATE").Attr("value")
	event, _ := result.Find("#__EVENTVALIDATION").Attr("value")
	v := url.Values{}
	v.Add("Button1", encoder.ConvertString("成绩统计"))
	v.Add("__EVENTTARGET", "")
	v.Add("__EVENTARGUMENT", "")
	v.Add("__VIEWSTATE", view)
	v.Add("hidLanguage", "")
	v.Add("ddlXN", "")
	v.Add("ddlXQ", "")
	v.Add("ddl_kcxz", "")
	v.Add("__EVENTVALIDATION", event)

	body := strings.NewReader(v.Encode())
	req, err = http.NewRequest("POST", resulturl, body)
	checkErr(err)
	req.Header.Add("Referer", resulturl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	/*
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}*/
	response, err = client.Do(req)
	checkErr(err)
	fmt.Println("结果页", response.Status)
	ma, xf, jd := matchcredit(response)
	//fmt.Println(jd)
	c.Data["Name"] = decoder.ConvertString(cname)
	c.Data["Num"] = c.Ctx.Request.Form["num"][0]
	c.Data["Xf"] = xf
	c.Data["Jd"] = jd
	c.Data["Res"] = ma
	c.TplName = "credit.html"
}

func (c *MainController) Querygrade() {
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
	}
	var err error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err = c.Ctx.Request.Cookie("ASP.NET_SessionId")
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	client.Jar.SetCookies(checkcodeurl, cookie)
	c.Ctx.Request.ParseForm()

	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	cname := encoder.ConvertString(c.Ctx.Request.Form["name"][0])
	resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + c.Ctx.Request.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"

	req, _ := http.NewRequest("GET", resulturl, nil)

	req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+c.Ctx.Request.Form["num"][0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	//checkErr(err)
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	fmt.Println("查询页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}

	doc := decoder.NewReader(response.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	view, _ := result.Find("#__VIEWSTATE").Attr("value")
	event, _ := result.Find("#__EVENTVALIDATION").Attr("value")
	v := url.Values{}
	v.Add("btn_xq", encoder.ConvertString("学期成绩"))
	v.Add("__EVENTTARGET", "")
	v.Add("__EVENTARGUMENT", "")
	v.Add("__VIEWSTATE", view)
	v.Add("hidLanguage", "")
	v.Add("ddlXN", "2017-2018")
	v.Add("ddlXQ", "1")
	v.Add("ddl_kcxz", "")
	v.Add("__EVENTVALIDATION", event)

	body := strings.NewReader(v.Encode())
	req, err = http.NewRequest("POST", resulturl, body)
	checkErr(err)
	req.Header.Add("Referer", resulturl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	/*
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}*/
	response, err = client.Do(req)
	checkErr(err)
	fmt.Println("结果页", response.Status)
	grade := matchgrade(response)
	c.Data["Name"] = decoder.ConvertString(cname)
	c.Data["Num"] = c.Ctx.Request.Form["num"][0]
	c.Data["Graderesult"] =grade
	c.TplName = "grade.html"
}
