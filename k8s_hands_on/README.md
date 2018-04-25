# Kubernetes ハンズオン(1時間で)

## 参考資料

- https://kubernetes.io/
- [O'Reilly Japan - 入門 Kubernetes](https://www.oreilly.co.jp/books/9784873118406/)

## k8s環境を用意する

まずはk8s環境が必要です。クラウド上にクラスターを作ってもよいのですが、今回は手元に用意しましょう。
minikubeという、1VM上にk8s環境を構築してくれるツールを使います。

- https://github.com/kubernetes/minikube

### minikubeのインストール

Homebrewを使えばシュッとインストールできます。

```
brew cask install minikube
```


### minikubeを使ってローカルにk8s環境を作る

`minikube start`コマンドで作ってくれます。VMを作る方法として、今回はVirtualBoxを使います。
Hyperkitやxhyveを使うことも可能です。
メモリサイズは各自の環境に応じて変更してもよいでしょう。

```
minikube start \
    --vm-driver=virtualbox \
    --kubernetes-version=v1.9.4 \
    --bootstrapper=localkube \
    --memory=4096
```

## コンテナを用いてWordpressを起動するには

いよいよk8s上にwordpressを作っていきます。今回は、store.docker.comにあるコンテナイメージを使います。

* mysql
* wordpress

## MySQL用のパスワードを機密情報としてk8sに登録する
いきなりmysqlを建てるまえに、準備として機密情報を扱えるようにしておきます。機密情報といってもいろいろありますが、最低限必要なのはwordpressからmysqlにアクセスするためのパスワードです。
mysqlを起動した時にパスワードを設定し、wordpressを起動したときにそのパスワードでmysqlに接続する、というやり方です。
この2つのコンテナに同じパスワード情報を渡すにはどうすればよいでしょうか？コンテナを起動する際、引数や環境変数で渡すという方法も考えられますが、管理が難しいですよね。
k8sには、ConfigMap、Secretという機能があり、ここに様々な情報を保存しておくことができます。今回はSecretを使ってパスワードをk8s上に保存し、mysqlコンテナとwordpressコンテナから利用できるようにしておきます。

```
$ cat password.txt
fhui2wqofwheaiuwhel
```

```
$ tr -d '\n' <password.txt >.strippedpassword.txt && mv .strippedpassword.txt password.txt   # bash / zsh
$ tr -d '\n' <password.txt >.strippedpassword.txt; and mv .strippedpassword.txt password.txt # fish
```

## mysqlを起動する

次はmysqlを建てていいきます。k8sでは、リソースの最小単位はコンテナではなくPodと呼びます。Podはコンテナとストレージボリュームの集合です。
さて、以下のyamlがmysqlのPodを定義したものです。

* apiVersion
* kind
	* リソースの種類です。
* metadata
	* Podに設定するメタデータです。この例では、Podの名前とラベルを設定しています。ラベルを設定しておくことで、同じアプリケーションが動作するPodが複数存在するときに絞り込みを行うことができます。
* spec
	* Podの仕様についての定義です。この例ではイメージは何を使うか、起動時に渡す環境変数は何か、公開するポートが何かを設定しています。
	* 環境変数として`MYSQL_ROOT_PASSWORD`を定義しています。ここで、シークレットに登録したパスワードを参照しています。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: wordpress-mysql
  labels:
    app: wordpress
spec:
  containers:
  - image: mysql:5.6
    name: mysql
    env:
      - name: MYSQL_ROOT_PASSWORD
        valueFrom:
          secretKeyRef:
            name: mysql-pass
            key: password.txt
    ports:
      - containerPort: 3306
        name: mysql
```

applyコマンドでPodを作成します。


```
➤ kubectl apply -f ./manifests/mysql.yaml
```


## wordpressを起動する

mysqlを建てたりので、次はwordpressを建てていきましょう。mysqlを同じく、Podを定義します。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: wordpress
  labels:
    app: wordpress
spec:
  containers:
  - image: wordpress
    name: wordpress
    env:
    - name: WORDPRESS_DB_HOST
      value: wordpress-mysql
    - name: WORDPRESS_DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: mysql-pass
          key: password.txt
    ports:
    - containerPort: 80
      name: wordpress
```

これだけだとエラーになる。

```
➤ kubectl logs po/wordpress
WordPress not found in /var/www/html - copying now...
Complete! WordPress has been successfully copied to /var/www/html

Warning: mysqli::__construct(): (HY000/2002): php_network_getaddresses: getaddrinfo failed: Temporary failure in name resolution in Standard input code on line 22

MySQL Connection Error: (2002) php_network_getaddresses: getaddrinfo failed: Temporary failure in name resolution
```

なぜか？wordpressからwordpress-mysqlに繋ぐことができないから。

なぜでしょうか？気づいているかたもいるかもしれませんが、wordpressからwordpress-mysqlに繋ぐことができないからです。Podだけでは、他から参照することができないんですね。ここで必要となるのは、mysqlの接続先情報はこれだ！というのを知っている、Serviceというリソースを作ることです。Serviceを作ることで外部からディスカバリーできるようになるわけですね。

Serviceの定義を書いてapplyしましょう。Serviceのspecとしては、ポート情報とどのPodにつなぐかを決めるselector、クラスターIPです。


```
---
apiVersion: v1
kind: Service
metadata:
  name: wordpress-mysql
  labels:
    app: wordpress
    tier: mysql
spec:
  ports:
    - port: 3306
  selector:
    app: wordpress
    tier: mysql
  clusterIP: None
```

Serviceを作ったら、wordpress Podを一度消します。
そして、mysqlのホストとしてmysql Serviceを使うように定義を追加し、再度applyしましょう。

これで、wordpressからwordpress-mysqlに繋がるはずです。しかし、k8sの外からwordpressに繋ぐことができません。そこで、wordpress用のサービスも作ります。mysqlとは違い、外部から接続したいので、typeとしてNodePortというものを選択しています。これは、PodのポートとNodeのポートを接続するというものです。

```
---
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
  name: wordpress
spec:
  type: NodePort
  ports:
    - port: 80
      targetPort: 80
      nodePort: 30180
  selector:
    app: wordpress
    tier: frontend
```

これをapplyしたら、ブラウザからアクセス可能となります。minikubeの場合、以下のコマンドでブラウザで開くことができます。Wordpressのセットアップ画面が表示されれば成功です。

```
minikube service wordpress
```

これでwordpressが起動しました！ここからは、より便利な機能（といってもk8sを使うならほぼ必須ですが）を使ってみましょう。

## アプリケーションがクラッシュした時、自動回復してほしい

さてさて無事にwordpressが起動したのですが、このままだと例えばwordpressコンテナがクラッシュした時に繋がらなくなってしまいます。Podはコンテナとストレージボリュームの集合というだけで、自分自身を管理するということをしていないためです。そこで、Deploymentというリソースを使います。Podとして定義していた箇所を、以下のようにDeploymentにします。Deploymentのspec/template/spec部分はPodのspecと同じであることがわかりますね。

```
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wordpress
  labels:
    app: wordpress
    tier: frontend
spec:
  selector:
    matchLabels:
      app: wordpress
      tier: frontend
  template:
    metadata:
      labels:
        app: wordpress
        tier: frontend
    spec:
      containers:
      - image: wordpress
        name: wordpress
        env:
        - name: WORDPRESS_DB_HOST
          value: wordpress-mysql:3306
        - name: WORDPRESS_DB_PASSWORD
          valueFrom:
            secretKeyRef:
              name: mysql-pass
              key: password.txt
        ports:
        - containerPort: 80
          name: wordpress
```

applyします。Deploymentによって自動的にwodpress Podが作成されることが確認できます。古いPod(wordpress)は削除しておきましょう。

```
➤ kubectl apply -f ./manifests/wordpress.yaml
deployment "wordpress" created
service "wordpress" unchanged
r_takaishi@PMC02V437VHV2R:~/s/g/t/h/wordpress_on_k8s|wordpress_on_k8s⚡?
➤ kubectl delete po/wordpress
pod "wordpress" deleted
```

Deploymentは定義した状態を維持しようとします。例えば、Podを削除しても新しいPodが自動的に作成され、Podが1台動いているという状態が維持されるわけです。

```
➤ kubectl delete po/wordpress-55448464cd-7hgdk
```

k8sの大きな特徴として、この状態を維持しようとする機能があります。これにより、アプリケーションがクラッシュしたりしても自動的に復旧させることができ、サービスの運用負担軽減に繋げられるというわけです。

Deploymentを使うように変えてもPodは1つのままです。しかし、Podの数を変更することももちろん可能なので、試してみましょう。


## DBのデータが永続ボリュームに保存されていない問題を解決する

現在、MysqlのDBデータはPod内にあります。そのため、Podを削除するとデータが消えてしまう。これはまずいですね。もちろん解決方法があって、永続ボリューム（PersistentVolume）を使うことができます。

まず、PersistentVolumeを定義します。

```
apiVersion: v1
kind: PersistentVolume
metadata:
  name: local-volume-1
  labels:
    type: local
spec:
  capacity:
    storage: 20Gi
  accessModes:
    - ReadWriteOnce
  hostPath:
    path: /tmp/data/lv-1
  persistentVolumeReclaimPolicy: Recycle
```

適用すると、PersistentVolumeが作成されます。

```
➤ kubectl apply -f ./manifests/volume.yaml
persistentvolume "local-volume-1" created
```

mysql PodがPersistentVolumeを使うようにマウント設定などを行います。ついでにDeployment化しておきましょう。

```
---
apiVersion: apps/v1
kind: Deployment
metadata:
  name: wordpress-mysql
  labels:
    app: wordpress
    tier: mysql
spec:
  selector:
    matchLabels:
      app: wordpress
      tier: mysql
  template:
    metadata:
      labels:
        app: wordpress
        tier: mysql
    spec:
      containers:
      - image: mysql:5.6
        name: mysql
        env:
          - name: MYSQL_ROOT_PASSWORD
            valueFrom:
              secretKeyRef:
                name: mysql-pass
                key: password.txt
        ports:
          - containerPort: 3306
            name: mysql
        volumeMounts:
          - name: mysql-local-storage
            mountPath: /var/lib/mysql
      volumes:
      - name: mysql-local-storage
        persistentVolumeClaim:
          claimName: mysql-lv-claim
---
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-lv-claim
  labels:
    app: wordpress
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```


## まとめ

以下の5リソースを使ってwordpress環境を作ってきました。k8sはyaml形式の定義が長くて難しい、という話も聞きますが、1つ1つ分解するとリソースを組み合わせているということがわかります。

* ConfigMap / Secret
* Pod
* Service
* Deployment
* PersistentVolume
