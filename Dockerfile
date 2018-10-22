FROM golang:1.8-alpine

EXPOSE 4000

WORKDIR /go/
COPY src/ src/

RUN go build node_gen
CMD ["./node_gen"]
