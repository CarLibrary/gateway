package carlibrary

import (
	"CarLibrary/gateway/serializer"
	"context"
	"encoding/json"
	"errors"
	car "github.com/CarLibrary/proto/carlibrary"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
	"io"
	"log"
	"os"
	"time"
)


func FindALLCarBand(ctx *gin.Context)  {

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("CARLIBRARY_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client :=car.NewCarLibraryServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cl,err:=client.FindALLCarBand(c,&car.Empty{})
	if err != status.Error(codes.OK,"ok") {
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

func FindAllCarSeries(ctx *gin.Context)  {

	band:=ctx.Param("band")
	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("CARLIBRARY_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client :=car.NewCarLibraryServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cl,err:=client.FindAllCarSeries(c,&car.CarSeriesRequest{CarBand: band})
	if err != status.Error(codes.OK,"ok") {
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

func FindAllCarModel(ctx *gin.Context)  {

	band:=ctx.Param("band")
	series:=ctx.Param("series")
	if band == "" || series == "" {
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        "请重试",
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg:errors.New("bind error").Error()},
		})
		return
	}

	// Set up a connection to the server.
	conn, err := grpc.Dial(os.Getenv("CARLIBRARY_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client :=car.NewCarLibraryServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cl,err:=client.FindAllCarModel(c,&car.CarModelRequest{
		CarBand:   band,
		CarSeries: series,
	})
	if err != status.Error(codes.OK,"ok") {
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