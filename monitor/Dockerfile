FROM ubuntu:18.04
WORKDIR /app
RUN apt-get update && apt-get install -y wget unzip tar

RUN wget https://github.com/grafana/loki/releases/download/v2.3.0/loki-linux-amd64.zip && unzip loki-linux-amd64.zip
RUN wget https://dl.grafana.com/oss/release/grafana-8.1.5.linux-amd64.tar.gz && tar -zxvf grafana-8.1.5.linux-amd64.tar.gz

COPY loki-local-config.yaml /app/
COPY docker-entrypoint.sh /app/
COPY grafana-datasources.yml /app/grafana-8.1.5/conf/provisioning/datasources/datasources.yml
COPY grafana-dashboards.yml /app/grafana-8.1.5/conf/provisioning/dashboards/dashboards.yml
COPY ipfs-dashboard.json /var/lib/grafana/dashboards/

CMD [ "./docker-entrypoint.sh" ]
EXPOSE 3000 3100
