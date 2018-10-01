FROM golang:1.8

WORKDIR .
RUN go build node_gen
CMD ["./node_gen"]
