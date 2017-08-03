FROM rabbitmq:3.6.10-management

COPY rabbitmq-cluster /usr/local/bin/
COPY standlone.conf /etc/rabbitmq/rabbitmq-standalone.config
COPY cluster.conf /etc/rabbitmq/rabbitmq-cluster.config
COPY pre-entrypoint.sh /docker-entrypoint.sh

RUN chmod u+x /docker-entrypoint.sh && chmod +x /usr/local/bin/rabbitmq-cluster

EXPOSE 5672 15672 25672 4369 9100 9101 9102 9103 9104 9105
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["rabbitmq-cluster"]
