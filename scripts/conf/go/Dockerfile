FROM golang:1.10.0-stretch

ADD . /go/src/github.com/ulule/amqpx
WORKDIR /go/src/github.com/ulule/amqpx

RUN go get -u github.com/alecthomas/gometalinter \
    && gometalinter --install \
    && go get -d -u github.com/stretchr/testify/require \
    && go get -d -u github.com/pkg/errors \
    && go get -d -u github.com/streadway/amqp

CMD go version
