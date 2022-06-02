FROM ubuntu:focal
WORKDIR /root

RUN sed -i '/path-exclude=\/usr\/share\/man\/*/c\#path-exclude=\/usr\/share\/man\/*' /etc/dpkg/dpkg.cfg.d/excludes

RUN apt-get update -qq && \
    apt-get install -y apt-transport-https \
                       ca-certificates \
                       software-properties-common \
                       man \
                       manpages-posix \
                       man-db \
                       vim \
                       screen \
                       curl \
                       atop \
                       htop \
                       dstat \
                       jq \
                       dnsutils \
                       tcpdump \
                       traceroute \
                       iputils-ping \
                       net-tools \
                       netcat \
                       iproute2 \
                       strace \
                       telnet \
                       openssl \
                       psmisc \
                       dsniff \
                       conntrack

RUN curl -fsSL https://download.docker.com/linux/ubuntu/gpg | apt-key add - && \
    add-apt-repository "deb [arch=amd64] https://download.docker.com/linux/ubuntu $(lsb_release -cs) stable" && \
    apt-get update -qq && \
    apt-get install -y docker-ce

CMD [ "/bin/bash" ]
