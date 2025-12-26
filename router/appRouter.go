package router

import (
	"fmt"
	"jiyu/api/appVersionApi"
	"jiyu/api/chargeApi"
	"jiyu/api/feedbackApi"
	"jiyu/api/globalApi"
	"jiyu/api/inviteApi"
	"jiyu/api/payApi"
	"jiyu/api/redeemCodeApi"
	"jiyu/api/taskApi"
	"jiyu/api/usageDetailApi"
	"jiyu/api/userApi"
	"jiyu/api/videoApi"
	"jiyu/api/watermarkApi"
	"jiyu/global"
	"jiyu/middleware"
	"time"

	"github.com/gin-gonic/gin"
)

func NewRouter() *gin.Engine {
	router := gin.New()

	router.
		Use(gin.Recovery()).
		Use(middleware.Cors()).
		Use(middleware.Trace()).
		Use(
			gin.LoggerWithFormatter(func(param gin.LogFormatterParams) string {
				traceId := param.Request.Context().Value(global.TraceIDKey{})
				return fmt.Sprintf("%s > %s - [%s] %s \"%s\" %s %d %s \"%s\" %s \n",
					traceId,
					param.ClientIP,
					param.TimeStamp.Format(time.RFC1123),
					param.Method,
					param.Path,
					param.Request.Proto,
					param.StatusCode,
					param.Latency,
					param.Request.UserAgent(),
					param.ErrorMessage,
				)
			}))

	routers(router)

	return router
}

