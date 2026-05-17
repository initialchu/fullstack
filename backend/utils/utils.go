package utils

import "golang.org/x/crypto/bcrypt"

//这里使用bcrypt算法对密码进行哈希处理，bcrypt是一个适合密码存储的哈希函数，具有内置的盐和可调节的工作因子，可以有效抵抗暴力破解和字典攻击。
//HashPassword函数接受一个字符串类型的密码作为输入，返回一个哈希后的密码字符串和一个错误对象。
func HashPassword(pwd string) (string, error) {
	hash, err := bcrypt.GenerateFromPassword([]byte(pwd), 12)
	return string(hash), err
}
