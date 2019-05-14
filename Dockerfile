FROM ubuntu
WORKDIR /root

RUN sed -i '/path-exclude=\/usr\/share\/man\/*/c\#path-exclude=\/usr\/share\/man\/*' /etc/dpkg/dpkg.cfg.d/excludes

RUN apt-get update -qq && \
    apt-get install -y man \
                       manpages-posix \
                       man-db \
                       vim \
                       screen \
                       curl \
                       jq \
                       docker.io \
                       dnsutils \
                       tcpdump \
                       traceroute \
                       iputils-ping \
                       net-tools \
                       netcat \
                       iproute2 \
                       strace

CMD [ "/bin/bash" ]