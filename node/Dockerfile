FROM ubuntu:18.04
WORKDIR /app
RUN apt-get update && apt-get install -y wget unzip tar git make build-essential

RUN wget https://github.com/grafana/loki/releases/download/v2.3.0/promtail-linux-amd64.zip && unzip ./promtail-linux-amd64.zip
RUN wget https://golang.org/dl/go1.17.1.linux-amd64.tar.gz && tar -C /usr/local -xzf go1.17.1.linux-amd64.tar.gz
ENV PATH="/usr/local/go/bin:${PATH}"
COPY promtail-local-config.yaml /app/
RUN git clone https://github.com/wcgcyx/go-libp2p-kad-dht.git && cd go-libp2p-kad-dht && git checkout more-logging && cd ..
RUN git clone https://github.com/wcgcyx/go-ipfs.git  && cd go-ipfs && git checkout more-logging && echo "replace github.com/libp2p/go-libp2p-kad-dht => ../go-libp2p-kad-dht" >> go.mod && make build
RUN mkdir /ipfs-tests/
COPY docker-entrypoint.sh /app/
COPY agent /app/

CMD [ "./docker-entrypoint.sh" ]
EXPOSE 4001
EXPOSE 4001/udp
EXPOSE 5001
EXPOSE 8080
