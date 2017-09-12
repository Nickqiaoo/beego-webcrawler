package controller

import (
	"fmt"
	//"io/ioutil"
	"log"
	"net/http"
	"net/url"
	//"os"
	"strings"
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

type Result struct {
	Name string
	Xf string
	Res  map[string]string
	Jd   string
}

func spider(username string, password string, imagecode string, c *http.Client) *Result {
	u,_:=url.Parse(Url2)
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
	checkErr(err)

	doc := decoder.NewReader(response.Body)
	result, err := goquery.NewDocumentFromReader(doc)
	checkErr(err)
	cname :=result.Find("#xhxm").Text()
	cname= strings.TrimRight(cname,"同学")
	cname= encoder.ConvertString(cname)
	resulturl:="http://xk1.ahu.cn/xscjcx.aspx?xh="+username+"&xm="+url.QueryEscape(cname)+"&gnmkdm=N121605"
	fmt.Println(resulturl)
	r,_=http.NewRequest("GET",resulturl,nil)
	//r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("Referer", "http://xk1.ahu.cn/xs_main.aspx?xh="+username)
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	response,err=c.Do(r)
	checkErr(err)

	doc = decoder.NewReader(response.Body)
	result, _ = goquery.NewDocumentFromReader(doc)
	view,_:=result.Find("#__VIEWSTATE").Attr("value")
	event,_:=result.Find("#__EVENTVALIDATION").Attr("value")
	v = url.Values{}
	v.Add("Button1", encoder.ConvertString("成绩统计"))
	v.Add("__EVENTTARGET", "")
	v.Add("__EVENTARGUMENT", "")
	v.Add("__VIEWSTATE",view )
	v.Add("hidLanguage", "")
	v.Add("ddlXN", "")
	v.Add("ddlXQ", "")
	v.Add("ddl_kcxz", "")
	v.Add("__EVENTVALIDATION", event)

	body = strings.NewReader(v.Encode())
	r, err = http.NewRequest("POST", resulturl, body)
	if err != nil {
		log.Fatal(err)
	}
	r.Header.Add("Referer", resulturl)
	r.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	r.Header.Add("Host", "xk1.ahu.cn")
	r.Header.Add("Origin", "http://xk1.ahu.cn")
	/*for k,va :=range r1.Header{
		fmt.Println(k,va)
	}
	c.CheckRedirect = func(req *http.Request, via []*http.Request) error {
		return http.ErrUseLastResponse
	}*/
	response, err = c.Do(r)
	fmt.Println(c.Jar.Cookies(u))
	if err != nil {
		log.Fatal(err)
	}
	fmt.Println(response.Status)
	ma,xf,jd := match(response)
	fmt.Println(jd)
	return &Result{Name: username, Xf:xf,Res: ma, Jd: jd}
}

func match(response1 *http.Response) (map[string]string, string,string) {
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
	}
	checkErr(err)*/
	result.Find(".datelist").Eq(0).Find("tr").Each(func(i int, s *goquery.Selection) {
		if i > 0 {
			ma[s.Find("td").Eq(0).Text()] = s.Find("td").Eq(2).Text()
			fmt.Println(s.Find("td").Eq(0).Text(),s.Find("td").Eq(2).Text())
		}
	})
	jd := result.Find("#pjxfjd").Text()
	xf:=result.Find("#xftj").Text()
	fmt.Println(jd)
	for k, v := range ma {
		fmt.Println(k, v)
	}
	return ma,xf,jd
}
func checkErr(err error){
	if err!=nil{
		log.Fatal(err)
	}
}
