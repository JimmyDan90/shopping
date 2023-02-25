package user

import "shopping/utils/hash"

// Service 用户结构体
type Service struct {
	r Repository
}

// NewUserService 实例化
func NewUserService(r Repository) *Service {
	r.Migration()
	r.InsertSampleData()
	return &Service{
		r: r,
	}
}

// Create 创建用户
func (s Service) Create(user *User) error {
	if user.Password == user.Password2 {
		return ErrMismatchedPasswords
	}
	// 用户名存在
	_, err := s.r.GetByName(user.Username)
	if err == nil {
		return ErrUserExistWithName
	}
	// 无效用户名
	if ValidateUserName(user.Username) {
		return ErrInvalidUsername
	}
	// 无效密码
	if ValidatePassword(user.Password) {
		return ErrInvalidPassword
	}
	// 创建用户
	err = s.r.Create(user)
	return err
}

func (s *Service) UpdateUser(user *User) error {
	return s.r.Update(user)
}

// GetUser 查询用户
func (s *Service) GetUser(username, password string) (User, error) {
	user, err := s.r.GetByName(username)
	if err != nil {
		return User{}, ErrUserNotFound
	}
	match := hash.CheckPasswordHash(password+user.Salt, user.Password)
	if !match {
		return User{}, ErrUserNotFound
	}
	return user, nil
}
