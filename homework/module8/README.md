# module8 作业

## 提交作业

### 导入镜像

```shell
# docker 
docker load < httpserver-v0.4.tar
# containerd
ctr -n k8s.io i import httpserver-v0.4.tar
```

### 创建 configMap

[configMap](config.yaml)

```shell
kubectl create -f config.yaml
```

### 创建应用

[httpserver](httpserver.yaml)

```shell
kubectl create -f httpserver.yaml
```

待完善
- [X] 使用configmap控制日志等级 —— 已实现
- [ ] 优化日志文件的管理

## 作业内容

### 代码接收sig
需要修改 go 代码，以实现优雅停止的 sigterm

[go代码 main.go](../../httpserver/main.go)

### 日志控制

~~已经把日志级别控制移动到 flag，需要考虑容器控制~~

构建时使用环境变量控制  
[dockerfile](../../Dockerfile)

### 生成镜像

```shell
docker build . -t httpserver:0.4
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

### 组织 deployment

```yaml
apiVersion: apps/v1
kind: Deployment
metadata:
  labels:
    app: httpserver
  name: httpserver
spec:
  replicas: 3
  template:
    spec:
      containers:
        - name: httpserver
          image: httpserver:0.4
          readinessProbe:
              failureThreshold: 3
              httpGet:
                path: /healthz
                port: 8080
                scheme: HTTP
              initialDelaySeconds: 5
              periodSeconds: 10
              successThreshold: 1
              timeoutSeconds: 1
          livenessProbe:
            failureThreshold: 3
            httpGet:
              path: /healthz
              port: 8080
              scheme: HTTP
            initialDelaySeconds: 5
            periodSeconds: 10
            successThreshold: 1
            timeoutSeconds: 1
          resources:
            limits:
              cpu: 200m
              memory: 200Mi
            requests:
              cpu: 100m
              memory: 100Mi
      restartPolicy: Always
```

相对简陋，参考助教老师修改
