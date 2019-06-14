FROM ubuntu:bionic

WORKDIR /nodejs 

RUN apt-get update && \
      apt-get install -y build-essential wget

RUN wget -O node.tar.gz https://nodejs.org/dist/v10.16.0/node-v10.16.0.tar.gz
RUN tar xf node.tar.gz

WORKDIR /nodejs/node-v10.16.0

RUN ./configure --with-dtrace
RUN make -j 32
