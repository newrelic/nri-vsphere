FROM golang:1.14

WORKDIR /go/src/nri-vsphere
COPY . .
RUN make deps

RUN apt-get update && \
    apt-get install docker-compose -qq > /dev/null && \
    rm -rf /var/lib/apt/lists/*

CMD /etc/init.d/docker start && make test compile