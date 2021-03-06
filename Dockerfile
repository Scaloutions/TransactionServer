FROM golang:1.9.2-alpine3.6 AS build

RUN apk add --no-cache git

RUN mkdir -p /app

ADD . /app/

WORKDIR /app

#RUN go get ./

RUN go get "github.com/golang/glog"

RUN go get "github.com/gin-gonic/gin"

RUN go get "github.com/go-sql-driver/mysql"

RUN go get "github.com/joho/godotenv"

RUN go get "github.com/garyburd/redigo/redis"


RUN go build -o server .

CMD [ "/app/server", "-stderrthreshold=ERROR", "-log_dir=./logs" ]

EXPOSE 9090
