# Hello custom controller

## 準備

kubernetes v1.12.3を使う。minikubeで環境を作成

```
$ minikube start --memory=8192 --cpus=4 \
    --kubernetes-version=v1.12.3 \
    --vm-driver=hyperkit \
    --bootstrapper=kubeadm
```

## Custom Resource

- [Custom Resources](https://v1-12.docs.kubernetes.io/docs/concepts/extend-kubernetes/api-extension/custom-resources/)

CustomResourceは独自に作ったPodのようなリソースのこと。CustomResourceDefinition(CDR)を使って作成する。

* [Extend the Kubernetes API with CustomResourceDefinitions](https://v1-12.docs.kubernetes.io/docs/tasks/access-kubernetes-api/custom-resources/custom-resource-definitions/)

```yaml
# https://github.com/kubernetes/sample-controller
apiVersion: apiextensions.k8s.io/v1beta1
kind: CustomResourceDefinition
metadata:
  # name must match the spec fields below, and be in the form: <plural>.<group>
  name: foos.samplecontroller.k8s.io
spec:
  # group name to use for REST API: /apis/<group>/<version>
  group: samplecontroller.k8s.io
  # version name to use for REST API: /apis/<group>/<version>
  version: v1alpha
  # either Namespaced or Cluster
  scope: Namespaced
  names:
    # plural name to be used in the URL: /apis/<group>/<version>/<plural>
    plural: foos
    # singular name to be used as an alias on the CLI and for display
    singular: foo
    # kind is normally the CamelCased singular type. Your resource manifests use this.
    kind: Foo
    # shortNames allow shorter string to match your resource on the CLI
    shortNames:
      - foo
```

CustomResourceDefinitionを適用すると、以下のようにcrdオブジェクトとして見えるようになる：

```
$ kubectl get crd
NAME                           AGE
foos.samplecontroller.k8s.io   34m
```

Fooオブジェクトを作ってみる。以下のようなyamlを用意し、 `kubectl apply`する。

```yaml
# https://github.com/kubernetes/sample-controller
apiVersion: samplecontroller.k8s.io/v1alpha
kind: Foo
metadata:
  name: foo-001
spec:
  deploymentName: deploy-foo-001
  replicas: 1
---
apiVersion: samplecontroller.k8s.io/v1alpha
kind: Foo
metadata:
  name: foo-002
spec:
  deploymentName: deploy-foo-002
  replicas: 1
```

fooオブジェクトが生まれる：

```
$ kubectl get foos
NAME      AGE
foo-001   50s
foo-002   50s
```

詳細情報を見ると、yamlで指定したspecが見えることが確認できる：

```
$ kubectl describe foo foo-001
Name:         foo-001
Namespace:    default
Labels:       <none>
Annotations:  kubectl.kubernetes.io/last-applied-configuration={"apiVersion":"samplecontroller.k8s.io/v1alpha","kind":"Foo","metadata":{"annotations":{},"name":"foo-001","namespace":"default"},"spec":{"deploymentNa...
API Version:  samplecontroller.k8s.io/v1alpha
Kind:         Foo
Metadata:
  Creation Timestamp:  2019-01-02T03:42:45Z
  Generation:          1
  Resource Version:    3649
  Self Link:           /apis/samplecontroller.k8s.io/v1alpha/namespaces/default/foos/foo-001
  UID:                 77506f4f-0e40-11e9-95ac-263ada282756
Spec:
  Deployment Name:  deploy-foo-001
  Replicas:         1
Events:             <none>
```



## Custom Controller

CRDだけではオブジェクトが作成されるだけで、何も起きないため、このオブジェクトを参照して何かするためのコントローラーを作らないといけない。

* [KubernetesのCRD(Custom Resource Definition)とカスタムコントローラーの作成](https://qiita.com/__Attsun__/items/785008ef970ad82c679c)
* [Extending Kubernetes with Custom Resources and Operator Frameworks](https://speakerdeck.com/ianlewis/extending-kubernetes-with-custom-resources-and-operator-frameworks)
  * Kubernetesを拡張するにはデータとロジックが必要
  * データはCustomResourceDefinition
  * ロジックはController
  * Operator Framework
    * operator-sdk
    * kubebuilder
* [kubernetes/sample-controller](https://github.com/kubernetes/sample-controller)
  * `Foo` というカスタムリソースを定義する
  * このリソースは `Deployment` を定義するためのカスタムリソース
    * 名前とレプリカ数を指定できる
  * client-goライブラリを使っている

## sample-controllerを素朴に実装してみる

* [kubernetes/code-generator](https://github.com/kubernetes/code-generator)
* [Kubernetesを拡張しよう](https://www.ianlewis.org/jp/extending-kubernetes-ja)
* [Extending Kubernetes: Create Controllers for Core and Custom Resources](https://medium.com/@trstringer/create-kubernetes-controllers-for-core-and-custom-resources-62fc35ad64a3)
  * コントローラのイベントフロー解説
* [KubernetesのCustom Resource Definition(CRD)とCustom Controller](https://www.sambaiz.net/article/182/)
* [Kubernetes Deep Dive: Code Generation for CustomResources](https://blog.openshift.com/kubernetes-deep-dive-code-generation-customresources/)



```
$ env GO111MODULE=off bash ~/src/k8s.io/code-generator/generate-groups.sh all github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/client github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/apis foo:v1alpha
Generating deepcopy funcs
Generating clientset for foo:v1alpha at github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/client/clientset
Generating listers for foo:v1alpha at github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/client/listers
Generating informers for foo:v1alpha at github.com/takaishi/hello2018/hello-custom-resource/my-sample-controller/pkg/client/informers
```



```
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/clientset.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/doc.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/fake/clientset_generated.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/fake/doc.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/fake/register.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/scheme/doc.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/scheme/register.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/typed/foo/v1alpha/doc.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/typed/foo/v1alpha/fake/doc.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/typed/foo/v1alpha/fake/fake_foo_client.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/typed/foo/v1alpha/foo_client.go
? hello-custom-resource/my-sample-controller/pkg/client/clientset/versioned/typed/foo/v1alpha/generated_expansion.go
```



```
cd my-sample-controller
env GO111MODULE=off OOS=linux GOARCH=amd64 go build -o controller-main main.go
docker build . -t rtakaishi/sample-controller-main
docker push rtakaishi/sample-controller-main:latest
```







