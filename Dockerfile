FROM 192.168.1.254:5000/alpine:3.9.5

COPY resources/ /opt/registryui/resources/
COPY bin/ /opt/registryui/

WORKDIR /opt/registryui

EXPOSE 8080

CMD ["/opt/registryui/registryui", "-debug", "false"]

