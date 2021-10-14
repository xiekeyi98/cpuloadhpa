# 编译用的
FROM golang:1.17.2 as builder 

WORKDIR /src/cpuloadhpa
# 拷贝源文件过去
COPY . .
# 编译

RUN go build -o cpuloadhpa ./main.go


FROM alpine:3.14

# 维护者标签
LABEL MAINTAINER = "keyixie@tencent.com"

RUN apk update && apk add --no-cache curl bash && apk add -U tzdata 
RUN rm -f /etc/localtime \
    && cp /usr/share/zoneinfo/Asia/Shanghai /etc/localtime \
    && echo "Asia/Shanghai" > /etc/timezone
# 暴露端口
#EXPOSE 8088

# 镜像工作目录
WORKDIR /app 
COPY --from=builder /src/cpuloadhpa/cpuloadhpa .


ENTRYPOINT [ "/app/cpuloadhpa"]

