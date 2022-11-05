FROM golang:latest
COPY ./ /go/src/wp-go
WORKDIR /go/src/wp-go
#ENV GOPROXY="https://goproxy.cn"
RUN go build -tags netgo

FROM alpine:latest
WORKDIR /opt/wp-go
COPY --from=0 /go/src/wp-go/wp-go ./
ENTRYPOINT ["./wp-go"]