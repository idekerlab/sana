FROM golang:latest

RUN go get -u google.golang.org/grpc

COPY sana-bin /bin/sana
COPY proto    /go/src/github.com/sana/proto
COPY main.go  /go/src/github.com/sana/main.go

CMD ["go","run","/go/src/github.com/sana/main.go"]
