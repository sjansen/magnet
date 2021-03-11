FROM golang:1.16-alpine as builder
ARG MODULE=github.com/sjansen/magnet
RUN apk add --update \
    ca-certificates \
    git
ADD go.mod go.sum main.go /go/src/${MODULE}/
RUN cd /go/src/${MODULE} && \
    go mod download

ARG GITSHA="(missing)"
ARG TIMESTAMP="(missing)"
ADD cmd      /go/src/${MODULE}/cmd
ADD internal /go/src/${MODULE}/internal
RUN cd /go/src/${MODULE} && \
    echo GITSHA=${GITSHA} && \
    echo TIMESTAMP="${TIMESTAMP}" && \
    CGO_ENABLED=0 GOOS=linux \
    go build \
        -o /lambda-handler \
        -a -installsuffix cgo \
        -ldflags="-s -w -X ${MODULE}/internal/build.GitSHA=${GITSHA} -X '${MODULE}/internal/build.Timestamp=${TIMESTAMP}'" \
        ./cmd/move

FROM scratch
COPY --from=builder /etc/ssl/certs/ca-certificates.crt /etc/ssl/certs/ca-certificates.crt
COPY --from=builder /lambda-handler /lambda-handler
EXPOSE 8000
WORKDIR /
ENTRYPOINT ["/lambda-handler"]
