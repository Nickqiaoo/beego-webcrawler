package controllers

import (
	"net/http"
	
	"github.com/PuerkitoBio/goquery"
	"github.com/axgle/mahonia"
)

// Result 学分结果
type Result struct {
	Num  string
	Name string
	Xf   string
	Res  map[string]string
	Jd   string
}

// Info 学号，姓名
type Info struct {
	Name string
	Num  string
}

// Grade 成绩
type Grade struct{
	Num  string
	Name string
	Graderesult [][]string 
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
