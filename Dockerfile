FROM golang:1.14

# needed for running integration tests on docker
RUN apt-get update && \
    apt-get install docker-compose -qq > /dev/null && \
    rm -rf /var/lib/apt/lists/*

WORKDIR /go/src/nri-vsphere
COPY . .
RUN make deps

CMD /etc/init.d/docker start && make test compile