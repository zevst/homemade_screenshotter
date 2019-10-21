FROM golang:1.13-alpine AS imgbld
LABEL stage=builder
RUN apk update && apk add openssh git gcc libc-dev ca-certificates && mkdir -p /go/src/
COPY . /go/src/hmsc
RUN cd /go/src/hmsc/ && GO111MODULE=on go build -a -o hmsc

FROM alpine:latest
RUN apk update && apk add ca-certificates && mkdir -p /go/src/hmsc
COPY --from=imgbld /go/src/hmsc/hmsc /go/src/hmsc
COPY ./.env /go/src/hmsc
WORKDIR /go/src/hmsc/
ENTRYPOINT ["./hmsc"]
