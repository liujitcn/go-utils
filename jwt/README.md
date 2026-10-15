# jwt

jwt 是独立 Go 子模块，基于 HMAC JWT 提供生成、验证、刷新和声明读取辅助函数。

## 安装

~~~bash
go get github.com/liujitcn/go-utils/jwt@latest
~~~

## 生成与验证

~~~go
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
	if err != nil {
		panic(err)
	}

	claims, err := jwtutil.VerifyJWT(token, secret)
	if err != nil {
		panic(err)
	}
	fmt.Println(claims["sub"])
}
~~~

主要生成函数：

- GenerateJWT：从 jwt.MapClaims 生成 token。
- GenerateGenericJWT：把 JSON 对象型泛型值转换为 claims 后生成 token。
- GenerateJWTWithHeader：生成 token 并追加自定义 Header。
- GenerateShortLivedJWT：写入当前时间加 duration 的 exp 后生成 token。
- RefreshJWT：验证 HMAC 签名后更新 exp，并使用 HS256 重新签名。

GenerateJWT、GenerateGenericJWT、GenerateJWTWithHeader 和 GenerateShortLivedJWT 都要求非空密钥、非空 claims 和非 nil 签名方法。

## 解析与校验

- VerifyJWT 验证签名，并只接受 HMAC 签名方法。
- ParseJWTPayload、ParseJWTClaimsToStruct、GetJWTClaims、GetJWTHeader 只解析 token，不验证签名。
- ValidateJWTAlgorithm 检查 Header 中的 alg 字段，但不验证签名。
- ValidateJWTAudience、ValidateJWTIssuer、IsJWTExpired、GetJWTIssuedAt 读取对应声明；这些辅助函数本身不替代签名验证。
- ExtractJWTFromRequest 从 Authorization: Bearer <token> 请求头提取 token。

不可信 token 的处理顺序应是先 VerifyJWT，再读取 claims 或执行业务授权判断。

## 令牌辅助

- NewRefreshToken 使用 32 字节随机数生成 Base64URL 刷新令牌。
- NewJWTId 使用 UUID v7 生成无连字符的 jti。
- ToMapClaims 将可编码为 JSON 对象的泛型值转换为 jwt.MapClaims；数组、字符串和数字等非对象值不能直接作为 claims。

## 注意事项

- secretKey 由业务安全保存，示例中的密钥仅用于演示。
- RefreshJWT 只接受 HS256、HS384、HS512，并会将刷新后的算法固定为 HS256。
- ParseJWTPayload 和声明检查函数不提供安全边界，不能单独用于身份认证或授权。

## 测试

~~~bash
go test ./...
~~~
