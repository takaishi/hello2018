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

## Custom Controller

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







