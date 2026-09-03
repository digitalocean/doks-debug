# match doks-debug version with DOKS worker node image version for kernel
# tooling compatibility reasons
FROM debian:13-slim

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
                       llvm-19 llvm-19-tools \
                       wget \
                       watch \
                       bpftool \
                       xdp-tools

# Install crictl
RUN wget https://github.com/kubernetes-sigs/cri-tools/releases/download/${CRICTL_VERSION}/crictl-${CRICTL_VERSION}-linux-amd64.tar.gz && \
    tar zxvf crictl-${CRICTL_VERSION}-linux-amd64.tar.gz -C /usr/local/bin && \
    rm -f crictl-${CRICTL_VERSION}-linux-amd64.tar.gz

# Specify the default image endpoint for crictl
RUN echo 'runtime-endpoint: unix:///run/containerd/containerd.sock' >> /etc/crictl.yaml
RUN echo 'image-endpoint: unix:///run/containerd/containerd.sock' >> /etc/crictl.yaml
RUN echo 'timeout: 2' >> /etc/crictl.yaml

# Install amd-smi (ROCm standalone; for AMD GPU node debugging)
ARG ROCM_MAJOR_MINOR="7.14"
ARG ROCM_PKG_VERSION="7.14.0-3"

RUN mkdir --parents --mode=0755 /etc/apt/keyrings && \
    wget https://repo.amd.com/rocm/packages-multi-arch/gpg/rocm.gpg -O - | \
      gpg --dearmor | tee /etc/apt/keyrings/amdrocm.gpg > /dev/null && \
    echo "deb [arch=amd64 signed-by=/etc/apt/keyrings/amdrocm.gpg] https://repo.amd.com/rocm/packages-multi-arch/debian13 stable main" \
      > /etc/apt/sources.list.d/rocm.list && \
    apt-get update -qq && \
    apt-get install -y "amdrocm-amdsmi${ROCM_MAJOR_MINOR}=${ROCM_PKG_VERSION}" && \
    echo "/opt/rocm/core-${ROCM_MAJOR_MINOR}/lib" > /etc/ld.so.conf.d/amd-rocm-smi.conf && \
    ldconfig && \
    ln -sf "/opt/rocm/core-${ROCM_MAJOR_MINOR}/bin/amd-smi" /usr/local/bin/amd-smi && \
    rm -f /etc/apt/sources.list.d/rocm.list && \
    apt-get clean && rm -rf /var/lib/apt/lists/* && \
    rm /etc/apt/keyrings/amdrocm.gpg


CMD [ "/bin/bash" ]
