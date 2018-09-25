# Learn Rust

プログラミングRust

## 第1章 なぜRustなのか？

* システム向けプログラミングは難しい
  * メモリ管理が難しい
  * マルチスレッドのコードを書くことが難しい

## 第2章 Rustツアー

* 今回は`brew install rust`でインストールした
* Result型便利そう
* 


並列化しない場合：

``
$ time ./target/debug/mandelbrot mandel.png 1000x750 -1.20,0.35 -1,0.20
Hello, world!
        8.80 real         8.78 user         0.00 sys
```

並列化した場合：

```
$ time ./target/debug/mandelbrot mandel.png 1000x750 -1.20,0.35 -1,0.20
Hello, world!
        3.15 real         8.99 user         0.01 sys
```

## 第3章 基本的な型

- Rustの型の目的（安全性/効率性/簡潔性）
- 暗黙の型変換がない


