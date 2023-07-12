# 云服务模式按抽象层级由低到高有哪些？

### IaaS

Examples: Amazon Web Services (AWS), Microsoft Azure, Google Cloud Platform (GCP), Alibaba Cloud, DigitalOcean…

This is the least abstract layer and the closest to managing physical servers. IaaS provides the infrastructure such as servers, virtual machines, networks, operating system, storage, and other low-level resources. The user is still responsible for managing the server, including the applications, runtime, OS, middleware, and data.

### CaaS or **Managed Kubernetes Services**

Examples: Google Kubernetes Engine, Azure Kubernetes Service, Amazon EKS…

With CaaS (Container as a Service) platforms like GKE, you have more control over the underlying infrastructure. You can define how your applications should run, how they should communicate, and how they should scale. However, you don't need to manage the Kubernetes control plane, as GKE takes care of that. In other words, GKE provides a level of abstraction somewhere between Infrastructure as a Service (IaaS) and PaaS.

### PaaS

Examples: Cloud Foundry, Heroku, Google App Engine…

PaaS offers a higher level of abstraction by providing a platform that includes infrastructure, runtime environment, and development tools. Users focus only on developing their applications and deploying their code, while the service provider manages the underlying infrastructure, runtime, middleware, and operating system.

### FaaS or Serverless **Platforms**

Examples: AWS Lambda, Google Cloud Functions, Azure Functions, etc.

FaaS (Function as a Service) is the most abstract level, where developers only focus on writing individual functions or pieces of business logic, and the infrastructure is fully managed by the service provider. FaaS falls under the umbrella of **serverless** computing, because the server management is entirely abstracted away from the developer.