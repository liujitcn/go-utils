package crypto

import (
	"fmt"
	"strings"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestHashPassword(t *testing.T) {
	text := "admin123"
	hash, _ := HashPassword(text)
	fmt.Println(hash)
}

func TestCheckPasswordHash(t *testing.T) {
	text := "super"
	hash3 := "$2a$10$Q3DL31uqv9nJ8/.tIQRbRe8519czV4AMuBzQFyuch4NUaOt3ZfG0i"
	bMatched := CheckPasswordHash(text, hash3)
	assert.True(t, bMatched)

	bMatched = CheckPasswordHash(text, hash3)
	assert.True(t, bMatched)
}

func TestJwtToken(t *testing.T) {
	const bearerWord string = "Bearer"
	token := "Bearer eyJhbGciOiJIUzI1NiIsInR5cCI6IkpXVCJ9.eyJ1c2VyX2lkIjowfQ.XgcKAAjHbA6o4sxxbEaMi05ingWvKdCNnyW9wowbJvs"
	auths := strings.SplitN(token, " ", 2)
	assert.Equal(t, len(auths), 2)
	assert.Equal(t, strings.EqualFold(auths[0], bearerWord), true, "JWT token is missing")
}
