FROM golang:1.12.5-stretch

ENV DEBIAN_FRONTEND noninteractive
ENV GO111MODULE on

ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN go get || true

RUN go build

RUN apt-get update && apt-get install -y valgrind openssh-server iperf3 iputils-ping vim kcachegrind

ENTRYPOINT ["/bin/bash"]
