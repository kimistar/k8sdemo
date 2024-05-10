## kubectl常用命令

kubectl get deploy
kubectl get svc
kubectl get po
kubectl get no
kubectl get ns

apply: 用于在集群中创建或更新资源。与 create 不同的是，apply 可以部分更新资源，而不是覆盖整个资源配置
kubectl apply -f deployment.yml

describe：用于显示资源的详细信息
kubectl describe pod my-pod

delete：用于删除集群中的资源
kubectl delete pod my-pod

exec：用于在运行中的容器中执行命令
kubectl exec -it my-pod -- /bin/bash

scale：用于扩容或缩容部署的副本数量
kubectl scale deploy <deploy的名称> --replicas=6

自动扩容
kubectl autoscale deploy <deploy的名称> --cpu-percent=20 --min=2 --max=5

镜像更新后，重新部署
kubectl set image deployment/k8s-demo-deployment<deployment名称> k8s-demo<容器名称>=alleninnz/k8s-demo

edit: 进行deployment.yml文件中修改
kubectl edit deploy <deploy的名称>

暂停与恢复(修改template之后不会立刻生效，恢复后才生效)
暂停：kubectl rollout pause deploy <deploy的名称>
恢复：kubectl rollout deploy <deploy的名称>

## 滚动更新：

更新deployment.yml文件中的template之后会自动触发滚动更新

查看滚动更新的状态：kubectl rollout status deploy <deploy的名称>
查看全部版本列表：kubectl rollout history deploy <deploy的名称>
查看指定版本的信息：kubectl rollout history deploy <deploy的名称> --revision=2 <版本号>
回退到上个版本：kubectl rollout undo deploy <deploy的名称>
回退到指定版本：kubectl rollout undo deploy <deploy的名称> --to-revision=2 <版本号>
    
## HPA

Kubernetes 提供了自动水平扩展（Horizontal Pod Autoscaling，HPA）功能，可以根据请求流量或其他指标来自动扩展或缩减 Pod 的副本数量。HPA 是 Kubernetes 的核心功能之一，可以根据你的配置动态地调整 Pod 的副本数量，以满足应用程序的需求，从而实现更好的性能和资源利用率。

要实现自动水平扩展，你需要完成以下几个步骤：

部署指标服务器：
首先，你需要在你的 Kubernetes 集群中部署一个指标服务器，比如 Heapster、Metrics Server 或者 Prometheus Operator。指标服务器负责收集集群中各种资源的指标数据，比如 CPU 使用率、内存使用率等。

创建 HorizontalPodAutoscaler 对象：
接下来，你需要创建一个 HorizontalPodAutoscaler（HPA）对象，并指定你希望自动扩展的目标资源、触发条件、最小和最大副本数量等配置。HPA 将监视指标服务器收集到的指标数据，并根据配置来动态地调整 Pod 的副本数量。

设置触发条件：
在 HPA 对象中，你可以设置触发条件，比如 CPU 使用率超过某个阈值或者每秒请求的数量超过某个阈值。当触发条件满足时，HPA 将增加 Pod 的副本数量；当触发条件不再满足时，HPA 将减少 Pod 的副本数量。

监视和调整：
一旦创建了 HPA 对象，Kubernetes 将会自动监视和调整 Pod 的副本数量，以满足你配置的触发条件。HPA 将定期检查指标服务器收集到的指标数据，并根据触发条件来调整 Pod 的副本数量。

列出当前的 HPA 对象：kubectl get hpa
删除HPA对象：kubectl delete hpa <hpa名称>

## 使用addons

minikube addons enable metrics-server
minikube addons enable ingress-nginx

### 手动安装

wget https://github.com/kubernetes-sigs/metrics-server/releases/latest/download/components.yaml

vim components.yaml 搜索containers /containers 添加 - --kubelet-insecure-tls （可选 http）

kubectl apply -f components.yaml

验证：kubectl get deployment metrics-server -n kube-system

### 查看指标

kubectl top pod

## ~~helm安装Ingress~~
helm repo add ingress-nginx https://kubernetes.github.io/ingress-nginx
helm repo update
helm pull ingress-nginx/ingress-nginx
tar -zxvf ingress-nginx-3.34.0.tgz
cd ingress-nginx

创建namespace ingress-nginx
kubectl create ns ingress-nginx

