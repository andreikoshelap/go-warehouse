FROM alpine:latest

RUN mkdir /app

COPY productServiceApp /app

CMD [ "/app/productServiceApp"]