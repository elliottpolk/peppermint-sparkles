FROM alpine
MAINTAINER Elliott Polk <elliott@tkwcafe.com>

RUN rm -rf /var/cache/apk/*

COPY confgr /usr/bin/confgr
RUN mkdir -p /var/lib/confgr

WORKDIR /var/lib/confgr
ENTRYPOINT ["/usr/bin/confgr", "server"]