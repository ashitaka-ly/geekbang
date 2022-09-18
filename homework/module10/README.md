# module10  作业

## 提交作业

### 导入镜像

```shell
# docker 
docker load < httpserver-metrics-v0.6.tar
# containerd
ctr -n k8s.io i import httpserver-metrics-v0.6.tar
```

### 创建应用

[httpserver-metrics](httpserver-v0.6.yaml)

```shell
kubectl create -f httpserver-v0.6.yaml
```