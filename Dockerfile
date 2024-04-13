FROM golang:1.22.2-alpine as gobulidIso
COPY ./ /go/src/wp-go
WORKDIR /go/src/wp-go
RUN go build -ldflags "-w" -tags netgo -o wp-go app/cmd/main.go

FROM alpine:latest
WORKDIR /opt/wp-go
COPY --from=gobulidIso /go/src/wp-go/wp-go ./
ENTRYPOINT ["./wp-go"]