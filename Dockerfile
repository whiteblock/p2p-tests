FROM ubuntu:19.04

ENV DEBIAN_FRONTEND noninteractive
ENV GO111MODULE on

ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN apt-get update && apt-get install -y valgrind openssh-server iperf3 iputils-ping vim kcachegrind snapd

RUN snap install go --classic 

RUN go get || true

RUN go build


ENTRYPOINT ["/bin/bash"]
