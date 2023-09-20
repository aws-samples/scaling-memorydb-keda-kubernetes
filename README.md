# Auto-scale Redis based data processing applications using Amazon EKS, Amazon MemoryDB and KEDA

This repo contains sample code and configuration to demonstrate auto scale Redis-based data processing applications. We will be using:
- [Amazon MemoryDB for Redis](https://docs.aws.amazon.com/memorydb/latest/devguide/what-is-memorydb-for-redis.html),
- Kubernetes cluster on [Amazon Elastic Kubernetes Service](https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html) (Amazon EKS) - That's where the data processing application that needs to be auto-scaled will be deployed to, and, 
- [KEDA](keda.sh) - An open source, Kubernetes-based, event-driven auto scaler.

For details, [check out the blog post].

## Installation and Configuration

### Install KEDA in your K8s cluster

See [scripts/keda.sh](scripts/keda.sh)
```bash
kubectl apply -f https://github.com/kedacore/keda/releases/download/v2.10.0/keda- 2.10.0.yaml
kubectl get crd
kubectl get deployment -n keda
kubectl logs -f $(kubectl get pod -l=app=keda-operator -o jsonpath='{.items[0].metadata.name}' -n keda) -n keda
```

### Create environment variables

Create `.env.sh` use [scripts/.env.sh.example](scripts/.env.sh.example)
```bash
export AWS_ACCOUNT_ID=<enter AWS account ID>
export AWS_REGION=<enter AWS region e.g. us-east-1> aws ecr get-login-password --region ${AWS_REGION} | docker login --username AWS -- password-stdin ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com
export CONSUMER_APP=redis-consumer-app
export PRODUCER_APP=redis-producer-app aws ecr create-repository --repository-name ${CONSUMER_APP} --region ${AWS_REGION} aws ecr create-repository --repository-name ${PRODUCER_APP} --region ${AWS_REGION}
```

### Build Proucer and Consumer apps and push them to ECR
```bash
# if you're on Mac M1, use: export DOCKER_DEFAULT_PLATFORM=linux/amd64
# for producer 
docker build -t ${PRODUCER_APP} producer 
docker tag ${PRODUCER_APP}:latest ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${PRODUCER_APP}:latest
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${PRODUCER_APP}:latest
# for consumer
docker build -t ${CONSUMER_APP} consumer
docker tag ${CONSUMER_APP}:latest ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${CONSUMER_APP}:latest
docker push ${AWS_ACCOUNT_ID}.dkr.ecr.${AWS_REGION}.amazonaws.com/${CONSUMER_APP}:latest
```

### Deploy the application to Amazon EKS

1. Deploy the consumer application with the following code:
```bash
kubectl apply -f deploy/consumer.yaml
# check logs
kubectl logs -f $(kubectl get pod -l=app=redis-consumer -o jsonpath='{.items[0].metadata.name}')
```

2. Before you create the `ScaledObject`, start a watch on the consumer application deployment:
```bash
kubectl get deployment/redis-consumer -w
```

3. In a separate terminal, deploy the `ScaledObject`:
```bash
kubectl apply -f deploy/scaled-object.yaml
```

### Auto scaling in action

1. In a different terminal, deploy the producer app, which starts generating data and sending it to the Redis List:
```bash
kubectl apply -f deploy/producer.yaml
```

2. Go back to the terminal where you’re tracking the consumer deployment:
```bash
kubectl get pod -l=app=redis-consumer
```

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.