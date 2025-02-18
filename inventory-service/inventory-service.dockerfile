FROM alpine:latest

RUN mkdir /app

COPY inventoryServiceApp /app

CMD [ "/app/inventoryServiceApp"]