FROM alpine
ENV VERSION 2.0.0
MAINTAINER Elliott Polk <elliott@tkwcafe.com>

RUN rm -rf /var/cache/apk/*

COPY confgr /usr/bin/confgr
RUN mkdir -p /var/lib/confgr

WORKDIR /usr/bin

CMD ["/usr/bin/confgr", "server"]
