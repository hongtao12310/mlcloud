FROM golang:1.8.3

MAINTAINER Hongtao Zhang

WORKDIR /go/src/github.com/deepinsight/mlcloud

COPY ./src ./src
COPY ./vendor ./vendor


RUN cd src && \
    go build -o mlcloud && \
    mv mlcloud $GOPATH/bin

EXPOSE 9090

ENTRYPOINT ["mlcloud", "-port", "9090"]
