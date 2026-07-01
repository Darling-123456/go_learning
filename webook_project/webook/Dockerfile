#基础镜像
FROM ubuntu:20.04
#把编译后的打包进来这个镜像 放到工作目录/app
COPY webook /app/webook
#设定我的工作目录
WORKDIR /app

ENTRYPOINT ["/app/webook"]
#LABEL authors="darling-ge"

#ENTRYPOINT ["top", "-b"]