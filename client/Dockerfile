FROM golang:1.12.3-stretch

ADD . /libp2p

WORKDIR /libp2p

RUN go get || true

RUN go build

ENTRYPOINT ["/bin/bash"]