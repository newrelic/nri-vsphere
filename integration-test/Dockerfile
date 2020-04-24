FROM golang:1.10 as builder
ARG CGO_ENABLED=0

#Compiling mock
WORKDIR /go/src/github.com/vmware/
RUN git clone https://github.com/vmware/govmomi.git
WORKDIR /go/src/github.com/vmware/govmomi/vcsim
RUN make all

WORKDIR /go/src/github.com/newrelic/nri-vsphere/
COPY . .

#Compiling nri-vsphere, it will avaible at /go/src/github.com/newrelic/nri-vsphere/bin/linux/nri-vsphere
RUN make compile-linux

