package userService

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"jiyu/config"
	"jiyu/global"
	"jiyu/global/redis"
	"jiyu/model/contextModel"
	"jiyu/model/globalModel"
	"jiyu/model/positionModel"
	"jiyu/model/userModel"
	"jiyu/repo/globalRepo"
	"jiyu/repo/userRepo"
	"jiyu/service/adminUserService"
	"jiyu/service/globalService"
	"jiyu/service/memberService"
	"jiyu/util"
	"jiyu/util/dysms"
	"jiyu/util/email"
	"jiyu/util/gaode"

	// "jiyu/util/ishumei"
	"log"
	"strconv"
	"time"

	"github.com/samber/lo"
	"github.com/spf13/cast"
	"gorm.io/gorm"
)

func Login(nickname string, ip string) (jwtToken string, err error) {
	var (
		uid  uint
		user *userModel.User
	)

	if !ExistByNickname(nickname) {
		return "", errors.New("用户不存在")

	} else {
		user, err = userRepo.FindByUserName(nickname)
		if err != nil {
			log.Println("userService.Login user登陆失败, error: ", err)
			return "", errors.New("用户登陆失败")
		}

		uid = user.ID
	}

	// 生成jwt
	jwtToken = util.CreateJWT(uid, nickname)
	return
}

func LoginSMSCode(phone string) (string, error) {
	code := util.RandInt(123456, 987654)

	err := dysms.Do(phone, code)
	if err != nil {
		log.Println("userService.LoginSMSCode 发送验证码失败, error: ", err)
		return "", errors.New("发送验证码失败")
	}

	cd := globalModel.SMSCode{
		Phone: phone,
		Code:  strconv.Itoa(code),
		// CreatedAt: time.Now(),
	}
	global.DB.Create(&cd)

	// 	user := User{Name: "Jinzhu", Age: 18, Birthday: time.Now()}

	// result := db.Create(&user) // 通过数据的指针来创建

	// user.ID             // 返回插入数据的主键....
	// result.Error        // 返回 error
	// result.RowsAffected // 返回插入记录的条数
	// global.DB.Update()
	// globalModel.SMSCodeResponse
	// 记录到redis
	_, err = redis.RDB.Set(redis.Ctx, redis.KeySMSCode+fmt.Sprintf("::%s", phone), strconv.Itoa(code), time.Minute*10).Result()
	if err != nil {
		log.Println("userService.LoginSMSCode redis写入失败, error: ", err)
		return "", errors.New("发送验证码失败")
	}

	return strconv.Itoa(code), nil
}

func ComparePassword(username string, password string) bool {
	user, err := userRepo.FindByUserName(username)
	if err != nil {
		return false
	}
	return user.ComparePassword(password)
}

// CheckSMSCode 短信验证码校验
func CheckSMSCode(phone string, code string) error {
	key := redis.KeySMSCode + fmt.Sprintf("::%s", phone)

	result, err := redis.RDB.Get(redis.Ctx, key).Result()
	if err != nil {
		log.Println("userService.CheckSMSCode redis读取失败, error: ", err)
		return errors.New("验证码读取失败")
	}

	if code == result {
		redis.RDB.Del(redis.Ctx, key)
		return nil
	}

	return errors.New("验证码错误")
}

// SendEmailCode 发送邮箱验证码
func SendEmailCode(emailAddr string) (string, error) {
	// 验证邮箱格式
	if !email.ValidateEmail(emailAddr) {
		return "", errors.New("邮箱格式错误")
	}

	// 生成验证码
	code := util.RandInt(123456, 987654)
	codeStr := strconv.Itoa(code)

	// 从配置文件读取邮箱配置
	emailConfig := email.EmailConfig{
		SMTPHost:     config.EmailConfig.SMTPHost,
		SMTPPort:     config.EmailConfig.SMTPPort,
		FromEmail:    config.EmailConfig.FromEmail,
		FromPassword: config.EmailConfig.FromPassword,
		FromName:     config.EmailConfig.FromName,
	}

	// 发送邮件
	err := email.SendEmailCode(emailConfig, emailAddr, codeStr)
	if err != nil {
		log.Println("userService.SendEmailCode 发送验证码失败, error: ", err)
		return "", errors.New("发送验证码失败")
	}

	// 保存到Redis，有效期1分钟（验证码是临时数据，使用Redis更合适）
	_, err = redis.RDB.Set(redis.Ctx, redis.KeyEmailCode+fmt.Sprintf("::%s", emailAddr), codeStr, time.Minute*1).Result()
	if err != nil {
		log.Println("userService.SendEmailCode redis写入失败, error: ", err)
		return "", errors.New("发送验证码失败")
	}

	return codeStr, nil
}

