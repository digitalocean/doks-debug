# match doks-debug version with DOKS worker node image version for kernel
# tooling compatibility reasons
FROM debian:12-slim

# Specify the version of crictl to install
ARG CRICTL_VERSION="v1.31.1"

WORKDIR /root

# use same dpkg path-exclude settings that come by default with ubuntu:focal
# image that we previously used
RUN echo 'path-exclude=/usr/share/locale/*/LC_MESSAGES/*.mo' >> /etc/dpkg/dpkg.cfg.d/excludes
RUN echo 'path-exclude=/usr/share/doc/*' >> /etc/dpkg/dpkg.cfg.d/excludes
RUN echo 'path-include=/usr/share/doc/*/copyright' >> /etc/dpkg/dpkg.cfg.d/excludes
RUN echo 'path-include=/usr/share/doc/*/changelog.Debian.*' >> /etc/dpkg/dpkg.cfg.d/excludes

RUN apt-get update -qq && \
    apt-get install -y apt-transport-https \
                       ca-certificates \
                       software-properties-common \
                       httping \
                       man \
                       man-db \
                       vim \
                       screen \
                       curl \
                       gnupg \
                       atop \
                       htop \
                       dstat \
                       jq \
                       dnsutils \
                       tcpdump \
                       termshark \
                       traceroute \
                       iputils-ping \
                       iptables \
                       net-tools \
                       ncat \
                       iproute2 \
                       strace \
                       lsof \
                       telnet \
                       openssl \
                       psmisc \
                       dsniff \
                       mtr-tiny \
                       conntrack \
                       llvm-13 llvm-13-tools \
                       wget \
                       watch \
                       bpftool

# Install crictl
RUN wget https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz && \
    tar zxvf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz -C /usr/local/bin && \
    rm -f crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

# Specify the default image endpoint for crictl
RUN echo 'runtime-endpoint: unix:///run/containerd/containerd.sock' >> /etc/crictl.yaml
RUN echo 'image-endpoint: unix:///run/containerd/containerd.sock' >> /etc/crictl.yaml
RUN echo 'timeout: 2' >> /etc/crictl.yaml

CMD [ "/bin/bash" ]
