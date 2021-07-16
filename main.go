package main

import (
	"CarLibrary/gateway/account"
	"CarLibrary/gateway/carlibrary"
	"CarLibrary/gateway/middleware"
	"CarLibrary/gateway/score"
	"CarLibrary/gateway/testdrive"
	"github.com/gin-gonic/gin"
)

func main() {
	r:=gin.Default()

	//注册
	r.POST("/signup",account.Signup)

	//登录
	r.POST("login",account.Login)

	//车型大全
	//查看所有品牌
	r.GET("/band",carlibrary.FindALLCarBand)
	//查看品牌的全部车系
	r.GET("/:band/series",carlibrary.FindAllCarSeries)
	//查看某品牌的某车系的全部车型
	r.GET("/:band/:series/model",carlibrary.FindAllCarModel)


	app:=r.Group("/v1",middleware.CheckToken())

	//个人中心
	//查看个人信息
	app.GET("/user",account.GetUserInfo)

	//评分
	//查看我的评分
	app.GET("/user/car_point",score.FindMYScore)
	//修改评分
	app.PUT("/user/car_point",score.ModifyScore)
	//评分
	app.POST("/user/car_point",score.MakeScore)

	//试驾
	//查看我的试驾信息
	app.GET("/user/test_drive",testdrive.FindMyTestDrive)
	//试驾
	app.POST("/user/test_drive",testdrive.TestDrive)

	r.Run(":8080")
}
