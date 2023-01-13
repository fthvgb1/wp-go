FROM golang:latest as gobulidIso
COPY ./ /go/src/wp-go
WORKDIR /go/src/wp-go
ENV GOPROXY="https://goproxy.cn"
RUN go build -ldflags "-w" -tags netgo -o wp-go internal/cmd/main.go

FROM alpine:latest
WORKDIR /opt/wp-go
COPY --from=gobulidIso /go/src/wp-go/wp-go ./
ENTRYPOINT ["./wp-go"]