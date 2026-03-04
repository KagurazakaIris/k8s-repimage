# repimage

很多镜像都在国外。国内下载很慢，需要加速，每次都要手动修改yaml文件中的镜像地址，很麻烦。这个项目就是为了解决这个问题。

用于替换k8s中一些在国内无法访问的镜像地址，替换的镜像地址在 [public-image-mirror
](https://github.com/DaoCloud/public-image-mirror)中查看

# 快速上手
## 安装
```shell
kubectl create -f https://cdn.jsdelivr.net/gh/KagurazakaIris/k8s-repimage@refs/heads/main/repimage.yaml
kubectl rollout status deployment/repimage -n kube-system
```

# 使用后效果
自动替换yaml文件中的镜像地址，例如:
```
k8s.gcr.io/coredns/coredns => m.daocloud.io/k8s.gcr.io/coredns/coredns

nginx => m.daocloud.io/docker.io/library/nginx
```

# 配置选项
## 忽略指定域名
如果你有一些私有镜像仓库或者不需要加速的域名，可以通过 `--ignore-domains` 参数来忽略这些域名。

例如，在 deployment.yaml 中添加参数：
```yaml
containers:
  - command:
      - /repimage
      - --ignore-domains=myregistry.example.com,private.registry.local
```

这样，来自 `myregistry.example.com` 和 `private.registry.local` 的镜像将不会被替换。

## 自定义镜像前缀

[内网再部署一级缓存](https://github.com/DaoCloud/public-image-mirror/tree/main/docs/local-cache)

当前默认使用 `m.daocloud.io` 作为镜像前缀，可以通过 `--prefix` 参数自定义：

```yaml
containers:
  - command:
      - /repimage
      - --prefix=mirror.example.com
```

## 自定义镜像源映射（推荐）

除了使用 `--prefix` 参数进行统一前缀替换外，你还可以通过提供一个 JSON 配置文件来定义更灵活的镜像源映射规则。这使得你可以将特定的原始域名映射到不同的目标注册中心。

默认情况下，工具会尝试加载 `./config/registries.json` 文件。你可以使用 `--config` 参数指定自定义的配置文件路径。

配置文件 `registries.json` 的格式如下：

```json
{
  "docker.io": "docker.utilapi.bid",
  "quay.io": "quay.utilapi.bid",
  "gcr.io": "gcr.utilapi.bid",
  "k8s.gcr.io": "k8s-gcr.utilapi.bid",
  "registry.k8s.io": "k8s.utilapi.bid",
  "ghcr.io": "ghcr.utilapi.bid",
  "docker.cloudsmith.io": "cloudsmith.utilapi.bid"
}
```

在此配置中，键是原始镜像的域名，值是用于替换的新的注册中心地址。例如，`"docker.io": "docker.utilapi.bid"` 会将所有 `docker.io` 域名下的镜像替换为 `docker.utilapi.bid`。

**优先顺序：**

如果同时提供了 `--prefix` 和 `--config`，并且配置文件中定义了某个域名的映射，则配置文件中的映射会优先于 `--prefix` 进行替换。如果没有找到特定的域名映射，则会回退到 `--prefix` 行为。

**使用示例：**

在 `deployment.yaml` 中添加 `--config` 参数：

```yaml
containers:
  - command:
      - /repimage
      - --config=/etc/repimage/registries.json
```

通过这种方式，你可以方便地管理和升级你的镜像源配置，而无需修改代码。

## 忽略特定Pod
如果你希望某些Pod不被镜像替换，可以在Pod的annotation中添加 `repimage.kubernetes.io/skip=true`。

例如：
```yaml
apiVersion: v1
kind: Pod
metadata:
  name: example-pod
  annotations:
    repimage.kubernetes.io/skip: "true"
spec:
  containers:
  - name: example-container
    image: nginx
```

带有此annotation的Pod将不会被镜像替换。

# License

Apache-2.0

# 特别感谢

- [@shixinghong](https://github.com/shixinghong) 感谢原作者提供的灵感
- [DaoCloud](https://github.com/DaoCloud) 免费提供的镜像代理服务