// CheckEmailCode 邮箱验证码校验
func CheckEmailCode(emailAddr string, code string) error {
	key := redis.KeyEmailCode + fmt.Sprintf("::%s", emailAddr)

	result, err := redis.RDB.Get(redis.Ctx, key).Result()
	if err != nil {
		log.Println("userService.CheckEmailCode redis读取失败, error: ", err)
		return errors.New("验证码读取失败")
	}

	if code == result {
		redis.RDB.Del(redis.Ctx, key)
		return nil
	}

	return errors.New("验证码错误")
}

// Exist 检查用户是否存在（公开函数）
func Exist(email string) bool {
	return !(userRepo.Exist(email) != nil)
}

// ExistByNickname 检查昵称（账号）是否存在（公开函数）
func ExistByNickname(nickname string) bool {
	return !(userRepo.ExistByNickname(nickname) != nil)
}

// GenerateUniqueInviteCode 生成唯一的邀请码
func GenerateUniqueInviteCode() (string, error) {
	maxRetries := 10 // 最大重试次数
	for i := 0; i < maxRetries; i++ {
		// 生成8位随机邀请码（字母+数字）
		inviteCode := util.GenerateRandomString(8)

		// 检查邀请码是否已存在，如果不存在（返回错误），说明可以使用
		err := userRepo.ExistByMyInviteCode(inviteCode)
		if err != nil {
			// 邀请码不存在，可以使用
			return inviteCode, nil
		}
	}
	return "", errors.New("生成邀请码失败，请重试")
}

// Register 注册用户（公开函数）
func Register(data userModel.RegisterRequest) (*userModel.User, error) {
	// 数据库注册
	user := userModel.New(data.Email)
	user.Phone = data.Phone
	user.Nickname = data.Nickname
	user.Password = data.Password

	// 如果有填写邀请码，则根据邀请码查找邀请人，并记录到 RefId 和 RefCode
	if data.InviteCode != "" {
		inviter, err := userRepo.FindByMyInviteCode(data.InviteCode)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				return nil, errors.New("邀请码无效")
			}
			log.Println("userService.Register 查询邀请码对应用户失败, error: ", err)
			return nil, errors.New("邀请码校验失败")
		}
		user.RefId = int(inviter.ID)
		user.RefCode = data.InviteCode
	}

	// 生成并设置用户自身的邀请码
	myInviteCode, err := GenerateUniqueInviteCode()
	if err != nil {
		log.Println("userService.Register 生成邀请码失败, error: ", err)
		return nil, errors.New("生成邀请码失败")
	}
	user.MyInviteCode = myInviteCode

	// 加密密码
	if err := user.HashPassword(); err != nil {
		log.Println("userService.Register 密码加密失败, error: ", err)
		return nil, errors.New("密码加密失败")
	}

	_, err = userRepo.Create(user)

	if err != nil {
		return nil, err
	}

	return user, nil
}

func handleRef(user *userModel.User, ip string) {
	// 通过ip查询是否是邀请注册
	refId, err := globalRepo.FindRefIdByIp(ip)
	if err == nil {
		user.RefId = refId
	}

	// 增加邀请人的邀请数
	err = userRepo.IncreaseRefCountById(refId)
	if err != nil {
		log.Printf("userService.handleRef 增加邀请人的邀请数失败, error: %v", err)
	}
}

func UpAvatar(ctx contextModel.Context, avatar userModel.UserAvatarRequest) (url string, err error) {

	image := globalModel.Image{
		UUID:   avatar.Avatar,
		Width:  avatar.Width,
		Height: avatar.Height,
	}

	err = globalRepo.CreateImage(image)

	if err != nil {
		return "", err
	}

	url = config.QiniuConfig.ImgHost + avatar.Avatar

	return
}

func UploadAvatar(ctx contextModel.Context, filePath string) (url string, err error) {
	url, err = globalService.UploadImage(filePath)

	if err != nil {
		return "", errors.New("上传失败")
	}

	return
}

