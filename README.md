# elder-grater
EKS IRSA demonstration to access S3 and SQS

1. Build the application and deploy it to a container registry (i.e. AWS ECR, DockerHub, etc)
1. Create an IAM role "elder-grater" 
    1. From AWS console, select IAM and "Create role"
    1. Check "Web identity" 
        1. Identity provider is your EKS cluster
            1. Has the form "oidc.eks.region.amazonaws.com/id/many-hex-characters"
        1. Audience
            1. "sts.amazonaws.com"
    1. Select "Next"
    1. Add permissions
        1. AmazonS3FullAccess
        1. AmazonSQSFullAccess
    1. Select "Next"
    1. Add role name "elder-grater"
    1. Select "Create role"
    1. [example](https://github.com/guycole/elder-grater/blob/main/iam_role.json)
1. Deploy the k8s ServiceAccount
    1. Edit [service_account.yaml](https://github.com/guycole/elder-grater/blob/main/service_account.yaml)
        1. Change "replace-me" with your AWS account number
    1. Deploy the SA ("kubectl apply -f service_account.yaml")
1. Deploy the pod
    1. Edit [deployment.yaml](https://github.com/guycole/elder-grater/blob/main/deployment.yaml)
        1. Update image spec to reflect true location
    1. Deploy the pod ("kubectl apply -f deployment.yaml")
1. Review the log
    1. "kubectl logs elder-grater -f"
    1. Success means you see happy AWS login and SQS/S3 information
1. Cleanup
    1. kubectl delete -f pod.yaml
    1. kubectl delete -f service_account.yaml
    1. Delete the IAM role "elder grater" via AWS console
