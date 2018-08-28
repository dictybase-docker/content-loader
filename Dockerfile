FROM golang:1.10.4-alpine3.8
MAINTAINER Siddhartha Basu <siddhartha-basu@northwestern.edu>
RUN apk add --no-cache git build-base \
    && go get github.com/golang/dep/cmd/dep
RUN mkdir -p /go/src/github.com/dictybase-docker/content-loader
WORKDIR /go/src/github.com/dictybase-docker/content-loader
COPY Gopkg.* *.go ./
RUN dep ensure \
    && go build -o app

FROM alpine:3.8
RUN apk --no-cache add ca-certificates
COPY --from=0 /go/src/github.com/dictybase-docker/content-loader/app /usr/local/bin/
ENTRYPOINT ["/usr/local/bin/app"]
