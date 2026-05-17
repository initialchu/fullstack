package utils

import (
	"errors"
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
	return "Bearer " + Token, err
}

// CheckPassword函数接受一个明文密码和一个哈希密码作为输入，使用bcrypt的CompareHashAndPassword函数比较两者是否匹配，如果匹配返回true，否则返回false。
// 这个函数可以用于用户登录时验证输入的密码是否正确。
func CheckPassword(password string, hash string) bool {
	err := bcrypt.CompareHashAndPassword([]byte(hash), []byte(password))
	return err == nil

}

// ParseJWT函数接受一个JWT令牌字符串作为输入，解析并验证该令牌的有效性，如果令牌有效，返回其中包含的用户名，否则返回一个错误对象。
// 返回值是一个字符串类型的用户名和一个错误对象，如果解析成功，错误对象为nil；如果解析失败，用户名为空字符串，错误对象包含具体的错误信息。
func ParseJWT(tokenString string) (string, error) {
	if len(tokenString) > 7 && tokenString[:7] == "Bearer " {
		//如果令牌字符串以"Bearer "开头，去掉这个前缀，得到实际的JWT令牌字符串
		tokenString = tokenString[7:]
	}
	token, err := jwt.Parse(tokenString, func(token *jwt.Token) (interface{}, error) {
		if _, ok := token.Method.(*jwt.SigningMethodHMAC); !ok {
			return nil, errors.New("unexpected signing method")
		}
		//返回用于验证JWT令牌的密钥，这里使用了一个简单的字符串"secret"
		return []byte("secret"), nil

	})
	if err != nil {
		return "", err
	}
	if claims, ok := token.Claims.(jwt.MapClaims); ok && token.Valid {
		//如果令牌的声明部分可以转换为jwt.MapClaims类型，并且令牌有效，返回其中包含的用户名
		username, ok := claims["username"].(string)
		if !ok {
			return "", errors.New("username claim is not a string")

		}
		return username, nil
	}
	return "", err
}
