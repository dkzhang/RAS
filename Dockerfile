FROM golang

RUN go env -w GO111MODULE=on && \
    go env -w GOPROXY=https://goproxy.cn,direct

WORKDIR /go/src

RUN git clone https://github.com/dkzhang/RAS.git #20200716-2132

WORKDIR /go/src/RAS/webApiServer

RUN go build ./server.go

CMD ./server