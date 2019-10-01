FROM ubuntu:18.04

ARG VERSION

RUN apt-get update && apt-get install wget -y
WORKDIR /gopxy
RUN wget -c https://github.com/xxxsen/gopxy/releases/download/$VERSION/binary.tar.gz -O - | tar -xz
WORKDIR binary

EXPOSE 8080
VOLUME /etc/gopxy

RUN ls -alh

CMD ./gopxy_linux --config=/etc/gopxy/config.json

