FROM golang:1.10-stretch

RUN mkdir -p /go/src/github.com/pagarme/teleport/
WORKDIR /go/src/github.com/pagarme/teleport/
RUN go get -u github.com/jteeuwen/go-bindata/...

ADD . /go/src/github.com/pagarme/teleport/

RUN make install
