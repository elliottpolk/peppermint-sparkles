FROM openjdk:8-jdk-slim-buster
LABEL maintainer Elliott Polk <benjamin_elliott_polk@manulifeam.com>

RUN mkdir -p /root/.gnupg \
    && printf 'use-agent\npinentry-mode loopback' > /root/.gnupg/gpg.conf \
    && printf 'allow-loopback-pinentry' > /root/.gnupg/gpg-agent.conf

RUN  apt-get update && \
    apt-get install -y \
		unzip \
		jq \
		gnupg \
		curl \
		wget \
		fastjar \
		nodejs \
		npm

# install cf tools
# RUN wget -O /tmp/cf-cli.deb "https://cli.run.pivotal.io/stable?release=debian64&source=github" && \
RUN wget --no-check-certificate -O /tmp/cf-cli.deb "https://cli.run.pivotal.io/stable?release=debian64&source=github" && \
    dpkg -i /tmp/cf-cli.deb && \
    apt-get install -f

ENV GRADLE_HOME /opt/gradle
ENV GRADLE_VERSION 6.5

ARG GRADLE_CHECKSUM=23e7d37e9bb4f8dabb8a3ea7fdee9dd0428b9b1a71d298aefd65b11dccea220f
RUN set -o errexit -o nounset \
	&& mkdir -p /opt \
	&& echo "Downloading gradle" \
	&& wget --no-verbose --output-document=gradle.zip "https://services.gradle.org/distributions/gradle-${GRADLE_VERSION}-bin.zip" \
	&& echo "Checksum validation" \
	&& echo "${GRADLE_CHECKSUM} *gradle.zip" | sha256sum -c - \
	&& echo "Installing gradle" \
	&& unzip gradle.zip \
	&& rm -v gradle.zip \
	&& mv "gradle-${GRADLE_VERSION}/" "${GRADLE_HOME}/" \
	&& ln -s "${GRADLE_HOME}/bin/gradle" /usr/bin/gradle
