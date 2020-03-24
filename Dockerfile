FROM golang

RUN go get github.com/tencentcloud/tencentcloud-sdk-go && \
    go get github.com/julienschmidt/httprouter && \
    go get golang.org/x/crypto/ssh && \
    go get github.com/jmoiron/sqlx && \
    go get github.com/lib/pq && \
    go get github.com/gomodule/redigo/redis

WORKDIR /go/src

RUN git clone https://github.com/dkzhang/RAS.git #20200318

WORKDIR /go/src/RAS/webApiServer

CMD go run ./server.go