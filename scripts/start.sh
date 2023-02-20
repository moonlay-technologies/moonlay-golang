#bin/bash
aws ecr get-login-password --region ap-southeast-1 | docker login --username AWS --password-stdin 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com
docker pull 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/order-service:"$DEPLOYMENT_GROUP_NAME"-latest
docker run --restart always --log-opt awslogs-stream=order-service-api -d -p 8000:8000 --name order-service 650142038379.dkr.ecr.ap-southeast-1.amazonaws.com/order-service:"$DEPLOYMENT_GROUP_NAME"-latest
cloudfront_id=""

if [ "$DEPLOYMENT_GROUP_NAME" == "staging" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id StgCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$DEPLOYMENT_GROUP_NAME" == "development" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id DevCloudFrontID --query SecretString --output text --region ap-southeast-1)
elif [ "$DEPLOYMENT_GROUP_NAME" == "production" ]; then
  cloudfront_id=$(aws secretsmanager get-secret-value --secret-id ProdCloudFrontID --query SecretString --output text --region ap-southeast-1)
fi

aws cloudfront create-invalidation --distribution-id "$cloudfront_id" --paths "/*"