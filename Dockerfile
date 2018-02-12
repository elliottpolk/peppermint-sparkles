FROM debian:latest
MAINTAINER Elliott Polk <benjamin_elliott_polk@manulifeam.com>

RUN rm -rf /var/cache/apk/* && \
    mkdir -p /var/lib/peppermint-sparkles

WORKDIR /var/lib/peppermint-sparkles

COPY build/bin/* /usr/bin/
ENTRYPOINT ["/usr/bin/psparkles"]