func TagsMap() (map[int]string, error) {
	m := make(map[int]string)
	m[10000] = "学霸"
	m[10010] = "学渣"
	m[10020] = "文艺"
	m[10030] = "体育"
	m[10040] = "音乐"
	m[10050] = "美食"
	m[10060] = "旅行"

	return m, nil
}

func UpdateNickname(ctx contextModel.Context, data userModel.UserNickName) error {
	user := ctx.User

	err := userRepo.Updates(user)

	return err
}

func GetInfo(ctx contextModel.Context) (userModel.InfoResponse, error) {
	user := ctx.User

	memberList, err := memberService.ListMemberConfig()
	if err != nil {
		log.Println("userService.GetInfo 获取会员配置失败, error: ", err)
		return userModel.InfoResponse{}, err
	}

	// 从用户支付信息中获取数据
	jimicoin := user.JimiCoin
	usedcoin := 0     // TODO: 需要添加已使用币字段到Payment结构体
	requestCount := 0 // TODO: 需要添加请求次数字段到Payment结构体

	return userModel.InfoResponse{
		Id:           int(ctx.User.ID),
		Phone:        ctx.User.Phone,
		Email:        ctx.User.Email,
		Nickname:     ctx.User.Nickname,
		MyInviteCode: user.MyInviteCode,
		City:         ctx.User.City,
		VipInfo:      user.GetVipInfo(memberList),
		Payment:      userModel.PaymentResponse{JimiCoin: jimicoin, UsedCoin: usedcoin, RequestCount: requestCount},
		ApiKey:       user.ApiKey,
	}, nil
}

// SelectHotUsersTable 选取热门用户列表, 选择的是当前符合条件的用户列表, 会清空当前备选用户列表
func SelectHotUsersTable() error {
	// 获取用户体型映射表
	attr, err := adminUserService.GetAttrByTypeAndKey("info", "shape")
	if err != nil {
		return err
	}

	// 获取热门用户数量
	count := config.HotUserConfig.Count

	// 清空表
	err = userRepo.DropHotUsers()
	if err != nil {
		log.Printf("userService.SelectHotUsersTable 清空热门用户表失败, error: %v", err)
		return err
	}

	for shapeId := range attr.Range.Enum {
		hotUsers, err := userRepo.SelectHotUser(shapeId, count)
		if err != nil && !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("userService.SelectHotUsersTable 选取热门用户失败: ", err)
			continue
		}

		if len(hotUsers) == 0 {
			continue
		}

		for i := 0; i < len(hotUsers); i++ {
			hotUsers[i].Sort = i + 1
		}

		err = userRepo.SaveHotUsers(hotUsers)
		if err != nil {
			log.Println("userService.SelectHotUsersTable 保存热门用户失败: ", err)
			continue
		}
	}

	return nil
}

// FindAlternativeHotUsers 获取当前备选热门用户
func FindAlternativeHotUsers() (map[int][]userModel.HotUser, error) {
	users, err := userRepo.FindHotUsers()
	if err != nil {
		log.Printf("userService.FindAlternativeHotUsers 查询备选热门用户失败, error: %v", err)
		return nil, errors.New("查询备选热门用户失败")
	}

	results := lo.GroupBy(users, func(user userModel.HotUser) int {
		return user.Shape
	})

	return results, nil
}

// FindReleaseHotUsers 获取当前正式热门用户
func FindReleaseHotUsers() (map[int][]userModel.HotUser, error) {
	result, err := redis.RDB.HGetAll(redis.Ctx, redis.KeyHotUsers).Result()
	if err != nil {
		log.Printf("userService.FindReleaseHotUsers 获取redis数据失败: %s", err)
		return nil, err
	}

	m := make(map[int][]userModel.HotUser)

	for k, v := range result {
		var data []userModel.HotUser
		err := json.Unmarshal([]byte(v), &data)
		if err != nil {
			log.Printf("userService.FindReleaseHotUsers 转换json失败: %s", err)
			continue
		}
		i, err := strconv.ParseInt(k, 10, 64)
		if err != nil {
			log.Println("userService.FindReleaseHotUsers 转换key失败: ", err)
		}

		m[int(i)] = data
	}

	return m, nil
}

