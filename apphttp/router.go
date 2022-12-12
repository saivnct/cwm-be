package apphttp

import (
	"encoding/hex"
	"errors"
	"fmt"
	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
	"golang.org/x/exp/slices"
	"log"
	"net/http"
	"sol.go/cwm/appws"
	"sol.go/cwm/dao"
	"sol.go/cwm/model"
	"sol.go/cwm/utils"
	"strconv"
)

var (
	HTTPAccessDeniedErr         = errors.New("Access Denied")
	HTTPUnauthenticatedTokenErr = errors.New("Invalid Token")
	HTTPUnauthenticateUserdErr  = errors.New("Invalid User")
)

func GinMiddleware(allowOrigin string) gin.HandlerFunc {
	return func(c *gin.Context) {
		c.Writer.Header().Set("Access-Control-Allow-Origin", allowOrigin)
		c.Writer.Header().Set("Access-Control-Allow-Credentials", "true")
		c.Writer.Header().Set("Access-Control-Allow-Methods", "POST, OPTIONS, GET, PUT, DELETE")
		c.Writer.Header().Set("Access-Control-Allow-Headers", "Accept, Authorization, Content-Type, Content-Length, X-CSRF-Token, Token, session, Origin, Host, Connection, Accept-Encoding, Accept-Language, X-Requested-With")

		if c.Request.Method == http.MethodOptions {
			c.AbortWithStatus(http.StatusNoContent)
			return
		}

		c.Request.Header.Del("Origin")

		c.Next()
	}
}

func WSMiddleware() gin.HandlerFunc {
	return func(ctx *gin.Context) {

		jwtToken := ctx.Query(appws.WS_CRED_AUTH)
		//log.Println("WS - WSMiddleware", jwtToken)

		payload, err := hex.DecodeString(jwtToken)
		if err != nil {
			log.Println("WS - Invalid decode jwtToken", err)
			ctx.AbortWithError(http.StatusUnauthorized, HTTPUnauthenticatedTokenErr)
			return
		}

		claims, err := utils.ParseJWTToken(string(payload))
		if err != nil {
			log.Println("WS - Invalid jwtToken", err)
			ctx.AbortWithError(http.StatusUnauthorized, HTTPUnauthenticatedTokenErr)
			return
		}

		phoneFull := claims.Subject
		sessionId := claims.ID

		user, err := dao.GetUserDAO().FindByPhoneFull(ctx, phoneFull)
		if err != nil {
			log.Println("WSMiddleware - not found user", phoneFull)
			ctx.AbortWithError(http.StatusForbidden, HTTPAccessDeniedErr)
			return
		}

		idx := slices.IndexFunc(user.Sessions, func(c model.UserSession) bool { return c.SessionId == sessionId })
		if idx < 0 {
			log.Println("WSMiddleware - not found session:", sessionId)
			ctx.AbortWithError(http.StatusForbidden, HTTPAccessDeniedErr)
			return
		}

		//fmt.Println("on WS connect, path:", ctx.FullPath())
		//fmt.Println("phoneFull:", phoneFull)
		//ctx.Set(appws.WS_CTX_KEY_USER, user)

		ctx.Next()
	}
}

func StartServer(httpPort int, ws *socketio.Server) (*gin.Engine, error) {
	//testHTTPController := TestHTTPController{}
	//fileHTTPController := FileHTTPController{}

	gin.SetMode(gin.ReleaseMode)

	// Creates a gin router with default middleware:
	// logger and recovery (crash-free) middleware
	router := gin.Default()
	//	https://pkg.go.dev/github.com/gin-gonic/gin#readme-don-t-trust-all-proxies
	router.SetTrustedProxies([]string{"localhost", "10.61.60.41", "192.168.2.1", "192.168.1.7"})
	router.Use(GinMiddleware("*"))
	//router.Use(GinMiddleware("http://localhost:63343"))

	wsRouter := router.Group("ws")
	wsRouter.Use(WSMiddleware())
	{
		wsRouter.GET("/*any", gin.WrapH(ws))
		wsRouter.POST("/*any", gin.WrapH(ws))
	}

	router.GET("/healthCheck", func(ctx *gin.Context) {
		ctx.JSON(http.StatusOK, gin.H{"status": "success"})
	})

	//download media msg
	//fileRouter := router.Group("/media")
	//{
	//	fileRouter.GET("/get/:msgId/:fileId", fileHTTPController.getFile)
	//}

	//testRouter := router.Group("/test")
	//{
	//	testRouter.GET("/sendmsg/:name", testHTTPController.SendMsg)
	//	testRouter.GET("/test", testHTTPController.Test)
	//}

	go func() {
		go func() {
			err := ws.Serve()
			if err != nil {
				log.Fatalf("Failed to serve ws: %v", err)
			}
		}()
		defer ws.Close()

		fmt.Printf("Starting http server on: %v\n", httpPort)
		err := router.Run(":" + strconv.Itoa(httpPort))
		if err != nil {
			log.Fatalf("Failed to serve http: %v", err)
		}
	}()

	return router, nil
}
