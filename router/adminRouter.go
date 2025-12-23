package router

import (
	"jiyu/api/adminChargeApi"
	"jiyu/api/adminHomeApi"
	"jiyu/api/adminInviteApi"
	"jiyu/api/adminMemberApi"
	"jiyu/api/adminRouterApi"
	"jiyu/api/adminSettingApi"
	"jiyu/api/adminTransactionApi"
	"jiyu/api/adminUserApi"
	"jiyu/api/appVersionApi"
	"jiyu/api/globalApi"
	"jiyu/api/payApi"
	"jiyu/api/userApi"
	"jiyu/middleware"

	"github.com/gin-gonic/gin"
)

func NewAdminRouter() *gin.Engine {
	router := gin.Default()

	router.Use(middleware.Cors())

	adminRouters(router)

	return router
}

func adminRouters(r *gin.Engine) {
	publicRouter := r.Group("/api").Use(middleware.Context())
	{
		publicRouter.GET("/app/version", appVersionApi.GetVersion)
		publicRouter.POST("/app/version/edit", appVersionApi.UpdateVersion)
		publicRouter.GET("/ping", adminUserApi.Ping)
		publicRouter.POST("/login", adminUserApi.Login)
		publicRouter.GET("/user/attrMap", globalApi.UserAttributesMap)
		publicRouter.GET("/qiniu/upload/token", globalApi.QiniuUploadToken)
	}

	homeRouter := r.Group("/api").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		homeRouter.GET("/home", adminHomeApi.Home)
	}

	// 前端路由相关
	routerRouter := r.Group("/api/router").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		routerRouter.GET("/list", adminRouterApi.RouterList)
		routerRouter.POST("/save", adminRouterApi.SaveRouter)
		routerRouter.POST("/del", adminRouterApi.DeleteRouter)
	}

	userRouter := r.Group("/api/user").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		userRouter.GET("/admin/list", adminUserApi.AdminUserList)
		userRouter.POST("/admin/add", adminUserApi.AddAdminUser)
		userRouter.POST("/admin/del", adminUserApi.DelAdminUser)
		userRouter.POST("/admin/password", adminUserApi.ResetAdminUserPassword)
		userRouter.POST("/admin/auth", adminUserApi.ResetAdminUserAuth)
		userRouter.POST("/changePwd", adminUserApi.ChangePwd)
		userRouter.GET("/info", adminUserApi.Info)
		userRouter.POST("/logout", adminUserApi.Logout)
		userRouter.GET("/list", adminUserApi.UserList)
		userRouter.GET("/:uid", adminUserApi.UserByUid)
		userRouter.POST("/ban", adminUserApi.Ban)
		userRouter.POST("/delete", adminUserApi.Delete)
		userRouter.GET("/detail/:uid", adminUserApi.Detail)
		userRouter.POST("/update", adminUserApi.Update)
		userRouter.GET("/attr", adminUserApi.GetUserAttrs)
		userRouter.GET("/attr/:type/:key", adminUserApi.GetUserAttrByKey)
		userRouter.POST("/attr", adminUserApi.UpdateUserAttrs)
		userRouter.GET("/tags", adminUserApi.GetUserTags)
		userRouter.POST("/updateTag", adminUserApi.UpdateUserTags)
		userRouter.POST("/createTag", adminUserApi.CreateUserTags)
		userRouter.GET("/deleteTag", adminUserApi.DeleteUserTag)

		userRouter.GET("/hot/reselect", userApi.ReselectHotUser)             // 重新选择热门用户
		userRouter.GET("/hot/alternative/list", userApi.AlternativeHotUsers) // 获取备选热门用户
		userRouter.GET("/hot/release/list", userApi.ReleaseHotUsers)         // 获取正式热门用户
		userRouter.POST("/hot/release", userApi.UpdateReleaseHotUsers)       // 发布正式热门用户
	}

	settingRouter := r.Group("/api/setting").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		settingRouter.GET("", adminSettingApi.Get)
		settingRouter.POST("", adminSettingApi.Update)
	}

	transactionRouter := r.Group("/api/transaction").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		transactionRouter.GET("/list", adminTransactionApi.List)
	}

	payRouter := r.Group("/api/pay").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		payRouter.GET("/list", payApi.ListForAdmin)
		payRouter.GET("/list/withdrawal", payApi.WithdrawalListForAdmin)
		payRouter.GET("/withdrawal/edit", payApi.EditWithdrawal)
	}

	smsRouter := r.Group("/api/sms").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		smsRouter.GET("/list", globalApi.SMSCodeList)
	}

	docRouter := r.Group("/api/doc").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		docRouter.GET("", globalApi.FindDoc)
		docRouter.POST("", globalApi.SaveDoc)
	}

	chargeRouter := r.Group("/api/charge").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		chargeRouter.GET("list", adminChargeApi.ChargeConfigList)
		chargeRouter.POST("", adminChargeApi.UpdateChargeConfig)
	}

	memberRouter := r.Group("/api/member").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		memberRouter.GET("/config/list", adminMemberApi.MemberConfigList)
		memberRouter.GET("/list", adminMemberApi.MemberList)
		memberRouter.POST("/edit", adminMemberApi.EditMember)
		memberRouter.POST("/config/edit", adminMemberApi.UpdateMemberConfig)
	}

	inviteRouter := r.Group("/api/invite").Use(middleware.JWTAuth()).Use(middleware.ContextAdmin())
	{
		inviteRouter.GET("/list", adminInviteApi.InviteList)
		inviteRouter.GET("/superior/:uid", adminInviteApi.InviteSuperior) // 查询上级
		inviteRouter.GET("/inferior/:uid", adminInviteApi.InviteInferior) // 查询下级
		inviteRouter.GET("/rank", adminInviteApi.InviteRank)              // 查询排行榜
	}

}