修改values.yaml文件
hostNetwork: true
dnsPolicy: ClusterFirstWithHostNet
kind: DaemonSet
nodeSelector:
   ingress: "true"
service:
   type: ClusterIP
admissionWebhooks:
   enabled: false

`节点上打上标签`： kubectl label node minikube ingress=true

helm install ingress-nginx -n ingress-nginx .


## 运行NFS

docker run -d --name nfs-server \
    --privileged \
    --restart=always \
    -v /data/nfs:/nfsshare \
    -e SHARED_DIRECTORY=/nfsshare \
    -e SYNC=true \
    -p 2049:2049 \
    itsthenetwork/nfs-server-alpine:latest

docker exec -it <containerName> /bin/bash

### K8S NFS示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: nfs-server
spec:
  containers:
  - name: nfs-server
    image: itsthenetwork/nfs-server-alpine:latest
    securityContext:
      privileged: true
    env:
    - name: SHARED_DIRECTORY
      value: "/nfsshare"
    - name: SYNC
      value: "true"
    volumeMounts:
    - name: nfs-share
      mountPath: /nfsshare
  volumes:
  - name: nfs-share
    hostPath:
      path: /path/to/nfs/share
---
apiVersion: v1
kind: Service
metadata:
  name: nfs-service
spec:
  selector:
    pod: nfs-server
  ports:
    - protocol: TCP
      port: 2049
      targetPort: 2049

```

## ubuntu install NFS

apt install nfs-kernel-server

vim /etc/exports

/path/to/share *(rw,sync,no_subtree_check,no_root_squash,insecure)

exportfs -ra

systemctl start nfs-kernel-server

挂载 mount -t nfs xx.xx.xx.xx:/data/nfs ./nfs

umount /mnt/nfs

## helm安装 nfs provisioner

helm repo add nfs-subdir-external-provisioner https://kubernetes-sigs.github.io/nfs-subdir-external-provisioner/

helm install nfs-subdir-external-provisioner nfs-subdir-external-provisioner/nfs-subdir-external-provisioner \
    --set nfs.server=xx.xx.xx.xx \
    --set nfs.path=/data/nfs

## 其他
如果希望能够通过本地 IP 地址和端口访问 Kubernetes 服务，有几种方法可以实现：

1. 使用 type=NodePort：
    1. minikube service <service name>
2. 使用 type=LoadBalancer：
    1. minikube tunnel
    2. kubectl get svc 查看 EXTERNAL-IP 
    3. 访问 EXTERNAL-IP:port 8080:30791/TCP port为8080
3. 使用 type=NodePort + Ingress：
    1. minikube addons enable ingress
    2. kubectl apply -f ingress-nginx.yml
    3. minikube tunnel 开启ip隧道
    4. curl --resolve "hello-k8s.info:80:127.0.0.1" -i http://hello-k8s.info/app/hello
    5. Or add a line to the bottom of the /etc/hosts file on your computer 127.0.0.1 hello-k8s.info curl -i http://hello-k8s.info.info/app/hello 测试HPA：while true; do curl http://hello-k8s.info/app/hello >/dev/null 2>&1; done

### Pod调度

1. 调度约束（Scheduling Constraints） 例如

2. 节点选择器

节点标签：
* 节点1：拥有标签 zone=us-west
* 节点2：拥有标签 zone=us-east
* 节点3：拥有标签 zone=eu-central

示例

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
spec:
  containers:
  - name: nginx
    image: nginx
  nodeSelector:   # 节点选择器
    zone: us-west   # 选择具有标签 zone=us-west 的节点
  resources:   # 资源限制
    limits:
      cpu: "500m"   # CPU 限制为 0.5 核
      memory: "256Mi"   # 内存 限制为 256Mi

```

3. 亲和力与反亲和力
4. 污点和容忍

示例

假设我们有一个 Kubernetes 集群，其中一个节点（node-1）上有一个污点 special=true:NoSchedule (NoSchedule/PreferNoSchedule/NoExecute)，表示这个节点是一个特殊节点，不希望调度一般的 Pod 

打上污点

```
kubectl taint nodes node-1 special=true:NoSchedule
```

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: special-pod
spec:
  containers:
  - name: nginx
    image: nginx
  tolerations:   # 容忍
  - key: "special"
    operator: "Equal"
    value: "true"
    effect: "NoSchedule"

```