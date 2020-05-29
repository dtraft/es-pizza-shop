FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY order_projection_replay order_projection_replay

CMD ["/app/order_projection_replay"]