FROM golang:latest as build_env

COPY . $GOPATH/src/github.com/romanyx/integral_db
WORKDIR $GOPATH/src/github.com/romanyx/integral_db/cmd/server

ENV GO111MODULE=on
RUN go get -v
ENV CGO_ENABLED=0 GOOS=linux GOARCH=amd64
RUN go build -o /go/bin/db

FROM alpine
LABEL maintainer="Roman Budnikov work@romanyx.ru"

COPY --from=build_env /go/bin/db /bin/db

ENTRYPOINT ["/bin/db"]
