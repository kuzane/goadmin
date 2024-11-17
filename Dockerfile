FROM golang:1.22.5-alpine3.20  as application
ENV WORKDIR /app
# ENV GOPROXY https://goproxy.cn
WORKDIR $WORKDIR
COPY . $WORKDIR
RUN cd /app && CGO_ENABLED=0 GOOS=linux GOARCH=amd64 go build -ldflags "-s -w" cmd/main.go

FROM node:18 as frontend
ENV WORKDIR /app
WORKDIR $WORKDIR
COPY web $WORKDIR
# RUN cd /app && npm install --registry=https://registry.npmmirror.com && npm run build:prod
RUN cd /app && npm install && npm run build:prod

FROM alpine:latest
ENV WORKDIR /app
WORKDIR $WORKDIR

# RUN sed -i 's/dl-cdn.alpinelinux.org/mirrors.aliyun.com/g' /etc/apk/repositories && \
#     apk update && \
#     mkdir -pv web/dist &&\
#   	apk --no-cache add tzdata && \
#   	cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
#   	apk del tzdata && \
#   	rm -rf /var/cache/apk/*
RUN apk update && \
    mkdir -pv web/dist &&\
  	apk --no-cache add tzdata && \
  	cp -f /usr/share/zoneinfo/Asia/Shanghai /etc/localtime && \
  	apk del tzdata && \
  	rm -rf /var/cache/apk/*

COPY --from=application /app/main .
COPY --from=frontend /app/dist ./web/dist

EXPOSE 8000
CMD ["/app/main"]
