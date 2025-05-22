package service

import (
	"context"
	"errors"
	"strconv"
	"time"

	"github.com/zgsm/review-manager/internal/model"
	"github.com/zgsm/review-manager/internal/repository"
	"github.com/zgsm/review-manager/pkg/logger"
	// "github.com/dosun/review-manager/pkg/utils"
	"golang.org/x/crypto/bcrypt"
)

// UserService 用户服务接口
type UserService interface {
	Register(ctx context.Context, username, email, password string) (*model.User, error)
	ValidateRegisterParams(username, email, password string) error
	Login(ctx context.Context, username, password string) (*model.User, error, string)
	ValidateLoginParams(username, password string) error
	GetUserByID(ctx context.Context, id uint) (*model.User, error)
	UpdateUser(ctx context.Context, user *model.User) error
	UpdateUserInfo(ctx context.Context, userID uint, nickname, avatar string) error
	DeleteUser(ctx context.Context, id uint) error
	ValidateAndParseUserID(idStr string) (uint, error)
	ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error)
	ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error
	ValidateChangePasswordParams(oldPassword, newPassword string) error
	ValidateAndParsePageParams(pageStr, pageSizeStr string) (int, int)
}

// userService 用户服务实现
type userService struct {
	userRepo repository.UserRepository
}

// NewUserService 创建用户服务
func NewUserService() UserService {
	return &userService{
		userRepo: repository.NewUserRepository(),
	}
}

// ValidateRegisterParams 验证注册参数
func (s *userService) ValidateRegisterParams(username, email, password string) error {
	if username == "" || email == "" || password == "" {
		return errors.New("用户名、邮箱和密码不能为空")
	}
	return nil
}

// Register 用户注册
func (s *userService) Register(ctx context.Context, username, email, password string) (*model.User, error) {
	// 检查用户名是否已存在
	existUser, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		logger.Error("检查用户名失败", "error", err)
		return nil, err
	}
	if existUser != nil {
		return nil, errors.New("用户名已存在")
	}

	// 检查邮箱是否已存在
	existUser, err = s.userRepo.GetByEmail(ctx, email)
	if err != nil {
		logger.Error("检查邮箱失败", "error", err)
		return nil, err
	}
	if existUser != nil {
		return nil, errors.New("邮箱已存在")
	}

	// 加密密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	if err != nil {
		logger.Error("密码加密失败", "error", err)
		return nil, err
	}

	// 创建用户
	user := &model.User{
		Username:  username,
		Email:     email,
		Password:  string(hashedPassword),
		Nickname:  username,
		Role:      "user",
		Status:    1,
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}

	// 保存用户
	if err := s.userRepo.Create(ctx, user); err != nil {
		logger.Error("创建用户失败", "error", err)
		return nil, err
	}

	return user, nil
}

// ValidateLoginParams 验证登录参数
func (s *userService) ValidateLoginParams(username, password string) error {
	if username == "" || password == "" {
		return errors.New("用户名和密码不能为空")
	}
	return nil
}

// Login 用户登录
func (s *userService) Login(ctx context.Context, username, password string) (*model.User, error, string) {
	// 根据用户名获取用户
	user, err := s.userRepo.GetByUsername(ctx, username)
	if err != nil {
		logger.Error("获取用户失败", "error", err)
		return nil, err, ""
	}
	if user == nil {
		return nil, errors.New("用户不存在"), ""
	}

	// 检查用户状态
	if user.Status != 1 {
		return nil, errors.New("用户已被禁用"), ""
	}

	// 验证密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(password)); err != nil {
		return nil, errors.New("密码错误"), ""
	}

	// 更新最后登录时间
	now := time.Now()
	user.LastLogin = &now
	if err := s.userRepo.Update(ctx, user); err != nil {
		logger.Error("更新最后登录时间失败", "error", err)
		// 不返回错误，继续登录流程
	}

	// 生成token（实际项目中应使用JWT等认证方式）
	token := "sample_token_" + user.Username

	return user, nil, token
}

// GetUserByID 根据ID获取用户
func (s *userService) GetUserByID(ctx context.Context, id uint) (*model.User, error) {
	return s.userRepo.GetByID(ctx, id)
}

// UpdateUser 更新用户
func (s *userService) UpdateUser(ctx context.Context, user *model.User) error {
	return s.userRepo.Update(ctx, user)
}

// DeleteUser 删除用户
func (s *userService) DeleteUser(ctx context.Context, id uint) error {
	return s.userRepo.Delete(ctx, id)
}

// ListUsers 获取用户列表
func (s *userService) ListUsers(ctx context.Context, page, pageSize int) ([]*model.User, int64, error) {
	return s.userRepo.List(ctx, page, pageSize)
}

// ValidateChangePasswordParams 验证修改密码参数
func (s *userService) ValidateChangePasswordParams(oldPassword, newPassword string) error {
	if oldPassword == "" || newPassword == "" {
		return errors.New("旧密码和新密码不能为空")
	}
	return nil
}

// ChangePassword 修改密码
func (s *userService) ChangePassword(ctx context.Context, userID uint, oldPassword, newPassword string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 验证旧密码
	if err := bcrypt.CompareHashAndPassword([]byte(user.Password), []byte(oldPassword)); err != nil {
		return errors.New("旧密码错误")
	}

	// 加密新密码
	hashedPassword, err := bcrypt.GenerateFromPassword([]byte(newPassword), bcrypt.DefaultCost)
	if err != nil {
		return err
	}

	// 更新密码
	user.Password = string(hashedPassword)
	user.UpdatedAt = time.Now()
	return s.userRepo.Update(ctx, user)
}

// UpdateUserInfo 更新用户信息
func (s *userService) UpdateUserInfo(ctx context.Context, userID uint, nickname, avatar string) error {
	// 获取用户
	user, err := s.userRepo.GetByID(ctx, userID)
	if err != nil {
		return err
	}
	if user == nil {
		return errors.New("用户不存在")
	}

	// 更新用户信息
	if nickname != "" {
		user.Nickname = nickname
	}
	if avatar != "" {
		user.Avatar = avatar
	}
	user.UpdatedAt = time.Now()

	return s.userRepo.Update(ctx, user)
}

// ValidateAndParsePageParams 验证并解析分页参数
func (s *userService) ValidateAndParsePageParams(pageStr, pageSizeStr string) (int, int) {
	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	pageSize, err := strconv.Atoi(pageSizeStr)
	if err != nil || pageSize < 1 || pageSize > 100 {
		pageSize = 10
	}

	return page, pageSize
}

// ValidateAndParseUserID 验证并解析用户ID
func (s *userService) ValidateAndParseUserID(idStr string) (uint, error) {
	id, err := strconv.Atoi(idStr)
	if err != nil {
		return 0, errors.New("无效的用户ID")
	}
	return uint(id), nil
}
