FROM alpine

RUN mkdir -p /usr/app/accurics-service/.
WORKDIR /usr/app/accurics-service/.

COPY accurics-service .
COPY config.json ./config.json

RUN apk --no-cache add ca-certificates

CMD ["./accurics-service"]