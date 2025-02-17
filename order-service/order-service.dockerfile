FROM alpine:latest

RUN mkdir /app

COPY orderServiceApp /app

CMD [ "/app/orderServiceApp"]