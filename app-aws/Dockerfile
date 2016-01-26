# Dockerfile

FROM pottava/golang:1.5
MAINTAINER pottava

RUN go get -u github.com/justinas/alice

RUN go get -u github.com/aws/aws-sdk-go/aws
RUN go get -u github.com/aws/aws-sdk-go/service/dynamodb
RUN go get -u github.com/aws/aws-sdk-go/service/ec2

LABEL jp.co.supinf.works.application="golang-microservices-aws" \
      jp.co.supinf.works.license="MIT"
