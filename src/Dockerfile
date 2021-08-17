FROM golang:1.15-alpine
COPY main.go .
COPY exec.go .
COPY vendor /go/src
RUN go build -o main *.go

FROM alpine
COPY --from=0 /go/main .
CMD ["./main"]
