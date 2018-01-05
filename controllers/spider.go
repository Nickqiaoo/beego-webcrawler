package controllers

import (
	"fmt"
	"net/http"
	"net/url"
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
			graderesult=append(graderesult, row)
		}
	})
	//for k,v:=range graderesult{
		//fmt.Println(k,v)
	//}
	return   graderesult
}

func matchcredit(response1 *http.Response) (map[string]string, string, string) {
	ma := make(map[string]string)
	dec := mahonia.NewDecoder("gbk")
	doc := dec.NewReader(response1.Body)
	result, _ := goquery.NewDocumentFromReader(doc)
	result.Find(".datelist").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			ma[s.Find("td").Eq(0).Text()] = s.Find("td").Eq(2).Text()
		}
	})
	jd := result.Find("#pjxfjd").Text()
	xf := result.Find("#xftj").Text()
	return ma, xf, jd
}
