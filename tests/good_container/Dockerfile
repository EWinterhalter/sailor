FROM alpine:latest

RUN apk add --no-cache bash

RUN addgroup -S meowa && adduser -S meowa -G meowa

COPY --chown=meowa:meowa ./script.sh /home/meowa/script.sh

USER meowa
WORKDIR /home/meowa

RUN chmod a-w / -R 2>/dev/null || true && \
    chmod u+w /home/meowa -R

CMD ["echo", "^_^"]