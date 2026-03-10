# go-utils

`go-utils` 是一个通用 Go 工具库，覆盖常用的 ID 生成、时间处理、集合操作、字符串转换、文件路径与 TLS/JWT/密码学能力。

## 安装

```bash
go get github.com/liujitcn/go-utils@latest
```

## 模块结构

- `github.com/liujitcn/go-utils`：根包 `utils`（统计时间与金额转换）
- `github.com/liujitcn/go-utils/byte`：`int` 与 `[]byte` 转换、ASCII 大小写字节转换
- `github.com/liujitcn/go-utils/crypto`：密码学工具（独立子模块）
- `github.com/liujitcn/go-utils/id`：Snowflake、UUIDv4/v7、KSUID、XID、Mongo ObjectID
- `github.com/liujitcn/go-utils/io`：文件读取、路径匹配、文件属性判断
- `github.com/liujitcn/go-utils/jwt`：JWT 生成/解析/校验（独立子模块）
- `github.com/liujitcn/go-utils/map`：泛型 map 工具
- `github.com/liujitcn/go-utils/slice`：泛型 slice 工具
- `github.com/liujitcn/go-utils/string`：字符串与 JSON 数组转换、脱敏、随机数字串
- `github.com/liujitcn/go-utils/stringcase`：大小驼峰、蛇形、短横线等命名转换
- `github.com/liujitcn/go-utils/time`：时间格式化、区间、差值、protobuf 时间转换
- `github.com/liujitcn/go-utils/tls`：TLS 配置加载（文件/内存）
- `github.com/liujitcn/go-utils/trans`：值与指针/切片互转、Map Key/Value 提取

## 快速示例

### 根包 `utils`

```go
package main

import (
	"fmt"

	utils "github.com/liujitcn/go-utils"
)

func main() {
	fmt.Println(utils.CalcGrowthRate(100, 120)) // 2000 => 20.00%
	fmt.Println(utils.ConvertYuanToFen("12.34"))
}
```

### `id`

```go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/id"
)

func main() {
	sf, _ := id.NewSnowflake()
	sid := sf.NextVal()
	fmt.Println(sid, id.GetGenTime(sid))

	fmt.Println(id.NewGUIDv7(true))
	fmt.Println(id.NewXID())
}
```

### `time`

```go
package main

import (
	"fmt"

	timeutil "github.com/liujitcn/go-utils/time"
)

func main() {
	start, end := timeutil.GetTodayRangeTimeString()
	fmt.Println(start, end)
}
```

### `slice`

```go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/slice"
)

func main() {
	arr := []int{1, 2, 3, 4}
	evens := slice.Filter(arr, func(v, _ int, _ []int) bool { return v%2 == 0 })
	fmt.Println(evens) // [2 4]
}
```

### `map`

```go
package main

import (
	"fmt"

	_map "github.com/liujitcn/go-utils/map"
)

func main() {
	m := map[string]int{"a": 1, "b": 2}
	fmt.Println(_map.Keys(m))
}
```

### `string`

```go
package main

import (
	"fmt"

	strutil "github.com/liujitcn/go-utils/string"
)

func main() {
	fmt.Println(strutil.ConvertStringToInt64Array("1,2, 3"))
	fmt.Println(strutil.DesensitizePhone("13812345678"))
}
```

### `trans`

```go
package main

import (
	"fmt"

	"github.com/liujitcn/go-utils/trans"
)

func main() {
	name := trans.String("alice")
	fmt.Println(trans.StringValue(name))
	fmt.Println(trans.MapKeys(map[string]int{"a": 1, "b": 2}))
}
```

## 子模块说明

- `crypto`：见 [crypto/README.md](crypto/README.md)
- `stringcase`：见 [stringcase/README.md](stringcase/README.md)
- ECDH 使用建议：先通过 `Encrypt` 交换公钥，再用 `DeriveSharedSecret` 生成共享密钥；`Verify` 的第一个参数需传共享密钥的 Base64 字符串。

## 测试

在仓库根目录：

```bash
go test ./...
```

进入独立子模块执行：

```bash
cd crypto && go test ./...
cd jwt && go test ./...
```

## 打 Tag

```bash
make tag     # 根模块：仅根据远程仓库更新状态决定是否打并推送远程 tag（不提交代码）
make sub-tag # 多模块：递归检查 go.mod 目录，仅根据远程仓库更新状态为模块打并推送远程 tag（不提交代码）
```

说明：上述命令通过 `python3 scripts/tag_release.py` 执行统一的版本计算与远程更新检测逻辑。
