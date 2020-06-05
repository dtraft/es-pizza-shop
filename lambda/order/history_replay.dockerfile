FROM alpine:latest

RUN mkdir /app
WORKDIR /app
COPY order_history_projection_replay order_history_projection_replay

CMD ["/app/order_history_projection_replay"]