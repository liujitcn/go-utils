# stringcase

stringcase 提供 Go 标识符和常见分隔符命名转换，同时提供通用字符串拆分和判断函数。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/stringcase@latest
~~~

## 命名转换

| 函数 | 输出 |
| --- | --- |
| UpperCamelCase、PascalCase、ToPascalCase | 大驼峰，例如 UserProfile |
| LowerCamelCase、CamelCase、ToCamelCase | 小驼峰，例如 userProfile |
| ToGoPascalCase | 大驼峰，并把 API、ID、URL、UUID 等常见缩写保持为全大写 |
| ToGoCamelCase | 小驼峰，并保持常见缩写规则 |
| SnakeCase、ToSnakeCase | 小写蛇形，例如 user_profile |
| UpperSnakeCase | 大写蛇形，例如 USER_PROFILE |
| KebabCase | 小写短横线，例如 user-profile |
| UpperKebabCase | 大写短横线，例如 USER-PROFILE |

~~~go
fmt.Println(stringcase.ToGoPascalCase("base_api_id")) // BaseAPIID
fmt.Println(stringcase.ToGoCamelCase("base_api_id"))   // baseAPIID
fmt.Println(stringcase.SnakeCase("UserProfile"))       // user_profile
fmt.Println(stringcase.KebabCase("UserProfile"))       // user-profile
~~~

CommonAbbreviations 是只读语义上的缩写集合，可通过 Contains 查询某个单词是否属于当前规则。集合的底层 map 不应被修改。

## 拆分与清洗

- Split 按驼峰、数字和分隔符拆分字符串；通过 noSplit 参数可以保留不拆分的词。
- SplitByNonAlphanumeric 按非字母数字字符拆分。
- SplitAndKeepDelimiters 拆分字符串并保留分隔符。
- ReplaceNonAlphanumeric 将连续的非字母数字字符替换为指定字符串。
- IsSnakeCase 判断字符串是否符合小写蛇形规则。
- ContainsFn 使用自定义比较函数判断切片是否包含指定值。

~~~go
parts := stringcase.Split("HTTPServerV2", "HTTP")
clean := stringcase.ReplaceNonAlphanumeric("user/profile", "-")
valid := stringcase.IsSnakeCase("user_profile")
fmt.Println(parts, clean, valid)
~~~

空字符串和无效 UTF-8 的拆分结果遵循 Split 的当前实现行为；命名转换不会替调用方做业务语义判断，缩写和边界词应通过示例验证。
