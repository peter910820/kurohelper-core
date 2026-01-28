# kurohelper-core

KuroHelper 底層核心專案

- 核心專案有一些直接存取 env 的行為，在使用前請先確保啟動端有設置特定 env，避免使用上受影響(可參考[kurohelper/.env.example](https://github.com/kuro-helper/kurohelper/blob/main/.env.example))

## 安裝

使用 Go Modules：

> [!IMPORTANT]
> 此專案被golang proxy快取的最終版本為**v3.1.0**
> 之後的版本採用手動建置的方式，請自行clone以及replace

- v3(僅到**v3.1.0**)

```bash
go get github.com/kuro-helper/kurohelper-core/v3
```

- v2

```bash
go get github.com/kuro-helper/core/v2
```

- v1

```bash
go get github.com/peter910820/kurohelper-core
```
