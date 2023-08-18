FROM alpine:latest

WORKDIR /FileServer
COPY FileServer ./FileServer

RUN apk update \
    && apk add --no-cache tzdata \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone \
    && chmod +x ./FileServer \
    && mkdir -p /data/aria2 \
    && chmod -R 766 /data/aria2

EXPOSE 5212
VOLUME ["/FileServer/uploads", "/FileServer/avatar", "/data"]

ENTRYPOINT ["./FileServer"]
