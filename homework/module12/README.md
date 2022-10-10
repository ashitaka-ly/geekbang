
# module 12 作业


1. 创建应用

```shell
kubectl create ns homework
kubectl label ns homework istio-injection=enabled
kubectl create -f app.yaml -n homework
```

2. 配置 gateway

```shell
kubectl create -n istio-system secret tls httpserver-cert --key=cncamp.io.key --cert=cncamp.io.crt
kubectl create -f gw.yaml -n homework
```

3. 配置 VirtualService

```shell
kubectl -f create vs.yaml -n homework
```
