FROM golang:1.14

WORKDIR /go/src/nri-vsphere
COPY . .

RUN apt update && \
    apt install docker-compose -y && \
    rm -rf /var/lib/apt/lists/*

RUN make deps