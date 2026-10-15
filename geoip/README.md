# geoip

geoip 是独立 Go 子模块，统一定义 IP 归属地查询结果，并提供三种内置数据实现。查询结果类型为：

~~~go
type Result struct {
	IP       string
	Country  string
	Province string
	City     string
	ISP      string
}
~~~

实现都满足：

~~~go
type GeoIP interface {
	Query(queryIP string) (Result, error)
}
~~~

## 安装

~~~bash
go get github.com/liujitcn/go-utils/geoip@latest
~~~

## 实现选择

| 包 | 数据和地址范围 | 构造与释放 |
| --- | --- | --- |
| geoip/geolite | 内嵌 GeoLite2 City，支持 IPv4/IPv6 和城市级国家、省份、城市信息 | NewClient() (*Client, error)；使用后 Close() |
| geoip/qqwry | 内嵌 QQWry，主要支持 IPv4，并尽量返回省份、城市和运营商 | NewClient() *Client；不需要关闭 |
| geoip/ip2region | 内嵌 IP2Region v4/v6 数据，支持 IPv4/IPv6 和缓存策略 | NewClient() (*Client, error)；使用后 Close() |

所有数据文件都通过 Go embed 编译进对应包，运行时不需要额外下载数据库。数据覆盖范围和准确性取决于包内数据库版本。

## 使用示例

~~~go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/geoip/ip2region"
)

func main() {
	client, err := ip2region.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	result, err := client.Query("2001:db8::1")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", result.Country, result.Province, result.City, result.ISP)
}
~~~

需要自定义 IP2Region 数据或并发搜索器数量时，可使用 NoCache、VIndexCache、BufferCache，分别传给 NewV4Config/NewV6Config，再通过 NewIp2Region 组合 IPv4 和 IPv6 配置；某一版本配置传入 nil 可关闭该版本查询。

## 选择建议

- 需要 IPv6 或城市级信息时优先使用 geolite 或 ip2region。
- 需要轻量 IPv4 查询并兼容 QQWry 数据时使用 qqwry。
- geolite 的私网地址会返回“局域网”；qqwry 会将国家设置为“中国”并把省份和城市设置为“局域网”；ip2region 直接返回数据库查询结果。
- 查询失败会返回错误；调用方应处理无效 IP、数据库未命中和资源关闭错误。

## 数据说明

- GeoLite2：<https://dev.maxmind.com/geoip/geolite2-free-geolocation-data>
- 纯真 QQWry：<https://update.cz88.net/geo-public>
- IP2Region：<https://ip2region.net/>

上述地址用于了解数据格式和更新信息，运行时使用的是本模块内嵌的数据文件。