func routers(r *gin.Engine) {
	// 下载关联用户ID
	r.GET("/download/ref/:uid", globalApi.DownloadRef)

	publicRouter := r.Group("/api")
	{
		// 服务响应测试
		publicRouter.GET("/ping", userApi.Ping)
		// 登录
		publicRouter.POST("/login", userApi.Login)
		// 注册
		publicRouter.POST("/register", userApi.Register)
		// 发送登录验证码
		publicRouter.POST("/sendLoginCode", userApi.LoginSMSCode)
		// 发送注册邮箱验证码
		publicRouter.POST("/sendEmailCode", userApi.SendEmailCode)
		// 发送修改密码邮箱验证码
		publicRouter.POST("/sendUpdatePasswordEmailCode", userApi.SendUpdatePasswordEmailCode)
		// 忘记密码（未登录时使用）
		publicRouter.POST("/forgotPassword", userApi.ForgotPassword)
		// 获取所有标签
		publicRouter.GET("/user/tagsMap", userApi.TagsMap)
		// 获取协议
		publicRouter.GET("/agreement", globalApi.FindDoc)
		// 获取省市区
		publicRouter.GET("/area", globalApi.FindArea)
		// 获取用户属性映射表
		publicRouter.GET("/user/attrMap", globalApi.UserAttributesMap)
		// 获取用户标签映射表
		publicRouter.GET("/user/tagMap", globalApi.UserTagsMap)
		// 一键登录
		publicRouter.POST("/user/autoLogin", userApi.AutoLogin)
		// 获取版本信息
		publicRouter.GET("/app/version", appVersionApi.GetVersion)
	}

	userValidateRouter := r.Group("/api/user/noauth").Use(middleware.JWTAuth())
	{
		//注销账号
		userValidateRouter.POST("/logout", userApi.Logout)
	}

	userRouter := r.Group("/api/user").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 退出登录
		userRouter.POST("/logout", userApi.LogoutToken)
		// 上传头像
		userRouter.POST("/avatar", userApi.UpAvatar)
		// 获取自己的基本信息
		userRouter.GET("/info", userApi.GetInfo)
		// 更新用户昵称
		userRouter.POST("/updateNickname", userApi.UpadteNickName)
		// 用户是否被禁言
		userRouter.GET("/isBan", userApi.IsBan)
		// 获取用户主页数据
		userRouter.GET("/home/:uid", userApi.UserHomeData)
		// 上传用户信息
		userRouter.POST("/position", userApi.UpdatePosition)
		// 修改密码
		userRouter.POST("/password", userApi.UpdatePassword)
		// 创建使用明细
		userRouter.POST("/usage/detail", usageDetailApi.CreateUsageDetail)
		// 获取使用明细列表
		userRouter.GET("/usage/list", usageDetailApi.GetUsageDetailList)
	}

	watermarkRouter := r.Group("/api/watermark").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 去水印接口
		watermarkRouter.POST("/remove", watermarkApi.RemoveWatermark)
	}

	videoRouter := r.Group("/api/video").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 图生视频接口
		videoRouter.POST("/generate", videoApi.GenerateVideo)
		// 查询视频任务接口
		videoRouter.GET("/:videoId", videoApi.GetVideo)
	}

	taskRouter := r.Group("/api/task").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 创建任务
		taskRouter.POST("", taskApi.CreateTask)
		// 更新任务
		taskRouter.PUT("/:taskId", taskApi.UpdateTask)
		// 获取任务列表
		taskRouter.GET("/list", taskApi.GetTaskList)
	}

	chargeRouter := r.Group("/api/charge").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 充值类型
		chargeRouter.GET("/list", chargeApi.List)
		// 获取会员等级
		chargeRouter.GET("/member/list", chargeApi.MemberList)
		// 创建订单
		chargeRouter.POST("/order/create", chargeApi.CreateOrder)
		// 兑换兑换码
		chargeRouter.POST("/redeem", redeemCodeApi.RedeemCode)
		// 充值明细列表
		chargeRouter.GET("/redeem/list", chargeApi.RedeemList)
	}

	// 创建兑换码接口（无需鉴权）
	chargePublicRouter := r.Group("/api/charge")
	{
		// 创建兑换码
		chargePublicRouter.POST("/redeem/create", redeemCodeApi.CreateRedeemCode)
	}

	chargeNotifyRouter := r.Group("/api/charge")
	{
		// 充值回调
		chargeNotifyRouter.POST("/order/notify", chargeApi.Notify)
		chargeNotifyRouter.POST("/apple/notify", chargeApi.AppleNotify)
		chargeNotifyRouter.POST("/wechat/notify", chargeApi.WechatNotify)
	}

	globalRouter := r.Group("/api").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 上传图片
		// globalRouter.POST("/images", globalApi.UploadImages)
		globalRouter.GET("/image/uploadToken", globalApi.ImgUploadToken)
		// 投诉
		globalRouter.POST("/report", globalApi.Report)

	}

	paymentRouter := r.Group("/api/payment").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 新增付款账户
		paymentRouter.POST("/account/add", payApi.AddAccount)
		// 获取提现账户列表
		paymentRouter.GET("/account/list", payApi.AccountList)
		// 删除付款账户
		paymentRouter.POST("/account/del", payApi.DeleteAccount)
		// 获取提现设置
		paymentRouter.GET("/setting", payApi.Setting)
		// 申请提现
		paymentRouter.POST("/withdrawal", payApi.Withdrawal)
		// 获取提现列表
		paymentRouter.GET("/withdrawal/list", payApi.WithdrawalList)
	}

	feedbackRouter := r.Group("/api/feedback").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 意见反馈
		feedbackRouter.POST("/add", feedbackApi.Add)
		// 意见已读
		feedbackRouter.GET("/read/:feedbackId", feedbackApi.Read)
		// 意见反馈列表
		feedbackRouter.GET("list", feedbackApi.List)
		// 常见问题列表
		feedbackRouter.GET("/qas", feedbackApi.QAList)
	}
	inviteRouter := r.Group("/api/invite").Use(middleware.JWTAuth()).Use(middleware.Context()).Use(middleware.Banned())
	{
		// 邀请排行榜
		inviteRouter.GET("/ranking", inviteApi.Ranking)
		// 邀请记录
		inviteRouter.GET("/record", inviteApi.Record)
		// 获取前N名排行榜
		inviteRouter.GET("/top/:n", inviteApi.TopN)
		// 获取用户排名
		inviteRouter.GET("/user/rank", inviteApi.UserRank)
	}
}
