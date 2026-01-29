FROM alpine:3.12

RUN apk add --no-cache bash

ENV DB_PASSWORD="supersecret123" \
    API_KEY="123456789"

USER root

EXPOSE 8080

CMD ["sh", "-c", "echo 'Этот контейнер УМЫШЛЕННО небезопасен ^_^'"]