// ReleaseHotUsers 发布热门用户
func ReleaseHotUsers(users map[string][]userModel.HotUser) error {

	for k, v := range users {
		j, err := json.Marshal(v)
		if err != nil {
			log.Println("userService.ReleaseHotUsers json转换失败, error: ", err)
			return err
		}

		_, err = redis.RDB.HSet(context.Background(), redis.KeyHotUsers, k, string(j)).Result()
		if err != nil {
			log.Printf("userService.ReleaseHotUsers redis写入失败, error: %v", err)
			return err
		}
	}

	return nil
}

// HomeData 用户主页数据
func HomeData(ctx contextModel.Context, uid int) (*userModel.HomeDataResponse, error) {
	user, err := userRepo.FindByID(uid)
	if err != nil || int(user.ID) != uid {
		if !errors.Is(err, gorm.ErrRecordNotFound) {
			log.Println("userService.HomeData 获取用户数据失败, error: ", err)
		}
		return nil, errors.New("用户不存在")
	}

	memberList, err := memberService.ListMemberConfig()
	if err != nil {
		log.Println("userService.HomeData 获取会员配置失败, error: ", err)
		return nil, err
	}

	result := &userModel.HomeDataResponse{
		Id:            int(user.ID),
		City:          user.City,
		Distance:      user.Position.Distance(ctx.User.Position.Longitude, ctx.User.Position.Latitude),
		LastLoginTime: user.LastRequestTime,
		VipInfo:       user.GetVipInfo(memberList),
	}

	return result, nil
}

func UpdatePosition(ctx contextModel.Context, position positionModel.Position) error {
	city := gaode.ReGeo(position)

	ctx.User.City = city

	return nil
}

func GetUserInfoByPhone(phone string) (*userModel.User, error) {
	user, err := userRepo.Find(phone)
	if err != nil {
		return nil, err
	}
	return user, nil
}

func Logout(user *userModel.User) error {
	log.Println("userService.Logout 注销用户, userID: ", user.ID)
	return userRepo.DeleteByID(cast.ToInt(user.ID))
}

// LogoutToken 退出登录，将token加入黑名单
func LogoutToken(token string) error {
	// 解析token获取过期时间
	jwtToken, err := util.ValidJWT(token)
	if err != nil {
		return errors.New("token无效")
	}

	// 计算token剩余有效时间
	expiresAt := time.Unix(jwtToken.ExpiresAt, 0)
	now := time.Now()
	ttl := expiresAt.Sub(now)

	// 如果token已过期，不需要加入黑名单
	if ttl <= 0 {
		return nil
	}

	// 将token加入Redis黑名单，设置过期时间为token的剩余有效时间
	key := redis.KeyTokenBlacklist + fmt.Sprintf("::%s", token)
	_, err = redis.RDB.Set(redis.Ctx, key, "1", ttl).Result()
	if err != nil {
		log.Println("userService.LogoutToken 加入黑名单失败, error: ", err)
		return errors.New("退出登录失败")
	}

	return nil
}

func UpdateUserVip(tx *gorm.DB, uid int, days int) error {
	user, err := userRepo.FindByID(uid)
	if err != nil {
		return err
	}
	log.Println("userService.UpdateUserVip 更新用户vip, uid: ", uid, "days: ", days)
	if user.VipTime == nil {
		vipTime := time.Now().AddDate(0, 0, days)
		user.VipTime = &vipTime
	} else {
		newVipTime := user.VipTime.AddDate(0, 0, days)
		user.VipTime = &newVipTime
	}
	return userRepo.UpdateVipExpireTimeTx(tx, user)
}

// UpdatePassword 修改密码
func UpdatePassword(ctx contextModel.Context, data userModel.UpdatePasswordRequest) error {
	user := ctx.User

	// 验证旧密码
	if !user.ComparePassword(data.OldPassword) {
		return errors.New("旧密码错误")
	}

	// 设置新密码并加密
	user.Password = data.NewPassword
	if err := user.HashPassword(); err != nil {
		log.Println("userService.UpdatePassword 密码加密失败, error: ", err)
		return errors.New("密码加密失败")
	}

	// 只更新密码字段
	err := userRepo.UpdatePassword(user)
	if err != nil {
		log.Println("userService.UpdatePassword 更新密码失败, error: ", err)
		return errors.New("更新密码失败")
	}

	return nil
}
