package controller

import (
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"net/url"
	"os"
	"strings"
	"github.com/axgle/mahonia"
)

func spider(username string, password string, imagecode string, c *http.Client) {
	url1 := "http://xk1.ahu.cn/default2.aspx"
	//url2 := "http://xk1.ahu.cn/CheckCode.aspx?"
	//url3 := "http://xk1.ahu.cn/xs_main.aspx?xh=P71514011"
	v := url.Values{}
	enc := mahonia.NewEncoder("gbk")
	but := enc.ConvertString("学生")
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
	if err != nil {
		log.Fatal(err)
	}
	resulturl:="http://xk1.ahu.cn/xscjcx.aspx?xh=P71514011&xm=%u4e54%u82f1%u6770&gnmkdm=N121605"
	v1:=url.Values{}
	v1.Add("Button1", enc.ConvertString("成绩统计") )
	v1.Add("__EVENTTARGET","")
	v1.Add("__EVENTARGUMENT","")
	v1.Add("__VIEWSTATE","/wEPDwUKMTM4ODQwMjE5Nw8WEh4DZGczBQNiamceC3N0cl90YWJfYmpnBRN6Zl9jeGNqdGpfUDcxNTE0MDExHgpTb3J0RXhwcmVzBQRrY21jHgJ4aAUJUDcxNTE0MDExHgh6eGNqY3h4cwUBMR4IU29ydERpcmUFA2FzYx4IY2pjeF9sc2JkHgdkeWJ5c2NqBQExHgZzZmRjYmtlFgICAQ9kFiACBA8QDxYGHg1EYXRhVGV4dEZpZWxkBQJYTh4ORGF0YVZhbHVlRmllbGQFAlhOHgtfIURhdGFCb3VuZGdkEBUDAAkyMDE2LTIwMTcJMjAxNS0yMDE2FQMACTIwMTYtMjAxNwkyMDE1LTIwMTYUKwMDZ2dnZGQCCg8QDxYGHwkFBmtjeHptYx8KBQZrY3h6ZG0fC2dkEBUID+WFrOWFseWfuuehgOivvg/kuJPkuJrmoLjlv4Por74P5LiT5Lia6YCJ5L+u6K++Eui3qOS4k+S4mumAieS/ruivvhXntKDotKjmlZnogrLpgInkv67or74M5a6e6Le15pWZ6IKyD+S4k+S4muWfuuehgOivvgAVCAIwMQIwMwIwNAIwNgIwNwIwOAIwOQAUKwMIZ2dnZ2dnZ2dkZAITDw8WAh4HVmlzaWJsZWhkZAIYDw8WAh8MaGRkAiAPDxYCHwxoZGQCIg8PFgIeBFRleHRlZGQCJA8PFgQfDQUS5a2m5Y+377yaUDcxNTE0MDExHwxnZGQCJg8PFgQfDQUS5aeT5ZCN77ya5LmU6Iux5p2wHwxnZGQCKA8PFgQfDQUh5a2m6Zmi77ya55S15a2Q5L+h5oGv5bel56iL5a2m6ZmiHwxnZGQCKg8PFgQfDQUJ5LiT5Lia77yaHwxnZGQCLA8PFgQfDQUP54mp6IGU572R5bel56iLHwxnZGQCLg8PFgIfDQUN5LiT5Lia5pa55ZCROmRkAjAPDxYEHw0FGuihjOaUv+ePre+8mjE157qn54mp6IGU572RHwxnZGQCNA88KwALAQAPFgIfDGhkZAI2D2QWKgIBDw8WAh8MaGRkAgMPPCsACwEADxYCHwxoFgIeBXN0eWxlBQxESVNQTEFZOm5vbmVkAgUPZBYCAg0PPCsACwBkAgcPDxYEHw0FHuiHs+S7iuacqumAmui/h+ivvueoi+aIkOe7qe+8mh8MZ2RkAgkPPCsACwEADxYIHghEYXRhS2V5cxYAHgtfIUl0ZW1Db3VudGYeCVBhZ2VDb3VudAIBHhVfIURhdGFTb3VyY2VJdGVtQ291bnRmFgIfDgUNRElTUExBWTpibG9ja2QCDQ88KwALAQAPFgIfDGgWAh8OBQxESVNQTEFZOm5vbmVkAg8PPCsACwBkAhEPPCsACwEADxYCHwxoFgIfDgUMRElTUExBWTpub25lZAIVDzwrAAsAZAIXDzwrAAsBAA8WAh8MaBYCHw4FDERJU1BMQVk6bm9uZWQCGA88KwALAQAPFgIfDGgWAh8OBQxESVNQTEFZOm5vbmVkAhkPPCsACwEADxYCHwxoZGQCGw88KwALAQAPFgIfDGgWAh8OBQxESVNQTEFZOm5vbmVkAh0PPCsACwEADxYCHwxoFgIfDgUMRElTUExBWTpub25lZAIfDzwrAAsBARQrAAdkZDwrAAQBABYCHgpIZWFkZXJUZXh0BQzliJvmlrDlhoXlrrk8KwAEAQAWAh8TBQzliJvmlrDlrabliIY8KwAEAQAWAh8TBQzliJvmlrDmrKHmlbBkZGQCIQ8PFgQfDQUR5pys5LiT5Lia5YWxNTbkurofDGhkZAIjDw8WAh8MaGRkAisPDxYCHwxoZGQCMQ8PFgIfDGhkZAIzDw8WAh8NBQNBSFVkZAI0Dw8WAh4ISW1hZ2VVcmwFFS4vZXhjZWwvUDcxNTE0MDExLmpwZ2RkAjgPZBYCAgMPPCsACwBkZM99lZ66NeQGCqqB5JnUoP0uvChT")
	v1.Add("hidLanguage","")
	v1.Add("ddlXN","")
	v1.Add("ddlXQ","")
	v1.Add("ddl_kcxz","")
	v1.Add("__EVENTVALIDATION","/wEWGQKRhtjxBwLd4qCqDwLuwOmEBQK+q8LECwK/q/6ECgLfwOmEBQLQr8PqCQLRr8PqCQLSr8PqCQKfsIjIDgKfsLDIDgKfsLTIDgKfsLzIDgKfsKDIDgKfsOTLDgKfsOjLDgKP3+6lAgLwkp3BDALwksHADAKKxdH8DALukuXADAK7q7GGCAKM54rGBgKMk/3ADALf9dMTM1gxLrYgWuhI/FJkOEsD206ylsI=")
	
	/*for k,va :=range v1{
		fmt.Println(k,va)
	}*/
	body1:=strings.NewReader(v1.Encode())
	fmt.Println(v1.Encode())
	r1,err:=http.NewRequest("POST",resulturl,body1)
	if err != nil {
		log.Fatal(err)
	}
	r1.Header.Add("Referer", "http://xk1.ahu.cn/xscjcx.aspx?xh=P71514011&xm=%C7%C7%D3%A2%BD%DC&gnmkdm=N121605")
	r1.Header.Add("Content-Type", "application/x-www-form-urlencoded")
	r1.Header.Add("User-Agent", "Mozilla/5.0 (Windows NT 10.0; Win64; x64) AppleWebKit/537.36 (KHTML, like Gecko) Chrome/59.0.3071.104 Safari/537.36")
	r1.Header.Add("Host","xk1.ahu.cn")
	r1.Header.Add("Origin","http://xk1.ahu.cn")
	/*for k,va :=range r1.Header{
		fmt.Println(k,va)
	}*/
	c.CheckRedirect= func(req *http.Request, via []*http.Request) error {
		  return http.ErrUseLastResponse
		}
	response1, err:= c.Do(r1)
	//fmt.Println(c.)
	fmt.Println(response.Status)
	if err != nil {
		log.Fatal(err)
	}
	file3, err := os.OpenFile("spider.txt", os.O_RDWR|os.O_CREATE|os.O_TRUNC, os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
	defer file3.Close()
	body2, _ := ioutil.ReadAll(response1.Body)
	num, err := file3.Write(body2)
	if num != len(body2) {
		log.Fatal(err)
	}
	fmt.Println(response1.Status)
}
