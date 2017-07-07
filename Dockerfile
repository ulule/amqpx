FROM golang:1.8.3

ADD . /go/src/github.com/ulule/amqpx
WORKDIR /go/src/github.com/ulule/amqpx

RUN go get -u github.com/alecthomas/gometalinter && gometalinter --install

CMD make ci && make lint
