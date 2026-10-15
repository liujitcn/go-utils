# go-utils

`go-utils` 是面向 Go 业务服务的通用工具库，按功能拆分为多个可独立引入的包和子模块。项目要求 Go 1.27。

## 安装

根模块和独立子模块需要按实际使用的导入路径安装：

```bash
go get github.com/liujitcn/go-utils@latest
go get github.com/liujitcn/go-utils/crypto@latest
go get github.com/liujitcn/go-utils/jwt@latest
```

本地开发可以在仓库根目录执行：

```bash
make init  # 安装格式化相关命令
make fmt   # 整理 Go 导入和代码格式
```

## 功能索引

### 根模块

根模块不包含业务状态，提供以下通用包：

| 包 | 功能 |
| --- | --- |
| `byte` | `int` 与大端 `[]byte` 互转、ASCII 字节大小写转换 |
| `id` | Snowflake、UUID v4/v7、短 UUID、XID 生成及 Snowflake 时间提取 |
| `ip` | 从代理请求头或 `RemoteAddr` 提取客户端 IP |
| `io` | 文件读取、目录遍历、路径匹配、用户目录展开、文件状态和权限判断 |
| `map` | 泛型 map 合并、遍历、过滤和删除键 |
| `mapper` | DTO 与实体双向复制，以及时间、JSON、枚举类型转换器 |
| `set` | 泛型集合、集合运算、JSON 编解码和线程安全实现 |
| `slice` | 泛型切片过滤、映射、归约、查找、集合运算和展平 |
| `string` | JSON/逗号分隔数组、随机数字、手机号脱敏和金额元分转换 |
| `stringcase` | 驼峰、蛇形、短横线、Go 缩写保留和字符串拆分 |
| `time` | 时间格式化、常用日期区间、时间差和 protobuf 时间转换 |
| `tls` | 从文件或 PEM 内容创建客户端/服务端 TLS 配置 |
| `trans` | 基础值与指针、切片与指针切片互转 |

### 独立子模块

| 模块 | 功能 | 说明 |
| --- | --- | --- |
| `crypto` | 密码哈希、HMAC、AES、SM4、RSA、ECDSA/ECDH 和 OpenSSL 文件加解密 | 见 [crypto/README.md](crypto/README.md) |
| `geoip` | GeoLite2、QQWry、IP2Region 的 IP 归属地查询 | 见 [geoip/README.md](geoip/README.md) |
| `http` | 普通 HTTP、响应体限制、流式响应和 SSE | 见 [http/README.md](http/README.md) |
| `jwt` | JWT 生成、签名验证、刷新、声明读取和请求头提取 | 见 [jwt/README.md](jwt/README.md) |
| `translator` | 统一翻译接口和 Alibaba、Baidu、Google、Volc Provider | 见 [translator/README.md](translator/README.md) |

## 快速示例

### ID、时间和字符串

```go
package main

import (
	"fmt"

	strutil "github.com/liujitcn/go-utils/string"
	"github.com/liujitcn/go-utils/id"
	timeutil "github.com/liujitcn/go-utils/time"
)

func main() {
	fmt.Println(id.NewGUIDv7())
	fmt.Println(timeutil.GetTodayRangeTimeString())
	fmt.Println(strutil.ConvertYuanStringToFen("12.34")) // 1234
}
```

### 切片、map 和 set

```go
package main

import (
	"fmt"

	maputil "github.com/liujitcn/go-utils/map"
	"github.com/liujitcn/go-utils/set"
	"github.com/liujitcn/go-utils/slice"
)

func main() {
	values := slice.Filter([]int{1, 2, 3, 4}, func(value, _ int, _ []int) bool {
		return value%2 == 0
	})
	merged := maputil.Merge(map[string]int{"a": 1}, map[string]int{"a": 2, "b": 3})
	items := set.New("a", "b", "a")
	fmt.Println(values, merged, set.Sorted(items))
}
```

`map.Drop` 会直接修改传入的 map；`set.New` 创建线程安全集合，需要无锁访问时使用 `set.NewThreadUnsafeSet`。

### DTO 与实体映射

```go
package main

import (
	"fmt"
	"time"

	"github.com/liujitcn/go-utils/mapper"
)

type UserDTO struct { CreatedAt string }
type UserEntity struct { CreatedAt time.Time }

func main() {
	m := mapper.NewCopierMapper[UserDTO, UserEntity]()
	dto := m.ToDTO(&UserEntity{CreatedAt: time.Now()})
	fmt.Println(dto.CreatedAt)
}
```

