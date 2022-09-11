# module9  作业

## 提交作业

### 导入镜像

```shell
# docker 
docker load < httpserver-v0.5.tar
# containerd
ctr -n k8s.io i import httpserver-v0.5.tar
```

### 创建应用

[httpserver](httpserver.yaml)

```shell
kubectl create -f httpserver.yaml
```

