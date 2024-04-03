package domain

import (
	"context"
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"edu/internal/domain/enum"
	"encoding/base64"
	"fmt"

	"github.com/go-kratos/kratos/v2/log"
	uuid "github.com/satori/go.uuid"
)

type UserRepo interface {
	BaseRepo[User]
}

type User struct {
	Id          int32
	Username    string
	Password    string
	Salt        string
	Mobile      string
	Avatar      string
	Description string
	RoleId      int32
	RoleName    string

	Status       int32
	UpdatedAt    string
	Version      int32
	TenantId     int32
	IsSuperAdmin bool
	IsTenant     bool
}

// UserService is a Greeter usecase.
type UserService struct {
	BaseService[User]
	log *log.Helper
}

func (u *UserService) Create(ctx context.Context, user *User) (*User, error) {
	u4 := uuid.NewV4()
	salt := u4.String()
	h := hmac.New(sha256.New, []byte(salt))
	h.Write([]byte(fmt.Sprintf("%x", md5.Sum([]byte(user.Password)))))

	user.Salt = salt
	user.Password = base64.URLEncoding.EncodeToString(h.Sum(nil))
	nUser, err := u.repo.Save(ctx, user)

	if err != nil {
		return nil, err
	}
	return nUser, err
}

func (us *UserService) Authenticate(ctx context.Context, mobile, password string) (*User, error) {

	users, err := us.ListByMap(ctx, map[string]interface{}{"mobile": mobile, "status": enum.EnableStatusEnabled})
	if err != nil {
		return nil, err
	}
	if len(users) != 1 {
		return nil, fmt.Errorf("用户不存在,请联系管理员")
	}

	var user = users[0]
	h := hmac.New(sha256.New, []byte(user.Salt))
	h.Write([]byte(password))
	if user.Password != base64.URLEncoding.EncodeToString(h.Sum(nil)) {
		return nil, fmt.Errorf("用户名或密码错误")
	}

	return users[0], nil
}

func (us *UserService) ChangePassword(ctx context.Context, id int32, oldPassword, newPassword string) error {

	user, err := us.FindByID(ctx, id)
	if err != nil {
		return err
	}

	// 验证老密码
	h := hmac.New(sha256.New, []byte(user.Salt))
	h.Write([]byte(oldPassword))
	if user.Password != base64.URLEncoding.EncodeToString(h.Sum(nil)) {
		return fmt.Errorf("原密码错误")
	}

	// 新密码
	hn := hmac.New(sha256.New, []byte(user.Salt))
	hn.Write([]byte(newPassword))
	newEncryptPassword := base64.URLEncoding.EncodeToString(hn.Sum(nil))

	err = us.UpdateConcurrency(ctx, &User{Id: user.Id, Password: newEncryptPassword, Version: user.Version})
	if err != nil {
		return err
	}

	return nil
}

// NewUserService new a User usecase.
func NewUserService(repo UserRepo, logger log.Logger) *UserService {
	return &UserService{
		BaseService: BaseService[User]{repo: repo}, log: log.NewHelper(logger)}
}
