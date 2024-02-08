## summergo
![Go](https://img.shields.io/badge/go-%2300ADD8.svg?style=for-the-badge&logo=go&logoColor=white) 
[![No GPL](https://img.shields.io/badge/real_permissive-No_GPL-99CC33?style=for-the-badge&logo=opensourceinitiative&logoColor=white)](https://log.sda1.net/blog/no-gpl-badge/)<br>

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
