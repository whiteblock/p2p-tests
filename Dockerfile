FROM golang:1.12.3-stretch

ENV DEBIAN_FRONTEND noninteractive

ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN go get || true

RUN go build

RUN apt-get update && apt-get install -y valgrind

ENTRYPOINT ["/bin/bash"]
