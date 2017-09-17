package controller

import (
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"os"
	"html/template"
	"net/http/cookiejar"
	"strings"

	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

type Result struct {
	Num  string
	Name string
	Xf   string
	Res  map[string]string
	Jd   string
}

type Info struct {
	Name string
	Num  string
}

type Grade struct{
	Num  string
	Name string
	Graderesult [][]string 
}

func spider(username string, password string, imagecode string, c *http.Client) *Info {
	u, _ := url.Parse(Url2)
	fmt.Println(c.Jar.Cookies(u))
	url1 := "http://xk1.ahu.cn/default2.aspx"
	v := url.Values{}
	encoder := mahonia.NewEncoder("gbk")
	decoder := mahonia.NewDecoder("gbk")
	but := encoder.ConvertString("学生")
	v.Add("__VIEWSTATE", "/wEPDwUJODk4OTczODQxZGQhFC7x2TzAGZQfpidAZYYjo/LeoQ==")
	v.Add("txtUserName", username)
	v.Add("TextBox2", password)
	v.Add("txtSecretCode", imagecode)
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
	response, err := c.Do(r)
	fmt.Println("主页", response.Status)
	checkErr(err)

	doc := decoder.NewReader(response.Body)
	result, err := goquery.NewDocumentFromReader(doc)
	checkErr(err)
	cname := result.Find("title").Text()
	if strings.HasPrefix(cname, "欢迎") {
		fmt.Println("主页获取错误")
		return nil
	}
	cname = result.Find("#xhxm").Text()
	cname = strings.TrimRight(cname, "同学")
	return &Info{Name: cname, Num: username}
}

// Querycredit 查询学分绩点
func Querycredit(w http.ResponseWriter, r *http.Request) {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(Url2)
	c := http.Client{
		Jar: jar,
	}
	fmt.Println("method", r.Method)
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
		fmt.Println("name:", r.Form["name"])
		fmt.Println("num:", r.Form["num"])

		encoder := mahonia.NewEncoder("gbk")
		decoder := mahonia.NewDecoder("gbk")
		cname := encoder.ConvertString(r.Form["name"][0])
		resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + r.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
		fmt.Println(resulturl)
		req, _ := http.NewRequest("GET", resulturl, nil)

		req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+r.Form["num"][0])
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
		response, err := c.Do(req)
		//checkErr(err)
		if err != nil {
			fault(&w)
			return
		}
		fmt.Println("查询页", response.Status)
		if response.StatusCode != 200 {
			fault(&w)
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
		/*for k,va :=range r.Header{
			fmt.Println(k,va)
		}
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}*/
		response, err = c.Do(req)
		checkErr(err)
		fmt.Println("结果页", response.Status)
		ma, xf, jd := matchcredit(response)
		fmt.Println(jd)
		t, err := template.ParseFiles("view/credit.html", "view/footer.html", "view/header.html")
		checkErr(err)
		err = t.ExecuteTemplate(w, "credit", Result{Name: decoder.ConvertString(cname), Num: r.Form["num"][0], Xf: xf, Res: ma, Jd: jd})
		checkErr(err)
	} else if r.Method == "GET" {
		fault(&w)
		return
	}
}


// Querygrade 查询成绩
func Querygrade(w http.ResponseWriter, r *http.Request) {
	jar, _ := cookiejar.New(nil)
	u, _ := url.Parse(Url2)
	c := http.Client{
		Jar: jar,
	}
	fmt.Println("method", r.Method)
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
		fmt.Println("name:", r.Form["name"])
		fmt.Println("num:", r.Form["num"])

		encoder := mahonia.NewEncoder("gbk")
		decoder := mahonia.NewDecoder("gbk")
		cname := encoder.ConvertString(r.Form["name"][0])
		resulturl := "http://xk1.ahu.cn/xscjcx.aspx?xh=" + r.Form["num"][0] + "&xm=" + url.QueryEscape(cname) + "&gnmkdm=N121605"
		fmt.Println(resulturl)
		req, _ := http.NewRequest("GET", resulturl, nil)

		req.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+r.Form["num"][0])
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
		response, err := c.Do(req)
		//checkErr(err)
		if err != nil {
			fault(&w)
			return
		}
		fmt.Println("查询页", response.Status)
		if response.StatusCode != 200 {
			fault(&w)
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
		v.Add("ddlXN", "2016-2017")
		v.Add("ddlXQ", "2")
		v.Add("ddl_kcxz", "")
		v.Add("__EVENTVALIDATION", event)

		body := strings.NewReader(v.Encode())
		req, err = http.NewRequest("POST", resulturl, body)
		checkErr(err)
		req.Header.Add("Referer", resulturl)
		req.Header.Add("Content-Type", "application/x-www-form-urlencoded")
		req.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
		/*for k,va :=range r.Header{
			fmt.Println(k,va)
		}
		c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
			return http.ErrUseLastResponse
		}*/
		response, err = c.Do(req)
		checkErr(err)
		fmt.Println("结果页", response.Status)
		grade := matchgrade(response)
		t, err := template.ParseFiles("view/grade.html", "view/footer.html", "view/header.html")
		checkErr(err)
		err = t.ExecuteTemplate(w, "grade",Grade{Name: decoder.ConvertString(cname), Num: r.Form["num"][0],Graderesult:grade})
		checkErr(err)
	} else if r.Method == "GET" {
		fault(&w)
		return
	}
}

func matchgrade(response1 *http.Response) [][]string {
	dec := mahonia.NewDecoder("gbk")
	doc := dec.NewReader(response1.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	graderesult:=make([][]string,0)
	result.Find(".datelist").Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			row:=make([]string,6)
			row[0]=s.Find("td").Eq(3).Text()
			row[1]=s.Find("td").Eq(4).Text()
			row[2]=s.Find("td").Eq(6).Text()
			row[3]=s.Find("td").Eq(7).Text()
			row[4]=s.Find("td").Eq(8).Text()
			row[5]=s.Find("td").Eq(12).Text()
			//fmt.Println(row)
			graderesult=append(graderesult, row)
			/*for k,v:=range graderesult{
				fmt.Println(i,k,v)
			}*/
		}
	})
	for k,v:=range graderesult{
		fmt.Println(k,v)
	}
	return   graderesult
}

func matchcredit(response1 *http.Response) (map[string]string, string, string) {
	ma := make(map[string]string)
	dec := mahonia.NewDecoder("gbk")
	doc := dec.NewReader(response1.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	/*file, err := os.OpenFile("spider1.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	checkErr(err)
	defer file.Close()
	html, _ := ioutil.ReadAll(response1.Body)
	num, err := file.Write(html)
	if num != len(html) {
		log.Fatal(err)
	}*/
	result.Find(".datelist").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			ma[s.Find("td").Eq(0).Text()] = s.Find("td").Eq(2).Text()
			fmt.Println(s.Find("td").Eq(0).Text(), s.Find("td").Eq(2).Text())
		}
	})
	jd := result.Find("#pjxfjd").Text()
	xf := result.Find("#xftj").Text()
	return ma, xf, jd
}
func checkErr(err error) {
	if err != nil {
		log.Fatal(err)
	}
}
