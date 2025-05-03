FROM golang:1.24.0-alpine3.21 as app-builder
WORKDIR /go/src/app
RUN apk add alpine-sdk
COPY go.mod go.sum ./
RUN go mod tidy
COPY . .
RUN go build -o /everything

FROM alpine:3.16
RUN apk add tzdata
COPY --from=app-builder /everything /everything
ENTRYPOINT ["/everything"]
