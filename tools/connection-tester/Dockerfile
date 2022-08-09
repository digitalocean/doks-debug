FROM digitalocean/doks-debug:latest

ADD ./connection-tester ./
COPY ./connection-tester /usr/bin/connection-tester
RUN chmod a+x /usr/bin/connection-tester
