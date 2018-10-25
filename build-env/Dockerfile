FROM ubuntu:16.04
MAINTAINER Dmitrii Beliakov <dmitriy.b11@gmail.com>

RUN apt-get update && apt-get install -y wget curl
RUN wget https://dl.google.com/go/go1.11.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go*.tar.gz
ENV GOPATH /usr/local/go
ENV PATH=$PATH:/usr/local/go/bin

RUN apt-get update && apt-get install -y wget apt-transport-https ca-certificates curl software-properties-common
RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add -
RUN add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable"
RUN apt-get update && apt-get install -y docker-ce python fontconfig

RUN apt-get update && apt-get install -y openssh-client python-pip python-crypto bash rsync
RUN pip install fabric==1.12.0

RUN curl -sL https://deb.nodesource.com/setup_8.x | bash -
RUN apt-get install -y nodejs

