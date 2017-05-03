FROM nginx

RUN apt-get update && apt-get install -y supervisor

RUN apt-get clean && rm -rf /var/cache/apt/archives/*

ADD supervisord.conf /etc/supervisor/conf.d/supervisord.conf

ADD entrypoint.sh /entrypoint.sh

ADD sha1check.sh /usr/local/bin/nginx-reload

RUN mkdir /opt/nginx-k8sapi/ -p
ADD k8s-api /usr/local/bin/k8s-api
ADD conf.example /opt/nginx-k8sapi/conf.example

RUN chmod +x /usr/local/bin/nginx-reload
RUN chmod +x /usr/local/bin/k8s-api
RUN chmod +x /entrypoint.sh

ENTRYPOINT ["/usr/bin/supervisord"]
