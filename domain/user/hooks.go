package user

import (
	"gorm.io/gorm"
	"shopping/utils/hash"
)

// BeforeSave 保存用户之前的回调, 如果密码没有被加密，加密密码和salt
func (user User) BeforeSave(tx *gorm.DB) (err error) {
	if user.Salt == "" {
		// 为salt创建一个随机字符串
		salt := hash.CreateSalt()
		// 创建hash加密密码
		hashPassword, err := hash.HashPassword(user.Password + salt)
		if err != nil {
			return nil
		}
		user.Password = hashPassword
		user.Salt = salt
	}
	return
}
