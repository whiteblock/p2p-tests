FROM golang:1.12.3-stretch

ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN go get || true

RUN go build

ENTRYPOINT ["/bin/bash"]
