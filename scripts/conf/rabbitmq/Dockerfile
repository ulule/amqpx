FROM rabbitmq:3.6.10-management

COPY rabbitmq-ping /usr/local/bin/
COPY rabbitmq-cluster /usr/local/bin/
COPY standlone.conf /etc/rabbitmq/rabbitmq-standalone.config
COPY cluster.conf /etc/rabbitmq/rabbitmq-cluster.config
COPY docker-entrypoint.sh /docker-entrypoint.sh

ENV RABBITMQ_LOGS=/var/log/rabbitmq/rabbitmq.log RABBITMQ_SASL_LOGS=/var/log/rabbitmq/rabbitmq-sasl.log

RUN apt-get -y update \
    && apt-get upgrade -y \
    && apt-get install -y amqp-tools \
    && apt-get clean \
    && rm -rf /var/lib/apt/lists/* \
    && chmod u+x /docker-entrypoint.sh \
    && chmod +x /usr/local/bin/rabbitmq-cluster \
    && chmod +x /usr/local/bin/rabbitmq-ping

EXPOSE 5672 15672 25672 4369 9100 9101 9102 9103 9104 9105
ENTRYPOINT ["docker-entrypoint.sh"]
CMD ["rabbitmq-cluster"]
