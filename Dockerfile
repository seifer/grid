FROM golang:1.9.3 as build

ENV GOOS linux
ENV GOARCH amd64
ENV CGO_ENABLED 0

ADD . /go/src/grid
WORKDIR /go/src/grid

RUN go build -a -o app

FROM alpine:3.6

MAINTAINER Timur Nurutdinov <timur.nurutdinov@lamoda.ru>

RUN apk add --no-cache --virtual tzdata
RUN apk add --no-cache --virtual ca-certificates

COPY --from=build /go/src/grid/app /bin/app

EXPOSE 8080

CMD /bin/app