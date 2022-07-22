FROM digitalocean/doks-debug:latest

ADD ./test-server ./
COPY ./test-server/ /usr/bin/test-server
RUN chmod a+x /usr/bin/test-server
