package middleware

import (
	"CarLibrary/gateway/serializer"
	"context"
	account "github.com/CarLibrary/proto/account"
	"github.com/gin-gonic/gin"
	"google.golang.org/grpc"
	"log"
	"os"
	"time"
)


func CheckToken() gin.HandlerFunc  {
	return func(ctx *gin.Context) {
		toeknstr:=ctx.Request.Header.Get("Authorization")


		// Set up a connection to the server.
		conn, err := grpc.Dial(os.Getenv("ACCOUNT_ADRR"), grpc.WithInsecure(), grpc.WithBlock())
		if err != nil {
			log.Fatalf("did not connect: %v", err)
		}
		defer conn.Close()
		client := account.NewAccountServiceClient(conn)

		c, cancel := context.WithTimeout(context.Background(), 3*time.Second)
		defer cancel()
		u,err:=client.CheckToken(c,&account.TokenRequest{Token:toeknstr })
		if err != nil {
			ctx.JSON(200,serializer.CommonResponse{
				Statuscode: serializer.GRPC_ERROR,
				Data:       nil,
				Msg:        "服务器内部错误",
				Error:      serializer.ErrorResponse{serializer.MIDDLEWARE,err.Error()},
			})
		}

		ctx.Next()
		ctx.Set("id",u.Id)
		ctx.Set("username",u.Username)
	}
}
