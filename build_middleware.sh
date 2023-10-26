#!/usr/bin/env sh

#rm kitbook || true
#go mod tidy
#GOOS=linux GOARCH=arm go build -tags=k8s -o kitbook .
##GOOS=linux GOARCH=arm go build -o kitbook .
#docker rmi -f nmoek/kitbook:v0.0.1
#docker build -t nmoek/kitbook:v0.0.1 .


# 关闭kitbook mysql管理+mysql服务
kubectl.exe delete deployment kitbook-mysql
kubectl.exe delete service kitbook-mysql


# 关闭kitbook redis管理+redis服务
kubectl.exe delete deployment kitbook-redis
kubectl.exe delete service kitbook-redis


# 关闭kitbook Web管理+Web服务
#kubectl.exe delete deployment kitbook
#kubectl.exe delete service kitbook

# 关闭ingress 路由转发规则
kubectl.exe delete ingress kitbook-ingress


# 开启kitbook mysql管理+mysql服务
kubectl.exe apply -f kitbook-mysql-deployment.yaml
kubectl.exe apply -f kitbook-mysql-service.yaml
kubectl.exe apply -f kitbook-mysql-pvc.yaml
kubectl.exe apply -f kitbook-mysql-pv.yaml

# 开启kitbook redis管理+redis服务
kubectl.exe apply -f kitbook-redis-deployment.yaml
kubectl.exe apply -f kitbook-redis-service.yaml

# 开启kitbook Web管理+Web服务
#kubectl.exe apply -f kitbook-deployment.yaml
#kubectl.exe apply -f kitbook-service.yaml

# 开启kitbook ingress服务
kubectl.exe apply -f kitbook-ingress.yaml
