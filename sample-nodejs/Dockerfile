# TODO Improve with https://github.com/nodejs/docker-node/blob/5d8827883ba24066cec73f6d6c3e70c3ec4e1b28/12/alpine/Dockerfile
FROM ubuntu:18.04

ENV NODE_VERSION 12.4.0

WORKDIR /nodejs 

RUN apt-get update && \
      apt-get install -y build-essential wget python systemtap-sdt-dev

RUN wget -O node.tar.gz "https://nodejs.org/dist/v$NODE_VERSION/node-v10.16.0.tar.gz"
RUN tar xf node.tar.gz

WORKDIR "/nodejs/node-v$NODE_VERSION"

RUN ./configure --with-dtrace
RUN make -j 32

WORKDIR /code
COPY . .

CMD ["sh", "-c", "node ./app.js"]


