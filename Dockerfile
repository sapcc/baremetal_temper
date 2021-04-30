FROM golang:1.13.4-alpine3.10 as builder
WORKDIR /go/src/github.com/sapcc/baremetal_temper
RUN apk add --no-cache make git
COPY . .
ARG VERSION
RUN make all

FROM alpine:3.12
LABEL maintainer="Stefan Hipfel <stefan.hipfel@sap.com>"
LABEL source_repository="https://github.com/sapcc/baremetal_temper"

RUN apk add --no-cache curl
RUN curl -Lo /bin/dumb-init https://github.com/Yelp/dumb-init/releases/download/v1.2.2/dumb-init_1.2.2_amd64 \
	&& chmod +x /bin/dumb-init \
	&& dumb-init -V
COPY --from=builder /go/src/github.com/sapcc/baremetal_temper/bin/linux/temper /usr/local/bin/
COPY /etc/ /etc/
ENTRYPOINT ["dumb-init", "--"]
CMD ["temper"]