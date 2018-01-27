package controllers

import (
	"time"
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

// Checkcodeurl 验证码url
var Checkcodeurl = "http://xk1.ahu.cn/CheckCode.aspx?"

type MainController struct {
	beego.Controller
}


// Login 返回登陆界面
func (c *MainController) Login() {
	c.TplName = "login.html"
}

// Checkcode 获取验证码
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
	c.Ctx.Output.Cookie(cook[0].Name, cook[0].Value)
	imagecode, _ := ioutil.ReadAll(req.Body)
	c.Ctx.Output.Body(imagecode)
}

// Craw 登陆函数
func (c *MainController) Craw() {
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
		Timeout:time.Second*10,
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
	//log.Println("username:", c.Ctx.Request.Form["username"])
	//log.Println("password:", c.Ctx.Request.Form["password"])
	//log.Println("yzm", c.Ctx.Request.Form["yzm"])
	if len(c.Ctx.Request.Form["username"])==0{
		c.TplName = "fault.html"
		return
	}

	//获取主页
	url1 := "http://xk1.ahu.cn/default2.aspx"
	v := url.Values{}
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	but := encoder.ConvertString("学生")
	v.Add("__VIEWSTATE", "/wEPDwUJODk4OTczODQxZGQhFC7x2TzAGZQfpidAZYYjo/LeoQ==")
	v.Add("txtUserName", c.Ctx.Request.Form["username"][0])
	v.Add("TextBox2", c.Ctx.Request.Form["password"][0])
	v.Add("txtSecretCode", c.Ctx.Request.Form["yzm"][0])
	v.Add("RadioButtonList1", but)
	v.Add("Button1", "")
	v.Add("lbLanguage", "")
	v.Add("hidPdrs", "")
	v.Add("hidsc", "")
	v.Add("__EVENTVALIDATION", "/wEWDgKX/4yyDQKl1bKzCQLs0fbZDAKEs66uBwK/wuqQDgKAqenNDQLN7c0VAuaMg+INAveMotMNAoznisYGArursYYIAt+RzN8IApObsvIHArWNqOoPqeRyuQR+OEZezxvi70FKdYMjxzk=")
	
	//建立client发送POST请求
	body := strings.NewReader(v.Encode())
	r, _ := http.NewRequest("POST", url1, body)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Referer", "http://xk1.ahu.cn/default2.aspx")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(r)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println( c.Ctx.Request.Form["username"][0],"登陆-主页获取成功", response.Status)

	//解析主页，如果有欢迎则说明获取失败
	doc := decoder.NewReader(response.Body)
	result, err := goquery.NewDocumentFromReader(doc)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	cname := result.Find("title").Text()
	if strings.HasPrefix(cname, "欢迎") {
		log.Println( c.Ctx.Request.Form["username"][0],"主页获取错误")
		c.TplName = "fault.html"
		return
	}
	cname = result.Find("#xhxm").Text()
	cname = strings.TrimRight(cname, "同学")
	//return &Info{Name: cname, Num: username}
	client.Get("https://sc.ftqq.com/SCU20914Teefb444fcce3027f14828723ca1cd65e5a6c2b88500ab.send?text="+
	url.QueryEscape(c.Ctx.Request.Form["username"][0]+" "+cname+" 登陆"))
	c.Data["Name"] = cname
	c.Data["Num"] = c.Ctx.Request.Form["username"][0]
	c.TplName = "welcome.html"
}

// Querycredit 查询学分
func (c *MainController) Querycredit() {
	//初始化client
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
		Timeout:time.Second*10,
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

	//获取查询页
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	cname := encoder.ConvertString(c.Ctx.Request.Form["name"][0])
	resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + c.Ctx.Request.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
	req, _ := http.NewRequest("GET", resulturl, nil)
	req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+c.Ctx.Request.Form["num"][0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["num"][0],"学分查询页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}

	//获取view，event
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

	//构建新请求
	body := strings.NewReader(v.Encode())
	req, err = http.NewRequest("POST", resulturl, body)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	req.Header.Add("Referer", resulturl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err = client.Do(req)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["num"][0],c.Ctx.Request.Form["name"][0],"学分结果页", response.Status)
	ma, xf, jd := matchcredit(response)
	c.Data["Name"] = decoder.ConvertString(cname)
	c.Data["Num"] = c.Ctx.Request.Form["num"][0]
	c.Data["Xf"] = xf
	c.Data["Jd"] = jd
	c.Data["Res"] = ma
	c.TplName = "credit.html"
}

// Querygrade 查询成绩
func (c *MainController) Querygrade() {
	//初始化client
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
		Timeout:time.Second*10,
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

	//获取查询页
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	cname := encoder.ConvertString(c.Ctx.Request.Form["name"][0])
	resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + c.Ctx.Request.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
	req, _ := http.NewRequest("GET", resulturl, nil)
	req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+c.Ctx.Request.Form["num"][0])
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["num"][0],"成绩查询页", response.Status)
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

	//构造新请求
	body := strings.NewReader(v.Encode())
	req, err = http.NewRequest("POST", resulturl, body)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	req.Header.Add("Referer", resulturl)
	req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err = client.Do(req)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["num"][0],c.Ctx.Request.Form["name"][0],"成绩结果页", response.Status)
	grade := matchgrade(response)
	c.Data["Name"] = decoder.ConvertString(cname)
	c.Data["Num"] = c.Ctx.Request.Form["num"][0]
	c.Data["Graderesult"] =grade
	c.TplName = "grade.html"
}
