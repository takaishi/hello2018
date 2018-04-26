# Kubernetes ハンズオン(1時間で)

## ゴール

- ローカルにKubernetes環境を構築し、触れるようになる
- Kubernetes上でWordpressを動かし、以下のリソースを使ってみる
  - ConfigMap / Secret
  - Pod
  - Service
  - Deployment
  - PersistentVolume


## k8s環境を用意する

まずはk8s環境が必要です。クラウド上にクラスターを作ってもよいのですが、今回は手元に用意しましょう。
[minikube](https://github.com/kubernetes/minikube)という、1VM上にk8s環境を構築してくれるツールを使います。

### minikubeのインストール

2018年04月25日現在ではminikubeの最新バージョンはv0.26.1なのですが、今回はv0.25.2を使います。curlを使って実行ファイルをダウンロードし、権限付与を行います。

```
➤ curl -L https://github.com/kubernetes/minikube/releases/download/v0.25.2/minikube-darwin-amd64 -o /path/to/minikube
➤ chmod +x /path/to/minikube
```

```
➤ minikube version
minikube version: v0.25.2
```

### kubectlのインストール

Kubernetesを操作するため、コマンドラインツールが必要です。`kubectl`というツールが提供されており、Homebrewでインストールできます。

```
➤ brew install kubectl
```

### minikubeを使ってローカルにk8s環境を作る

`minikube start`コマンドを使います。Kubernetesのバージョンは、`v1.9.4`を指定します。メモリサイズは各自の環境に応じて変更してもよいでしょう。

※ 2018年04月25日時点で、`bootstrapper`オプションを使ってlocalkubeを指定しておかないとうまく起動しないようです。

```
minikube start \
    --vm-driver=virtualbox \
    --kubernetes-version=v1.9.4 \
    --bootstrapper=localkube \
    --memory=4096
```



```
➤ minikube start \
      --vm-driver=virtualbox \
      --kubernetes-version=v1.9.4 \
      --bootstrapper=localkube \
      --memory=4096
There is a newer version of minikube available (v0.26.1).  Download it here:
https://github.com/kubernetes/minikube/releases/tag/v0.26.1

To disable this notification, run the following:
minikube config set WantUpdateNotification false
Starting local Kubernetes v1.9.4 cluster...
Starting VM...
Downloading Minikube ISO
 142.22 MB / 142.22 MB [============================================] 100.00% 0s
Getting VM IP address...
Moving files into cluster...
Setting up certs...
Connecting to cluster...
Setting up kubeconfig...
Starting cluster components...
Kubectl is now configured to use the cluster.
Loading cached images from config file.
```

`minikube status`コマンドで、`minikube`と`cluster`がRunningになっていればOKです。

```
➤ minikube status
minikube: Running
cluster: Running
kubectl: Correctly Configured: pointing to minikube-vm at 192.168.99.100
```

kubectlがminikubeにアクセスできていることも確認します。`kubectl get nodes`コマンドで、Nodeとして`minikube`が1台いればOKです。

```
➤ kubectl get nodes
NAME       STATUS    ROLES     AGE       VERSION
minikube   Ready     <none>    1m        v1.9.4
```

## コンテナを用いてWordpressを起動するには

いよいよk8s上にwordpressを作っていきます。今回は、store.docker.comにあるコンテナイメージを使います。

* [mysql](https://store.docker.com/images/mysql)
* [wordpress](https://store.docker.com/images/wordpress)

## MySQL用のパスワードを機密情報としてk8sに登録する
いきなりmysqlを建てるまえに、準備として機密情報を扱えるようにしておきます。
機密情報といってもいろいろありますが、最低限必要なのはwordpressからmysqlにアクセスするためのパスワードです。
mysqlを起動した時にパスワードを設定し、wordpressを起動したときにそのパスワードでmysqlに接続する、というやり方です。
この2つのコンテナに同じパスワード情報を渡すにはどうすればよいでしょうか？コンテナを起動する際、引数や環境変数で渡すという方法も考えられますが、管理が難しいですよね。
k8sには、ConfigMap、Secretという機能があり、ここに様々な情報を保存しておくことができます。今回はSecretを使ってパスワードをk8s上に保存し、mysqlコンテナとwordpressコンテナから利用できるようにしておきます。

```
➤ cat password.txt
fhui2wqofwheaiuwhel
```

`password.txt`から改行文字を消します。

```
➤ tr -d '\n' <password.txt >.strippedpassword.txt && mv .strippedpassword.txt password.txt   # bash / zsh
➤ tr -d '\n' <password.txt >.strippedpassword.txt; and mv .strippedpassword.txt password.txt # fish
```

`kubectl ceate secret`コマンドで登録します。

```
➤ kubectl create secret generic mysql-pass --from-file=password.txt
```

`kubectl get secrets`コマンドで、Secretの一覧を取得することができます。

```
➤ kubectl get secrets
NAME                  TYPE                                  DATA      AGE
default-token-v2tf5   kubernetes.io/service-account-token   3         7m
mysql-pass            Opaque                                1         1s
```

## mysqlを起動する

次はmysqlを起動します。Kubernetesでは、リソースの最小単位はコンテナではなくPodというものです。Podはコンテナとストレージボリュームの集合です。

`./manifests/mysql-pod.yaml`という名前で、以下のyamlを作成します。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: wordpress-mysql
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
```

- apiVersion
- kind
  - リソースの種類です。
- metadata
  - Podに設定するメタデータです。この例では、Podの名前とラベルを設定しています。ラベルを設定しておくことで、同じアプリケーションが動作するPodが複数存在するときに絞り込みを行うことができます。
- spec
  - Podの仕様についての定義です。この例ではイメージは何を使うか、起動時に渡す環境変数は何か、公開するポートが何かを設定しています。
  - 環境変数として`MYSQL_ROOT_PASSWORD`を定義しています。ここで、シークレットに登録したパスワードを参照しています。

`kubectl apply`コマンドでPodを作成します。


```
➤ kubectl apply -f ./manifests/mysql-pod.yaml
```

`kubectl get pods`コマンドで、mysqlが起動していることを確認しましょう。`-l`ラベルを使って絞り込みをするためのオプションです。

```
➤ kubectl get pods -l app=wordpress -l tier=mysql
NAME              READY     STATUS    RESTARTS   AGE
wordpress-mysql   1/1       Running   0          4m
```

mysqlのログも確認してみましょう。`kubectl logs`コマンドを使います。

```
➤ kubectl logs wordpress-mysql | tail -n 5
2018-04-25 02:20:59 1 [Warning] Insecure configuration for --pid-file: Location '/var/run/mysqld' in the path is accessible to all OS users. Consider choosing a different directory.
2018-04-25 02:20:59 1 [Warning] 'proxies_priv' entry '@ root@wordpress-mysql' ignored in --skip-name-resolve mode.
2018-04-25 02:20:59 1 [Note] Event Scheduler: Loaded 0 events
2018-04-25 02:20:59 1 [Note] mysqld: ready for connections.
Version: '5.6.40'  socket: '/var/run/mysqld/mysqld.sock'  port: 3306  MySQL Community Server (GPL)
```

`mysqld: ready for connections.`というログが見えますね。さらに、mysqlクライアントから繋がるかどうかも試してみましょう。

`kubectl port-forward`コマンドを使うことで、Podのポートをフォワードできます。

```
➤ kubectl port-forward wordpress-mysql 13306:3306
Forwarding from 127.0.0.1:13306 -> 3306
Handling connection for 13306
```

最後に、mysqlコマンドで`127.0.0.1:13306`に繋がるかチェックしてみます。パスワードはSecretとして登録したものを使います。

```
➤ mysql -s -uroot -p -h127.0.0.1 --port=13306
Enter password:
mysql> show databases;
Database
information_schema
mysql
performance_schema
```

つながりました！やったー！

## DBのデータが永続ボリュームに保存されていない問題を解決する

さて、MysqlのDBデータはPod内にあります。そのため、Podを削除するとデータが消えてしまう。これはまずいですね。もちろん解決方法があって、永続ボリューム（PersistentVolume）を使うことができます。

まず、PersistentVolumeを定義します。`./manifests/mysql-volume.yaml`という名前で、以下のyamlを作成します。

```yaml
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
➤ kubectl apply -f ./manifests/mysql-volume.yaml
persistentvolume "local-volume-1" created
```

PersistentVolume一覧を見てみましょう。local-volume-1が作成されているはずです。

```
➤ kubectl get persistentvolume
NAME             CAPACITY   ACCESS MODES   RECLAIM POLICY   STATUS      CLAIM     STORAGECLASS   REASON    AGE
local-volume-1   20Gi       RWO            Recycle          Available                                      8s
```

次に、PodがPersistentVolumeを取得できるようにします。これは、PersistentVolumeClaimというリソースを使います。

`./manifests/mysql-persistent-volume-claim.yaml`という名前で、以下のyamlを作成します。

```yaml
apiVersion: v1
kind: PersistentVolumeClaim
metadata:
  name: mysql-lv-claim
  labels:
    app: wordpress
    tier: mysql
spec:
  accessModes:
    - ReadWriteOnce
  resources:
    requests:
      storage: 20Gi
```

`kubectl apply`コマンドでPersistentVolumeClaimを作成します。

```
➤ kubectl apply -f ./manifests/mysql-persistent-volume-claim.yaml
```

最後に、mysql PodがPersistentVolumeを使うようにマウント設定を行います。ついでにDeployment化しておきましょう。

`./manifests/mysql-deployment-with-volume.yaml`という名前で、以下のyamlを作成します。

```yaml
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
```

`kubectl apply`して、mysql Deploymentを作ります。

```
➤ kubectl apply -f  manifests/mysql-deployment-with-volume.yaml
```

Deploymentを作ったので、Deploymentによって作られたPodと、最初に作成したPodがいる状態です。`wordpress-mysql-cf9449df-7kfhf` がDeploymentによって作られたPodですね。

```
➤ kubectl get pod -l app=wordpress -l tier=mysql
NAME                             READY     STATUS    RESTARTS   AGE
wordpress-mysql                  1/1       Running   0          19m
wordpress-mysql-cf9449df-7kfhf   1/1       Running   0          31s
```

古いwordpress-mysql Podは消しておきましょう。

```
➤ kubectl delete pod wordpress-mysql
```

Deploymentが作ったPodを見ると、Volumesにmysql-lv-claimが登録されていることがわかります。

```
➤ kubectl describe pod wordpress-mysql-cf9449df-7kfhf
Name:           wordpress-mysql-cf9449df-7kfhf

[snip]

Volumes:
  mysql-local-storage:
    Type:       PersistentVolumeClaim (a reference to a PersistentVolumeClaim in the same namespace)
    ClaimName:  mysql-lv-claim
    ReadOnly:   false
  default-token-v2tf5:
    Type:        Secret (a volume populated by a Secret)
    SecretName:  default-token-v2tf5
    Optional:    false

[snip]
```

## wordpressを起動する

mysqlを起動したので、次はwordpressです。mysqlと同じように、Podを定義します。

`./manifests/wordpress-pod.yaml`という名前で、以下のyamlを作成します。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: wordpress
  labels:
    app: wordpress
    tier: front-end
spec:
  containers:
  - image: wordpress
    name: wordpress
    env:
    - name: WORDPRESS_DB_PASSWORD
      valueFrom:
        secretKeyRef:
          name: mysql-pass
          key: password.txt
    ports:
    - containerPort: 80
      name: wordpress
```

`kubectl apply`します。

```
➤ kubectl apply -f ./manifests/wordpress-pod.yaml
```

起動しているように見えますが…

```
➤ kubectl get pods -l app=wordpress -l tier=front-end
NAME        READY     STATUS    RESTARTS   AGE
wordpress   1/1       Running   0          40s
```

ログを見ると、mysqlサーバーに繋げないというエラーになっています。

```
➤ kubectl logs wordpress | tail -n 5
Warning: mysqli::__construct(): php_network_getaddresses: getaddrinfo failed: Name or service not known in Standard input code on line 22

Warning: mysqli::__construct(): (HY000/2002): php_network_getaddresses: getaddrinfo failed: Name or service not known in Standard input code on line 22

MySQL Connection Error: (2002) php_network_getaddresses: getaddrinfo failed: Name or service not known
```

なぜか？wordpressのPod定義で、mysqlのホストを指定していないからです。指定するには、環境変数として、`WORDPRESS_DB_HOST`を渡してあげる必要があります。

さて、環境変数でmysqlのホストを指定できることはわかりました。では、ここで必要となるのは、mysqlの接続先情報はこれだ！というのを知っている、Serviceというリソースを作ることです。Serviceを作ることで外部からディスカバリーできるようになるわけですね。

mysql用Serviceの定義を書きましょう。Serviceのspecとしては、ポート情報とどのPodにつなぐかを決めるselector、クラスターIPです。

`./manifests/mysql-service.yaml`という名前で、以下のyamlを作成します。


```yaml
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

そして、`kubectl apply`してwordpress-mysqlサービスを作成します。

```
➤ kubectl apply -f manifests/mysql-service.yaml
```

Service一覧を見て、作成されていることを確認します。

```
➤ kubectl get service -l app=wordpress -l tier=mysql
NAME              TYPE        CLUSTER-IP   EXTERNAL-IP   PORT(S)    AGE
wordpress-mysql   ClusterIP   None         <none>        3306/TCP   49s
```

Serviceを作ったら、wordpress Podを一度消します。

```
➤ kubectl delete pod wordpress
```

Podが消えるまで、少し待つ必要があります。削除されたのを確認して、次に進みましょう。

```
➤ kubectl get pods -l app=wordpress -l tier=frontend
No resources found.
```

mysqlのホストとしてmysql Serviceを使うように定義を追加ます。spec/containers/envに、`WORDPRESS_DB_HOST`を追加していて、値としては、wordpress-mysql Serviceの3306を指定しています。

`./manifests/wordpress-pod-with-mysql-host.yaml`という名前で、以下のyamlを作成します。

```yaml
apiVersion: v1
kind: Pod
metadata:
  name: wordpress
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

`kubectl apply`して、再度Podを作成します。

```
➤ kubectl apply -f ./manifests/wordpress-pod-with-mysql-host.yaml
```

Podが起動したら、ログを見てみましょう。先ほどとは違い、エラーがでていないことがわかります。

```
➤ kubectl logs wordpress | tail -n 5
Complete! WordPress has been successfully copied to /var/www/html
AH00558: apache2: Could not reliably determine the server's fully qualified domain name, using 172.17.0.5. Set the 'ServerName' directive globally to suppress this message
AH00558: apache2: Could not reliably determine the server's fully qualified domain name, using 172.17.0.5. Set the 'ServerName' directive globally to suppress this message
[Wed Apr 25 02:47:15.619378 2018] [mpm_prefork:notice] [pid 1] AH00163: Apache/2.4.25 (Debian) PHP/7.2.4 configured -- resuming normal operations
[Wed Apr 25 02:47:15.619793 2018] [core:notice] [pid 1] AH00094: Command line: 'apache2 -D FOREGROUND'
```

これで、wordpressからwordpress-mysqlに繋がるはずです。しかし、k8sの外からwordpressに繋ぐことができません。そこで、wordpress用のサービスも作ります。mysqlとは違い、外部から接続したいので、typeとしてNodePortというものを選択しています。これは、PodのポートとNodeのポートを接続するというものです。

`./manifests/wordpress-service.yaml`という名前で、以下のyamlを作成します。

```yaml
apiVersion: v1
kind: Service
metadata:
  labels:
    app: wordpress
    tier: frontend
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

`kubectl apply`コマンドで、wordpress Serviceを作成します。

```
➤ kubectl apply -f ./manifests/wordpress-service.yaml
```

サービスを見てみましょう。

```
➤ kubectl get service -l app=wordpress -l tier=frontend
NAME        TYPE       CLUSTER-IP       EXTERNAL-IP   PORT(S)        AGE
wordpress   NodePort   10.101.181.195   <none>        80:30180/TCP   48s
```

これをapplyしたら、ブラウザからアクセス可能となります。minikubeの場合、以下のコマンドでブラウザで開くことができます。Wordpressのセットアップ画面が表示されれば成功です。

```
➤ minikube service wordpress
Opening kubernetes service default/wordpress in default browser...
```

これでwordpressが起動しました！ここからは、より便利な機能（といってもk8sを使うならほぼ必須ですが）を使ってみましょう。

## アプリケーションがクラッシュした時、自動回復してほしい

さてさて無事にwordpressが起動したのですが、このままだと例えばwordpressコンテナがクラッシュした時に繋がらなくなってしまいます。Podはコンテナとストレージボリュームの集合というだけで、自分自身を管理するということをしていないためです。そこで、Deploymentというリソースを使います。Podとして定義していた箇所を、以下のようにDeploymentにします。Deploymentのspec/template/spec部分はPodのspecと同じであることがわかりますね。

`./manifests/wordpress-deployment.yaml`という名前で、以下のyamlを作成します。

```yaml
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
➤ kubectl apply -f ./manifests/wordpress-deployment.yaml
deployment "wordpress" created
```

`kubectl get deployment`コマンドで、wordpress Deploymentが作成されていることを確認しましょう。

```
➤ kubectl get deployment -l app=wordpress -l tier=frontend
NAME        DESIRED   CURRENT   UP-TO-DATE   AVAILABLE   AGE
wordpress   1         1         1            1           5m
```

DeploymentによってPodが自動作成されるので、それも確認しましょう。`wordpress`はPodとして作成したもので、`wordpress-55448464cd-mmq9w`がDeploymentによって自動作成されたものです。

```
r_takaishi@PMC02V437VHV2R:~/s/g/t/h/k8s_hands_on|k8s_hands_on⚡*?
➤ kubectl get pods -l app=wordpress -l tier=frontend
NAME                         READY     STATUS    RESTARTS   AGE
wordpress                    1/1       Running   0          18m
wordpress-55448464cd-mmq9w   1/1       Running   0          5m
```

DeploymentによってPodが作られる状態になったので、最初に作成したwordpress Podは消しておきましょう。

```
➤ kubectl delete po/wordpress
pod "wordpress" deleted
```

Deploymentは定義した状態を維持しようとします。例えば、Podを削除しても新しいPodが自動的に作成され、Podが1台動いているという状態が維持されるわけです。

```
➤ kubectl delete pod wordpress-55448464cd-mmq9w
pod "wordpress-55448464cd-mmq9w" deleted
```

このように、Deploymentが作ったPodを消してみると、すぐに新しいPodが作成されていることがわかります。

```
➤ kubectl get pods -l app=wordpress -l tier=frontend
NAME                         READY     STATUS              RESTARTS   AGE
wordpress-55448464cd-2sftv   0/1       ContainerCreating   0          1s
wordpress-55448464cd-mmq9w   0/1       Terminating         0          9m
```

k8sの大きな特徴として、この状態を維持しようとする機能があります。これにより、アプリケーションがクラッシュしたりしても自動的に復旧させることができ、サービスの運用負担軽減に繋げられるというわけです。

Deploymentを使うように変えてもPodは1つのままです。しかし、Podの数を変更することももちろん可能なので、試してみましょう。`kubectl scale`コマンドでレプリカ数を変更することで、Podの数が変わります。

```
➤ kubectl scale deployment wordpress --replicas 5
deployment "wordpress" scaled
```

`kubectl get pods`でPodの様子を見てみると、4台が追加されていることがわかりますね。

```
➤ kubectl get pods -l app=wordpress -l tier=frontend
NAME                         READY     STATUS              RESTARTS   AGE
wordpress-55448464cd-2sftv   1/1       Running             0          3m
wordpress-55448464cd-4h8zf   0/1       ContainerCreating   0          4s
wordpress-55448464cd-nk4n8   0/1       ContainerCreating   0          4s
wordpress-55448464cd-s8c7f   0/1       ContainerCreating   0          4s
wordpress-55448464cd-t4jrj   0/1       ContainerCreating   0          4s
```

これで、フロントエンドとなるwordpress Podがスケールアウトできる、wordpress環境の完成です！

## まとめ

以下の5リソースを使ってwordpress環境を作ってきました。k8sはyaml形式の定義が長くて難しい、という話も聞きますが、1つ1つ分解するとリソースを組み合わせているということがわかります。

* ConfigMap / Secret
* Pod
* Service
* Deployment
* PersistentVolume / PersistentVolumeClaim

## さらにKubernetesについて知りたい？

- https://kubernetes.io
- [O'Reilly Japan - 入門 Kubernetes](https://www.oreilly.co.jp/books/9784873118406/)