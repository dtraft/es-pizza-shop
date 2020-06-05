make package_order_history_projection_replay

docker tag order-history-projection-replay 849624642218.dkr.ecr.us-east-1.amazonaws.com/replay/order-history-projection:latest

docker push 849624642218.dkr.ecr.us-east-1.amazonaws.com/replay/order-history-projection:latest
