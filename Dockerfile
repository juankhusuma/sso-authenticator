############################
# STEP 1 build executable binary
############################
FROM golang:alpine AS builder

RUN apk update && apk add --no-cache git
WORKDIR $GOPATH/src/mypackage/myapp/
COPY .env .
COPY . .
# RUN go get -d -v
RUN go build -o /go/bin/hello ./cmd/http-server

FROM centos
COPY --from=builder /go/bin/hello /go/bin/hello
EXPOSE 3000
ENTRYPOINT ["/go/bin/hello"]