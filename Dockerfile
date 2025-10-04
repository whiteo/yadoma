FROM alpine:3.22

ARG DOCKER_GID=999
RUN addgroup -g $DOCKER_GID docker && adduser -D -H -s /sbin/nologin -G docker app

COPY build/yadoma-agent /usr/local/bin/yadoma-agent

USER app

EXPOSE 50001

ENTRYPOINT ["yadoma-agent"]