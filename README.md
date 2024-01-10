# Custom Kubernetes Pod Scheduler For Argo Project

This is sample code for deploying a controller for custom pod scheduling in Kubernetes for Argo project

***This project also supports Argo API for enhanced capabilities in managing pod scheduling.***



get inspired from [custom-kubernetes-scheduler](https://github.com/aws-samples/containers-blog-maelstrom/tree/main/custom-kubernetes-scheduler)

---
## [Customizing scheduling on Amazon EKS](https://aws.amazon.com/ko/blogs/containers/customizing-scheduling-on-amazon-eks/)
#### Introduction

In Kubernetes, pod scheduling is handled by the kube-scheduler process, which places pods on nodes based on resource requests, affinity rules, topology spread, and other considerations. The default behavior is to spread pods across nodes, but there are cases where more fine-grained control is needed.

#### Solution

This post presents a way to define a custom pod scheduling strategy using a mutating admission webhook. This strategy can use node labels to schedule pods proportionally.

#### Example

As an example, we looked at how to deploy pods across On-Demand nodes and Spot instances.

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  name: nginx
  namespace: test
  annotations:
    custom-pod-schedule-strategy: 'karpenter.sh/capacity-type=on-demand,base=2,weight=1:karpenter.sh/capacity-type=spot,weight=3'
spec:
  selector:
    matchLabels:
      app: nginx
  replicas: 10
  template:
    metadata:
      labels:
        app: nginx
    spec:
      containers:
      - image: public.ecr.aws/nginx/nginx:latest
        imagePullPolicy: Always
        name: nginx
        resources:
          limits:
            cpu:  400m
            memory: 1600Mi
          requests:
            cpu: 400m
            memory: 1600Mi        
        ports:
        - name: http
          containerPort: 80
```

this example, the **karpenter.sh/capacity-type** label is used to identify the node type. The **base** parameter specifies the minimum number of pods to be scheduled on each node type. The weight parameter specifies the relative distribution of pods across node types.

In this case, the base parameter for the karpenter.sh/capacity-type=on-demand label is set to 2. This means that at least 2 pods must be scheduled on On-Demand nodes. The **weight** parameter for the **karpenter.sh/capacity-type=on-demand** label is set to 1. This means that 1/4 of all pods should be scheduled on On-Demand nodes.

The weight parameter for the karpenter.sh/capacity-type=spot label is set to 3. This means that 3/4 of all pods should be scheduled on Spot instances.

#### Conclusion

This post demonstrated how you can use a mutating pod admission webhook to customize pod scheduling across nodes. You can use this solution for a variety of use cases such as prioritizing nodes in an Availability Zone to reduce data transfer costs, spreading workloads across Availability Zones, or running workloads across On-Demand and Spot instances.

#### Key Words
- Kubernetes
- Pod scheduling
- Admission webhook
- Node labels
- Proportional placement
- On-Demand nodes
- Spot instances
- Argo Projects


#### Readme

This repository contains the source code and deployment files for a proof of concept that demonstrates how to use a mutating admission webhook to customize pod scheduling across nodes in Kubernetes.

To use this solution, you will need:

- An Amazon Elastic Kubernetes Service (EKS) cluster
- The Kubernetes CLI (kubectl)
- Golang
- Docker
- Argo-cd
- Argo-rollouts
- cfssl
- cfssljson

To deploy the solution, follow these steps:

1. Clone the repository:
```bash
git clone https://github.com/paikend/custom-kubernetes-scheduler-for-argoproj.git
```
2. Change directory to the custom-kubernetes-scheduler directory:
```bash
cd custom-kubernetes-scheduler-for-argopro
```
3. Create an Amazon Elastic Container Registry (ECR) repository to store the container image for the admission webhook:
```bash
IMAGE_REPO="ACCOUNTID.dkr.ecr.${AWS_REGION}.amazonaws.com"
IMAGE_NAME="ECR_REPO"
export ECR_REPO_URI=$(aws ecr describe-repositories --repository-name ${IMAGE_NAME} | jq -r '.repositories[0].repositoryUri')

if [ -z "$ECR_REPO_URI" ]; then
  echo "${IMAGE_REPO}/${IMAGE_NAME} does not exist. So creating it..."
  ECR_REPO_URI=$(aws ecr create-repository \
    --repository-name $IMAGE_NAME \
    --region $AWS_REGION \
    --query 'repository.repositoryUri' \
    --output text)
  echo "ECR_REPO_URI=$ECR_REPO_URI"
else
  echo "${IMAGE_REPO}/${IMAGE_NAME} already exists..."
fi
```
