FROM ubuntu:19.04

ENV DEBIAN_FRONTEND noninteractive
ENV GO111MODULE on

RUN apt-get update && apt-get install -y valgrind openssh-server iperf3 iputils-ping vim kcachegrind snapd

RUN snapctl start

RUN snap install go --classic 

ADD ./client/ /p2p-tests

WORKDIR /p2p-tests

RUN go get || true

RUN go build


ENTRYPOINT ["/bin/bash"]
