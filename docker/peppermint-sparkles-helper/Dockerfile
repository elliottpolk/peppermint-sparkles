FROM debian:stretch-slim
LABEL maintainer Elliott Polk <benjamin_elliott_polk@manulife.com>

RUN mkdir -p /root/.gnupg \
    && printf 'use-agent\npinentry-mode loopback' > /root/.gnupg/gpg.conf \
    && printf 'allow-loopback-pinentry' > /root/.gnupg/gpg-agent.conf

RUN apt-get update && \
    apt-get install -y jq gnupg curl wget fastjar

# install cf tools
# RUN wget -O /tmp/cf-cli.deb "https://cli.run.pivotal.io/stable?release=debian64&source=github" && \
RUN wget --no-check-certificate -O /tmp/cf-cli.deb "https://cli.run.pivotal.io/stable?release=debian64&source=github" && \
    dpkg -i /tmp/cf-cli.deb && \
    apt-get install -f


