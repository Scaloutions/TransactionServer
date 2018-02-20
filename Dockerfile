FROM golang:latest

RUN mkdir -p /src

RUN go get "github.com/golang/glog"

RUN go get "github.com/gin-gonic/gin"

WORKDIR /src

ADD . /src

RUN go build ./server.go

CMD [ "./server" ]

