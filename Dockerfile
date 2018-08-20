FROM alpine:3.6

RUN apk add --update ca-certificates
RUN mkdir /lib64 && ln -s /lib/libc.musl-x86_64.so.1 /lib64/ld-linux-x86-64.so.2
RUN mkdir /root/.kube
COPY kube_config /root/.kube/config
COPY example_duncan.yml /root/.duncan.yml
