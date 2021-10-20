# Specification

## Components

### Project
A project represents a Kubernetes namespace. For example, project `prj_A` is a Kubernetes namespace named `prj_A`.

### App
An app represents a Kubernetes deployment and service. For example, an app `app_A` in project `prj_A` is a Kubernetes deployment named `app_A` in namespace `prj_A`.

### Deployment
A deployment represents a container image. For example, a deployment `dpl_A` of app `app_A` is a container image named `[CONTAINER_REGISTRY]/app_A:dpl_A`.
