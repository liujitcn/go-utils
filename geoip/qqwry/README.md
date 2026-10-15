# QQWry IPv4 查询

geoip/qqwry 使用包内嵌的 QQWry 数据，面向 IPv4 地址提供国家、省份、城市和运营商信息。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/geoip@latest
~~~

## 使用

~~~go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/geoip/qqwry"
)

func main() {
	client := qqwry.NewClient()
	result, err := client.Query("47.108.149.89")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s %s\n", result.Country, result.Province, result.City, result.ISP)
}
~~~

NewClient 直接读取包内嵌数据，不需要关闭。QQWry 主要支持 IPv4，IPv6 会返回 ip is not ipv4。内网地址会返回“局域网”省份和城市；数据库未命中或数据偏移无效时返回错误。

IsPrivateIP 可用于判断 IPv4 私网地址，SpiltAddress 可按 QQWry 地址格式拆分省份和城市。
