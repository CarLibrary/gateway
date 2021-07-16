package account

import (
	"CarLibrary/gateway/serializer"
	"context"
	"errors"
	account "github.com/CarLibrary/proto/account"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)

type SignUp struct {
	Username string `json:"username" binding:"required;max=6"`
	Password string `json:"password" binding:"required;min=6;max=16"`
	Head_url string `json:"head_url" binding:"max=255"`
	Sgin string `json:"sign" binding:"max=20"`
}

type LogIN struct {
	Username string `json:"username" binding:"required;max=6"`
	Password string `json:"password" binding:"required;min=6;max=16"`
}

// Signup 注册
func Signup(ctx *gin.Context)  {
	var s =&SignUp{}
	if err :=ctx.ShouldBind(s);err!=nil{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        serializer.BIND,
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg: err.Error()},
		})
		return
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("ACCOUNT_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	u,err:=client.Signup(c,&account.SignupRequset{
		Username: s.Username,
		Password: s.Password,
		HeadUrl:  &s.Head_url,
		Sgin:     &s.Sgin,
	})
	if err != nil {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.GRPC_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{serializer.OTHER,err.Error()},
		})
		return
	}
	ctx.JSON(200,&serializer.CommonResponse{
		Statuscode: serializer.OK,
		Data:       u,
		Msg:        "注册成功",
		Error:      serializer.ErrorResponse{},
	})
	return
}

//登录
func Login(ctx *gin.Context)  {
	var l = &LogIN{}
	if err:=ctx.ShouldBind(l);err!=nil{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        serializer.BIND,
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg: err.Error()},
		})
		return
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("ACCOUNT_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	t,err:=client.Login(c,&account.LoginRequest{
		Username: l.Username,
		Password: l.Password,
	})
	if err != nil {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.GRPC_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{serializer.OTHER,err.Error()},
		})
		return
	}
	ctx.JSON(200,&serializer.CommonResponse{
		Statuscode: serializer.OK,
		Data:       t,
		Msg:        "登录成功",
		Error:      serializer.ErrorResponse{},
	})
	return
}

// GetUserInfo 查看个人信息
func GetUserInfo(ctx *gin.Context)  {
	id,ok:=ctx.Get("id")
	if !ok {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.TOKEN_ERROR,
			Data:       nil,
			Msg:        "请重新登录！",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg: errors.New("token error").Error()},
		})
		return
	}
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("ACCOUNT_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := account.NewAccountServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()
	res,err:=client.GetUserInfo(c,&account.InfoRequest{Id: id.(int32)})
	if err != nil {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.GRPC_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:err.Error()},
		})
		return
	}
	ctx.JSON(200,&serializer.CommonResponse{
		Statuscode: serializer.OK,
		Data:       res,
		Msg:        "ok",
		Error:      serializer.ErrorResponse{},
	})
	return
}

