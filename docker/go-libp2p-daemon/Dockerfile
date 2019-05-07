FROM ubuntu:18.04
RUN apt-get update &&\
    DEBIAN_FRONTEND=noninteractive apt-get install -y apt-utils expect git \
    software-properties-common inetutils-tools wget ca-certificates curl \
    libssl-dev make iperf3 iputils-ping tmux build-essential git-extras

RUN wget https://dl.google.com/go/go1.12.2.linux-amd64.tar.gz
RUN tar -xvf go1.12.2.linux-amd64.tar.gz
RUN mv go /usr/local
ENV GOROOT=/usr/local/go
ENV GOPATH /root/go
ENV PATH=${PATH}:${GOPATH}/bin:${GOROOT}/bin:${PATH}

RUN git clone https://github.com/libp2p/go-libp2p-daemon.git
WORKDIR /go-libp2p-daemon
RUN go get ./...
RUN go install ./...

WORKDIR /go-libp2p-daemon/p2pd
RUN go build
RUN mv p2pd /usr/bin

ENTRYPOINT ["/bin/bash"]