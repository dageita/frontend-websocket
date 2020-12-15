FROM golang:alpine
COPY echo.go .
COPY exec.go .
COPY vendor /go/src
RUN go build -o websocket *.go

FROM alpine
COPY --from=0 /go/websocket .
CMD ["./websocket"]
