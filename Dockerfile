FROM alpine
MAINTAINER Elliott Polk <elliott@tkwcafe.com>

RUN rm -rf /var/cache/apk/* && \
    mkdir -p /var/lib/confgr

WORKDIR /var/lib/confgr

COPY build/bin/* /usr/bin/
ENTRYPOINT ["/usr/bin/confgr"]
