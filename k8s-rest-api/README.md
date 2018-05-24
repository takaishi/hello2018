# KubernetesのNamespaceをRuby DSLで管理してみる



```
⟩ find manifests -type f | xargs -I[] cat []
apiVersion: v1
kind: Namespace
metadata:
  name: foo002
  labels:
    name: production
    service: piyo
apiVersion: v1
kind: Namespace
metadata:
  name: foo004
apiVersion: v1
kind: Namespace
metadata:
  name: foo001
  labels:
    name: staging
    service: hoge
```

```
⟩ kubectl apply -f ./manifests
namespace "foo001" created
namespace "foo002" created
namespace "foo004" created
```

```
⟩ kubectl get namespace --show-labels
NAME          STATUS    AGE       LABELS
default       Active    32d       <none>
foo001        Active    10s       name=staging,service=hoge
foo002        Active    10s       name=production,service=piyo
foo004        Active    10s       <none>
kube-public   Active    32d       <none>
kube-system   Active    32d       <none>
```

```
⟩ bundle exec ruby ./test.rb
foo001にラベル test: 123を追加
foo001にラベル test2: 456を追加
foo001からラベル service: hogeを削除
foo002のラベルを更新 service: aaaaaaaaaaa
create: foo003
delete foo004
```

```
⟩ kubectl get namespace --show-labels
NAME          STATUS    AGE       LABELS
default       Active    32d       <none>
foo001        Active    35s       name=staging,test2=456,test=123
foo002        Active    35s       name=production,service=aaaaaaaaaaa
foo003        Active    10s       name=development,service=fuga
kube-public   Active    32d       <none>
kube-system   Active    32d       <none>
```