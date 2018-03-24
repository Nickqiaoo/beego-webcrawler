FROM golang:1.8

WORKDIR /go/src/beego-webcrawler

COPY . .
COPY golang.org $GOPATH/src/golang.org

RUN  go get github.com/PuerkitoBio/goquery \
	&& go get github.com/astaxie/beego \
	&& go get github.com/axgle/mahonia \
	&& go build main.go 

CMD ["./main"]

