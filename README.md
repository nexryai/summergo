## summergo
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) 
[![No GPL](https://img.shields.io/badge/real_permissive-No_GPL-99CC33?style=for-the-badge&logo=opensourceinitiative&logoColor=white)](https://log.sda1.net/blog/no-gpl-badge/)<br>

Summaly for golang

### これは何
[misskey-dev/summaly](https://github.com/misskey-dev/summaly) の非公式なGo版

### usage (deploy)
 - [Deploy with Docker](https://github.com/nexryai/summaly-go)
   * Google Cloud RunやFly.io、その他クラウド/オンプレミスサーバーにデプロイする場合におすすめです。
   * サービス側の問題でコールドスタートにタイムアウトするほどの時間がかかるため、**Azure Container Appsへのデプロイは推奨しません。**
     * デプロイする場合は最小インスタンス数を0にする必要があります。
     * ↑当然料金が跳ね上がるので自己責任でお願いします。
   * 自力でビルドしてコンテナを使わずに動かすことも可能です
 - [Deploy to AWS Lambda](https://github.com/nexryai/summaly-lambda)
   * Lambdaにデプロイしたい場合にお使いください。

### usage (as Go module)
`github.com/nexryai/summergo`をimportすることで、Goのライブラリとして使用できます。  
`model.go`に完全な取得できるデータの構造体があります

```go
summaly, err := summergo.Summarize("https://www.youtube.com/watch?v=U1yqKWN80EM")
if err != nil {
    panic(err)
}

fmt.Println(summaly.Title)
fmt.Println(summaly.Description)
fmt.Println(summaly.Player.Url)
```

### Security
脆弱性を発見した場合、GitHubのセキュリティアドバイザリ機能を使用して報告してください。  
SSRF攻撃の対策は基本的なものを行なっていますが、完全ではないため**プライベートネットワークや内部サービスにアクセス可能な環境ではホストしないでください。**  
