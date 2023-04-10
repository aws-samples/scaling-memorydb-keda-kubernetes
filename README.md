# Auto-scale Redis based data processing applications using Amazon EKS, Amazon MemoryDB and KEDA

This repo contains sample code and configuration to demonstrate auto scale Redis-based data processing applications. We will be using:
- [Amazon MemoryDB for Redis](https://docs.aws.amazon.com/memorydb/latest/devguide/what-is-memorydb-for-redis.html),
- Kubernetes cluster on [Amazon Elastic Kubernetes Service](https://docs.aws.amazon.com/eks/latest/userguide/getting-started.html) (Amazon EKS) - That's where the data processing application that needs to be auto-scaled will be deployed to, and, 
- [KEDA](keda.sh) - An open source, Kubernetes-based, event-driven auto scaler.

For details, [check out the blog post](TODO).

## Security

See [CONTRIBUTING](CONTRIBUTING.md#security-issue-notifications) for more information.

## License

This library is licensed under the MIT-0 License. See the LICENSE file.