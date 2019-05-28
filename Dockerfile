FROM ubuntu:19.04

ENV DEBIAN_FRONTEND noninteractive
ENV GO111MODULE on

RUN apt-get update && apt-get install -y valgrind openssh-server iperf3 iputils-ping vim kcachegrind snapd

RUN wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
RUN tar -xvf go1.12.2.linux-amd64.tar.gz && mv go /usr/local/
RUN ln -s /usr/local/go/bin/go /usr/local/bin/go
RUN apt-get install -y git
ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN go get || true

RUN go build


ENTRYPOINT ["/bin/bash"]
