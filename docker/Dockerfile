FROM golang:alpine

ENV TIME_ZONE=Asia/Shanghai

RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.ustc.edu.cn/g' /etc/apk/repositories
RUN   apk update \
  && apk add --no-cache shadow git bash perl openssh vim tar curl python3 python3-dev py-pip gcc g++ libcurl make\
  && usermod -s /bin/bash root \
  && rm -rf /var/cache/apk/* 


RUN \
  mkdir -p /go/src/infini.sh/ \
  && apk add --no-cache tzdata \
  && echo "${TIME_ZONE}" > /etc/timezone \ 
  && ln -sf /usr/share/zoneinfo/${TIME_ZONE} /etc/localtime 

WORKDIR /go/src/infini.sh/
COPY ../ /go/src/infini.sh/proxy
COPY ./entrypoint.sh /

EXPOSE 2900

#CMD npm run dev
ENTRYPOINT ["/entrypoint.sh"]