FROM golang:1.16 as builder
ARG CGO_ENABLED=0
WORKDIR /go/src/github.com/newrelic/nri-vsphere
COPY . .

RUN make clean compile

FROM newrelic/infrastructure:latest
COPY --from=builder /go/src/github.com/newrelic/nri-vsphere/bin /var/db/newrelic-infra/newrelic-integrations/bin
COPY --from=builder /go/src/github.com/newrelic/nri-vsphere/tools/vsphere-config.yml /etc/newrelic-infra/integrations.d/vsphere-config.yml
COPY --from=builder /go/src/github.com/newrelic/nri-vsphere/tools/vsphere-performance.metrics /etc/newrelic-infra/integrations.d/vsphere-performance.metrics

ARG NRIA_LICENSE_KEY
ARG NRIA_STAGING=false

ENV NRIA_LICENSE_KEY $NRIA_LICENSE_KEY
ENV NRIA_STAGING $NRIA_STAGING
ENV NRIA_VERBOSE 1
