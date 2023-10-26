#!/usr/bin/env sh

rm kitbook || true
go mod tidy
GOOS=linux GOARCH=arm go build -tags=k8s -o kitbook .
#GOOS=linux GOARCH=arm go build -o kitbook .
docker rmi -f nmoek/kitbook:v0.0.1
docker build -t nmoek/kitbook:v0.0.1 .

# 关闭kitbook Web管理+Web服务
kubectl.exe delete deployment kitbook
kubectl.exe delete service kitbook
sleep 1
# 开启kitbook Web管理+Web服务
kubectl.exe apply -f kitbook-deployment.yaml
kubectl.exe apply -f kitbook-service.yaml