FROM golang:alpine as app-builder
WORKDIR /go/src/app
COPY . .
RUN apk add alpine-sdk
RUN go mod init client
RUN go mod tidy
RUN go build -o /f1_tracker

FROM alpine:3.16
COPY --from=app-builder /f1_tracker /f1_tracker
ENTRYPOINT ["/f1_tracker"]
