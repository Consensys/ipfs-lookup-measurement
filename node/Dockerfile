FROM ubuntu:18.04
WORKDIR /app
RUN apt-get update && apt-get install -y wget unzip tar

RUN wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip && unzip ./promtail-linux-amd64.zip
RUN wget https://dist.ipfs.io/go-ipfs/v0.10.0/go-ipfs_v0.10.0_linux-amd64.tar.gz && tar -zxvf ./go-ipfs_v0.10.0_linux-amd64.tar.gz

COPY promtail-local-config.yaml /app/
COPY docker-entrypoint.sh /app/

CMD [ "./docker-entrypoint.sh" ]
EXPOSE 4001
EXPOSE 4001/udp
EXPOSE 5001
EXPOSE 8080
