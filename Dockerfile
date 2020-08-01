FROM alpine:3.9.5

COPY resources/ /opt/registryui/resources/
COPY bin/ /opt/registryui/

WORKDIR /opt/registryui

EXPOSE 8080

ENTRYPOINT ["/opt/registryui/registryui"]

CMD ["-debug", "false"]
