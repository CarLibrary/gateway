package testdrive

import (
	"CarLibrary/gateway/serializer"
	"context"
	"encoding/json"
	"errors"
	test "github.com/CarLibrary/proto/testdrive"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

type TestDriveReq struct {

	//品牌
	Carband string `json:"carband" binding:"required"`
	//车系
	CarSeries string `json:"car_series" binding:"required"`
	//城市
	City string `json:"city" binding:"required"`
	//姓名
	Username string `json:"username" binding:"required"`
	//手机号
	PhoneNumber string `json:"phone_number" binding:"required"`
}

// TestDrive 试驾
func TestDrive(ctx *gin.Context) {
	id, ok := ctx.Get("id")
	if !ok{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.TOKEN_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:errors.New("token error").Error()},
		})
		return
	}

	var t = &TestDriveReq{}
	if err:=ctx.ShouldBind(t);err!=nil{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        serializer.BIND,
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg:err.Error()},
		})
		return
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("SCORE_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := test.NewTestDriveServerClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res,err:=client.TestDrive(c,&test.TestDriveRequest{
		Userid:      id.(int32),
		Carband:     t.Carband,
		CarSeries:   t.CarSeries,
		City:        t.City,
		Username:    t.Username,
		PhoneNumber: t.PhoneNumber,
	})
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

//查看我的试驾
func FindMyTestDrive(ctx *gin.Context)  {
	id, ok := ctx.Get("id")
	if !ok{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.TOKEN_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:errors.New("token error").Error()},
		})
		return
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("SCORE_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := test.NewTestDriveServerClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cl,err:=client.FindMyTestDrive(c,&test.MyTestDriveRequest{Userid: id.(int32)})
	if err != nil {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.GRPC_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:err.Error()},
		})
		return
	}

	ctx.Status(200)
	for {
		data, err := cl.Recv()
		d,_:=json.Marshal(data)
		if err == io.EOF {
			break
		}
		if err != nil {
			continue
		}
		if _,err:=ctx.Writer.Write(d);err!=nil{
			continue
		}
	}
	return
}