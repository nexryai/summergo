## summergo
Summaly for golang

### これは何
[misskey-dev/summaly](https://github.com/misskey-dev/summaly) の非公式なGo版

### usage
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