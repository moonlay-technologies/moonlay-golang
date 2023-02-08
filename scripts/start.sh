#bin/bash
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com
docker pull 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/order-service:"$DEPLOYMENT_GROUP_NAME"-latest
docker run --restart always --log-opt max-size=10m --log-opt max-file=5 -d -p 8000:8000 --name order-service 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/order-service:"$DEPLOYMENT_GROUP_NAME"-latest