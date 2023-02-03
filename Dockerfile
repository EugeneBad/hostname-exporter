FROM golang:1.17.6  AS builder
WORKDIR /go/build
COPY ./cmd/ ./
RUN go mod tidy
RUN CGO_ENABLED=0 GOOS=linux go build -o hostname-exporter .

FROM alpine:3.16.0  
WORKDIR /root/
COPY --from=builder /go/build/hostname-exporter ./
CMD ["./hostname-exporter"]