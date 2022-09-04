# module8 作业

## 提交作业

### 导入镜像

```shell
docker load < httpserver-0.3.tar
```

### 创建 configMap

```shell
kubectl create -f config.yaml
```

### 创建应用

```shell
kubectl create -f httpserver.yaml
```

## 作业内容

### 代码接收sig
需要修改 go 代码，以实现优雅停止的 sigterm

[go代码 main.go](../../httpserver/main.go)

### 生成镜像

```shell
docker build . -t httpserver:0.3
```

### 组织 configMap

```yaml
apiVersion: v1
data:
  loglevel: "4"
kind: ConfigMap
metadata:
  name: myenv
```
