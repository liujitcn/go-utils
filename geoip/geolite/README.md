# GeoLite2 城市查询

geoip/geolite 使用包内嵌的 GeoLite2 City 数据，支持 IPv4 和 IPv6。返回国家、省份、城市和 IP 字段；运营商字段通常为空。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/geoip@latest
~~~

## 使用

~~~go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/geoip/geolite"
)

func main() {
	client, err := geolite.NewClient()
	if err != nil {
		panic(err)
	}
	defer client.Close()

	result, err := client.Query("8.8.8.8")
	if err != nil {
		panic(err)
	}
	fmt.Printf("%s %s %s\n", result.Country, result.Province, result.City)
}
~~~

NewClient 从内嵌数据库创建读取器；查询结束后调用 Close 释放资源。数据库默认使用 zh-CN 名称，可通过 SetLanguage 设置 GeoLite2 中存在的其他语言代码。

私网地址不会查询数据库，而是将国家、省份和城市均返回为“局域网”。无效 IP 会返回错误。
