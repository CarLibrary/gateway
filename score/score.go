package score

import (
	"CarLibrary/gateway/serializer"
	"context"
	"encoding/json"
	"errors"
	score "github.com/CarLibrary/proto/score"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"io"
	"log"
	"os"
	"time"
)

type ScoreReq struct {
	Score float32 `json:"score" binding:"required"`
}

// MakeScore 打分
func MakeScore(ctx *gin.Context)  {

	id,ok:=ctx.Get("id")
	if !ok{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.TOKEN_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:errors.New("token error").Error()},
		})
		return
	}
	s:=&ScoreReq{}
	if err:=ctx.ShouldBind(s);err!=nil{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        serializer.BIND,
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg:err.Error()},
		})
		return
	}

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
	conn, err := grpc.Dial(os.Getenv("SCORE_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := score.NewScoreServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res,err:=client.MakeScore(c,&score.ScoreRequest{
		Userid:    id.(int32),
		CarBand:   band,
		CarSeries: series,
		Score:     s.Score,
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

// ModifyScore 修改评分
func ModifyScore(ctx *gin.Context)  {

	id,ok:=ctx.Get("id")
	if !ok{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.TOKEN_ERROR,
			Data:       nil,
			Msg:        "服务器内部错误",
			Error:      serializer.ErrorResponse{ErrType: serializer.OTHER,Msg:errors.New("token error").Error()},
		})
		return
	}
	s:=&ScoreReq{}
	if err:=ctx.ShouldBind(s);err!=nil{
		ctx.JSON(200,&serializer.CommonResponse{
			Statuscode: serializer.BIND_ERROR,
			Data:       nil,
			Msg:        serializer.BIND,
			Error:      serializer.ErrorResponse{ErrType: serializer.BIND,Msg:err.Error()},
		})
		return
	}
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
	conn, err := grpc.Dial(os.Getenv("SCORE_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
	if err != nil {
		log.Fatalf("did not connect: %v", err)
	}
	defer conn.Close()
	client := score.NewScoreServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	res,err:=client.ModifyScore(c,&score.ScoreRequest{
		Userid:    id.(int32),
		CarBand:   band,
		CarSeries: series,
		Score:     s.Score,
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

func FindMYScore(ctx *gin.Context)  {
	id,ok:=ctx.Get("id")
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
	client := score.NewScoreServiceClient(conn)

	c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
	defer cancel()

	cl,err:=client.FindMYScore(c,&score.MyScoresRequest{Userid: id.(int32)})
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