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







