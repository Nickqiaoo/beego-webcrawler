FROM golang:1.8

WORKDIR /go/src/beego-webcrawler

COPY . .
COPY github.com $GOPATH/src/github.com
COPY golang.org $GOPATH/src/github.com

RUN  go build main.go 

CMD ["./main"]

