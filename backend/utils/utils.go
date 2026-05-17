package utils

import (
	"time"

	"github.com/golang-jwt/jwt"
	"golang.org/x/crypto/bcrypt"
)

// 这里使用bcrypt算法对密码进行哈希处理，bcrypt是一个适合密码存储的哈希函数，具有内置的盐和可调节的工作因子，可以有效抵抗暴力破解和字典攻击。
// HashPassword函数接受一个字符串类型的密码作为输入，返回一个哈希后的密码字符串和一个错误对象。
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}

func GenerateJWT(username string) (string, error) {
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, jwt.MapClaims{
		"username": username,
		"exp":      time.Now().Add(time.Hour * 24).Unix(), //设置过期时间为24小时
	})
	Token, err := token.SignedString([]byte("secret"))
	//使用一个密钥字符串来签名生成的JWT令牌，这里使用了一个简单的字符串"secret"，在实际应用中应该使用更复杂和安全的密钥。
	return "Bearer" + Token, err
}

// CheckPassword函数接受一个明文密码和一个哈希密码作为输入，使用bcrypt的CompareHashAndPassword函数比较两者是否匹配，如果匹配返回true，否则返回false。
// 这个函数可以用于用户登录时验证输入的密码是否正确。
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}
