FROM golang:1.8

EXPOSE 4000

WORKDIR /go/
ADD . .

RUN go build node_gen
CMD ["./node_gen"]
