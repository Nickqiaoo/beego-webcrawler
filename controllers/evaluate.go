package controllers

import (
	//"io/ioutil"
	"log"
	"net/http"
	"net/http/cookiejar"
	"net/url"
	"strconv"
	"strings"
	//"time"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

// Evaluate 教学评价
func (c *MainController) Evaluate() {
	//初始化client
	var course []string
	jar, _ := cookiejar.New(nil)
	checkcodeurl, _ := url.Parse(Checkcodeurl)
	client := http.Client{
		Jar: jar,
	}
	var err error
	cookie := make([]*http.Cookie, 1)
	cookie[0], err = c.Ctx.Request.Cookie("ASP.NET_SessionId")
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	client.Jar.SetCookies(checkcodeurl, cookie)
	c.Ctx.Request.ParseForm()
	//fmt.Println("name:", c.Ctx.Request.Form["name"])
	//fmt.Println("num:", c.Ctx.Request.Form["num"])

	//获取主页
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	resulturl := "http://jw3.ahu.cn/xs_main.aspx?xh=" + c.Ctx.Request.Form["num"][0]
	req, _ := http.NewRequest("GET", resulturl, nil)
	req.Close = true
	req.Header.Add("Referer", resulturl)
	req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response, err := client.Do(req)
	if err != nil {
		log.Println(err)
		c.TplName = "fault.html"
		return
	}
	log.Println(c.Ctx.Request.Form["num"][0], c.Ctx.Request.Form["name"][0], "评价-获取首页", response.Status)
	if response.StatusCode != 200 {
		c.TplName = "fault.html"
		return
	}

	//获取所有课程
	doc := decoder.NewReader(response.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	result.Find("div#headDiv").Find("ul.nav").Find("li.top").Eq(3).Find("ul.sub").Find("li").Each(func(i int, s *goquery.Selection) {
		ref, a := s.Find("a").Attr("href")
		if a == false {
			log.Println("未找到课程列表")
		}
		course = append(course, "http://jw3.ahu.cn/"+ref)
	})
	log.Println("课程数量：", len(course))

	var res string

	if len(course) == 0 {
		res = "您已评价过"
	}

	//遍历所有课程
	for i, Url := range course {
		v := url.Values{}
		req, _ := http.NewRequest("GET", Url, nil)
		req.Close = true
		req.Header.Add("Referer", "http://jw3.ahu.cn/xs_main.aspx?xh="+c.Ctx.Request.Form["num"][0])
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
		response, err := client.Do(req)
		if err != nil {
			log.Println(err)
			c.TplName = "fault.html"
			return
		}
		log.Println(c.Ctx.Request.Form["num"][0], c.Ctx.Request.Form["name"][0], "评价-获取课程页", response.Status)
		if response.StatusCode != 200 {
			c.TplName = "fault.html"
			return
		}

		//获取页面view
		doc := decoder.NewReader(response.Body)
		result, _ := goquery.NewDocumentFromReader(doc)
		view, _ := result.Find("#__VIEWSTATE").Attr("value")
		event, _ := result.Find("#__EVENTVALIDATION").Attr("value")

		//获取教师数目
		num := result.Find("#DataGrid1").Find("tbody").Find("tr.alt").Eq(0).Find("td").Length() - 2
		log.Println("教师数目：", num)
		//fmt.Println(Url[35:64])
		//构造post数据
		for k := 1; k <= num; k++ {
			for j := 2; j <= 8; j++ {
				var s1 string
				if k > 2 {
					s1 = "DataGrid1$ctl0" + strconv.Itoa(j) + "$js" + strconv.Itoa(k)
				} else {
					s1 = "DataGrid1$ctl0" + strconv.Itoa(j) + "$JS" + strconv.Itoa(k)
				}
				s2 := "DataGrid1$ctl0" + strconv.Itoa(j) + "$txtjs" + strconv.Itoa(k)
				if j == 2 {
					v.Add(s1, encoder.ConvertString("良好"))
				} else {
					v.Add(s1, encoder.ConvertString("优秀"))
				}
				v.Add(s2, "")
			}
		}
		v.Add("pkjc", Url[35:64])
		v.Add("__EVENTTARGET", "")
		v.Add("__EVENTARGUMENT", "")
		v.Add("__VIEWSTATE", view)
		v.Add("__LASTFOCUS", "")
		v.Add("pjxx", "")
		v.Add("txt1", "")
		v.Add("TextBox1", "0")
		v.Add("__EVENTVALIDATION", event)
		v.Add("Button1", encoder.ConvertString("保 存"))

		body := strings.NewReader(v.Encode())
		req, err = http.NewRequest("POST", Url, body)
		req.Close = true
		if err != nil {
			log.Println(err)
			c.TplName = "fault.html"
			return
		}
		req.Header.Add("Referer", Url)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
		response, err = client.Do(req)
		if err != nil {
			log.Println(err)
			c.TplName = "fault.html"
			return
		}
		log.Println(i, "保存成功", response.Status)
		//最后一门课时提交
		if i == len(course)-1 {
			v.Del("Button1")
			v.Add("Button2", encoder.ConvertString(" 提  交 "))
			body := strings.NewReader(v.Encode())
			req, err = http.NewRequest("POST", Url, body)
			req.Close = true
			if err != nil {
				log.Println(err)
				c.TplName = "fault.html"
				return
			}
			req.Header.Add("Referer", Url)
			req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
			req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
			response, err = client.Do(req)
			if err != nil {
				log.Println(err)
				c.TplName = "fault.html"
				return
			}
			log.Println(c.Ctx.Request.Form["num"][0], c.Ctx.Request.Form["name"][0], "提交成功", response.Status)
		}

		if i == len(course)-1 {
			res = "评价成功"
		}
		//time.Sleep(5 * time.Second)
	}
	c.Data["Res"] = res
	c.TplName = "evaluate.html"
}
