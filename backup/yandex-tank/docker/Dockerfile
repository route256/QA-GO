# Yandex.Tank
#
# VERSION 0.0.4

FROM ubuntu:focal
MAINTAINER Yandex Load Team <load-public@yandex-team.ru>
# Version for desription
ARG VERSION
# You may specify tag instead of branch to build for specific tag
ARG BRANCH=master

LABEL Description="Fresh Yandex.Tank from github master branch with phantom" \
    Vendor="Yandex" \
    maintainer="direvius@yandex-team.ru" \
    YandexTank.version="${VERSION}" \
    Telegraf.version="${TELEGRAF_VERSION}"

RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update -q && \
    apt-get install --no-install-recommends -yq \
        sudo     \
        vim      \
        wget     \
        curl     \
        less     \
        iproute2 \
        software-properties-common \
        telnet   \
        atop     \
        openssh-client \
        git            \
        gpg-agent      \
        python3-pip  && \
    add-apt-repository ppa:yandex-load/main -y && \
    apt-get update -q && \
    apt-get install -yq \
        phantom \
        phantom-ssl && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/* /tmp/* /var/tmp/*

ENV TELEGRAF_VERSION=1.19.1-1
RUN curl -sL https://repos.influxdata.com/influxdb.key | gpg --import && \
    wget https://dl.influxdata.com/telegraf/releases/telegraf_${TELEGRAF_VERSION}_amd64.deb.asc && \
    wget https://dl.influxdata.com/telegraf/releases/telegraf_${TELEGRAF_VERSION}_amd64.deb && \
    gpg --batch --verify telegraf_${TELEGRAF_VERSION}_amd64.deb.asc telegraf_${TELEGRAF_VERSION}_amd64.deb && \
    dpkg -i telegraf_${TELEGRAF_VERSION}_amd64.deb && \
    rm -f telegraf_${TELEGRAF_VERSION}_amd64.deb*

COPY . /home/tank
WORKDIR /home/tank
RUN ls -lah
#    pip3 install https://api.github.com/repos/yandex/yandex-tank/tarball/${BRANCH} && \
ENV BUILD_DEPS="python3-dev build-essential gfortran libssl-dev libffi-dev"
RUN export DEBIAN_FRONTEND=noninteractive && \
    apt-get update && \
    apt-get install -yq --no-install-recommends ${BUILD_DEPS} && \
    pip3 install --upgrade setuptools && \
    pip3 install --upgrade pip && \
    pip3 install . && \
    apt-get autoremove -y ${BUILD_DEPS} && \
    apt-get clean && \
    rm -rf /var/lib/apt/lists/* /var/cache/apt/archives/* /tmp/* /var/tmp/* /root/.cache/*

RUN curl -L -v -o pandora https://github.com/yandex/pandora/releases/download/v0.4.0/pandora_0.4.0_linux_amd64 && \
    chmod +x ./pandora && \
    mv ./pandora /usr/local/bin/

COPY docker/files/bashrc /root/.bashrc
COPY docker/files/inputrc /root/.inputrc

VOLUME ["/var/loadtest"]
WORKDIR /var/loadtest
ENTRYPOINT ["bash", "-c"]
#ENTRYPOINT ["/usr/local/bin/yandex-tank"]
