FROM golang:1.15-alpine as builder
ARG MODULE=github.com/sjansen/magnet
ARG GITSHA="(missing)"
ARG TIMESTAMP="(missing)"

RUN apk --update add ca-certificates git
ADD go.mod go.sum main.go /go/src/${MODULE}/
RUN cd /go/src/${MODULE} && \
    go mod download
ADD internal /go/src/${MODULE}/internal
RUN cd /go/src/${MODULE} && \
    echo GITSHA=${GITSHA} && \
    echo TIMESTAMP=${TIMESTAMP} && \
    CGO_ENABLED=0 GOOS=linux \
    go build \
        -a -installsuffix cgo \
        -ldflags="-s -w -X ${MODULE}/internal/build.GitSHA=${GITSHA} -X ${MODULE}/internal/build.Timestamp=${TIMESTAMP}" \
        -o /app

FROM scratch
COPY --from=builder /app /app
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["/app"]
