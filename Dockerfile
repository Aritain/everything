FROM golang:1.24.0-alpine3.21 as app-builder
WORKDIR /go/src/app
COPY . .
RUN apk add alpine-sdk
RUN go mod init everything
RUN go mod tidy
RUN go build -o /everything

FROM alpine:3.16
COPY --from=app-builder /everything /everything
ENTRYPOINT ["/everything"]
