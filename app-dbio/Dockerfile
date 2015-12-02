# Dockerfile

FROM golang:1.5.1
MAINTAINER pottava

RUN go get -u github.com/justinas/alice

RUN go get -u github.com/aws/aws-sdk-go/aws
RUN go get -u github.com/aws/aws-sdk-go/service/dynamodb

LABEL jp.co.supinf.works.application="golang-micro-services-dbio" \
      jp.co.supinf.works.license="MIT"