映射器默认注册 `time.Time` 与字符串、常用基础切片与 JSON 字符串的转换；未知类型可以通过 `AppendConverter` 或 `AppendConverters` 追加规则。复制配置或类型不匹配会触发 panic，应在开发阶段修正映射配置。

### HTTP 客户端

```go
package main

import (
	"fmt"

	httputil "github.com/liujitcn/go-utils/http"
)

func main() {
	client := httputil.NewClient(
		httputil.WithBaseURL("https://example.com"),
		httputil.WithDefaultHeader("X-App", "demo"),
		httputil.WithMaxResponseBodyBytes(8<<20),
	)
	var result struct { Message string `json:"message"` }
	if err := client.Get("/ping", &result); err != nil {
		panic(err)
	}
	fmt.Println(result.Message)
}
```

普通请求使用 `Get`、`Post` 等方法直接解码 JSON；需要读取原始响应时使用 `Do`，长连接使用 `DoStream` 或 `DoSSE`。普通请求默认 30 秒超时并记录最终 URL，流式请求会忽略客户端总超时，长连接仍应通过 `WithContext` 取消。

### JWT

```go
package main

import (
	"fmt"

	jwtutil "github.com/liujitcn/go-utils/jwt"
	"github.com/golang-jwt/jwt/v5"
)

func main() {
	secret := []byte("replace-with-a-secret")
	token, err := jwtutil.GenerateJWT(
		jwt.MapClaims{"sub": "user-1"},
		secret,
		jwt.SigningMethodHS256,
	)
	if err != nil { panic(err) }
	claims, err := jwtutil.VerifyJWT(token, secret)
	if err != nil { panic(err) }
	fmt.Println(claims["sub"])
}
```

`ParseJWTPayload`、`GetJWTHeader`、`ValidateJWTAlgorithm` 等解析辅助函数不验证签名；处理不可信 token 时应先使用 `VerifyJWT`，再读取声明。

### IP 归属地

```go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/geoip/qqwry"
)

func main() {
	client := qqwry.NewClient()
	result, err := client.Query("47.108.149.89")
	if err != nil { panic(err) }
	fmt.Printf("%s %s %s %s\n", result.Country, result.Province, result.City, result.ISP)
}
```

三种查询实现都使用包内嵌数据：`geolite` 支持 IPv4/IPv6 城市信息，`qqwry` 主要支持 IPv4，`ip2region` 支持 IPv4/IPv6 和可选缓存策略。`geolite.Client` 和 `ip2region.Client` 使用完毕后需要调用 `Close`。

## 重要行为

- `time` 默认按 `Asia/Shanghai` 解析字符串；可用 `RefreshDefaultTimeLocation` 切换时区。常用字符串格式为 `2006-01-02`、`2006-01-02 15:04:05`。
- `ip.GetClientRealIP` 按 `X-Forwarded-For`、`X-Real-IP`、`RemoteAddr` 顺序读取地址。只有在代理链可信时才应直接信任这些请求头。
- `io.ReadFile`、目录和权限判断函数以返回空值或 `false` 表示失败；需要区分具体原因时请使用标准库 `os` API。
- `tls.LoadServerTlsConfig*` 传入 CA 后会要求并校验客户端证书，形成双向 TLS；`insecureSkipVerify` 仅用于服务端配置中的 `tls.Config` 字段，应谨慎使用。
- AES-CBC、SM4-CBC 保留用于兼容既有数据；新数据优先使用带认证的 AES-GCM 或 SM4-GCM。

## 测试

根目录的 `go.work` 将所有本地模块纳入同一工作区。可在仓库根目录运行：

```bash
go test ./...
```

也可以分别测试独立子模块：

```bash
for module in crypto geoip http jwt translator translator/alibaba translator/baidu translator/google translator/volc; do
  (cd "$module" && go test ./...)
done
```

## 目录说明

```text
.
├── byte id ip io map mapper set slice string stringcase time tls trans  # 根模块工具包
├── crypto geoip http jwt translator                                    # 独立子模块
├── go.work                                                             # 本地多模块工作区
└── scripts/tag_release.py                                              # 版本标签辅助脚本
```
