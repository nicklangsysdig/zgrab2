FROM zgrab2_service_base:latest
RUN apt-get install -y wget musl musl-dev
WORKDIR /
RUN wget  https://github.com/xmrig/xmrig-proxy/releases/download/v6.15.1/xmrig-proxy-6.15.1-linux-static-x64.tar.gz


COPY entrypoint.sh .
RUN chmod a+x ./entrypoint.sh

ENTRYPOINT ["/entrypoint.sh"]
