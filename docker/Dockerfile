FROM golang:alpine

RUN  apk update \
  && apk add --no-cache shadow git bash perl openssh vim tar curl python3 python3-dev py-pip gcc g++ libcurl make\
  && usermod -s /bin/bash root \
  && rm -rf /var/cache/apk/*

RUN mkdir -p /go/src/infini.sh

WORKDIR /go/src/infini.sh/
COPY ./ /go/src/infini.sh/proxy
COPY ./docker/entrypoint.sh /

RUN cd proxy && make build

FROM alpine:latest
WORKDIR /
COPY  --from=0 /go/src/infini.sh/proxy/bin/proxy /
EXPOSE 2900
CMD ["/proxy"]