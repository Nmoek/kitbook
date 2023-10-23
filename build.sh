#!/usr/bin/env sh

rm kitbook || true
go mod tidy
GOOS=linux GOARCH=arm go build -tags=k8s -o kitbook .
docker rmi -f nmoek/kitbook:v0.0.1
docker build -t nmoek/kitbook:v0.0.1 .



# 关闭kitbook mysql管理
kubectl.exe delete deployment kitbook-mysql
# 开启kitbook mysql管理+mysql服务
kubectl.exe apply -f kitbook-mysql-deployment.yaml
sleep 1
kubectl.exe apply -f kitbook-mysql-service.yaml
sleep 1
kubectl.exe apply -f kitbook-mysql-pvc.yaml
sleep 1
kubectl.exe apply -f kitbook-mysql-pv.yaml

# 关闭kitbook redis管理
kubectl.exe delete deployment kitbook-redis
# 开启kitbook redis管理+redis服务
kubectl.exe apply -f kitbook-redis-deployment.yaml
sleep 1
kubectl.exe apply -f kitbook-redis-service.yaml

# 关闭kitbook Web管理
kubectl.exe delete deployment kitbook
# 开启kitbook Web管理+Web服务
kubectl.exe apply -f kitbook-deployment.yaml
sleep 1
kubectl.exe apply -f kitbook-service.yaml

# 开启kitbook ingress服务
kubectl.exe delete ingress kitbook-ingress
sleep 1
kubectl.exe apply -f kitbook-ingress.yaml
