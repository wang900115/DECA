package main

import (
	"context"
	"log"
	"os"

	"github.com/joho/godotenv"
	"github.com/robfig/cron"
	"github.com/wang900115/DESA/internal/adapter/controller"
	"github.com/wang900115/DESA/internal/adapter/route"
	"github.com/wang900115/DESA/internal/application/usecase"
	"github.com/wang900115/DESA/lib/bootstrap"
	"github.com/wang900115/DESA/lib/common/middleware"
	middlewareCORS "github.com/wang900115/DESA/lib/common/middleware/cors"
	middlewareJWT "github.com/wang900115/DESA/lib/common/middleware/jwt"
	middlewareLOGGER "github.com/wang900115/DESA/lib/common/middleware/logger"
	middlewarePermission "github.com/wang900115/DESA/lib/common/middleware/role"
	middlewareSecure "github.com/wang900115/DESA/lib/common/middleware/secure_header"
	response "github.com/wang900115/DESA/lib/common/response/json"
	"github.com/wang900115/DESA/lib/common/router"
	"github.com/wang900115/DESA/lib/implement"
)

func main() {
	conf := bootstrap.NewConfig()

	err := godotenv.Load()
	if err != nil {
		log.Fatal("Error loading .env file")
	}

	secretKey := os.Getenv("SECRET_KEY")

	dbGroup := bootstrap.NewDBGroup(conf)
	redisGroup := bootstrap.NewRedisGroup(conf)

	syslogger := bootstrap.NewLogger(bootstrap.NewSysLoggerOption(conf))
	applogger := bootstrap.NewLogger(bootstrap.NewAppLoggerOption(conf))

	ctx := context.Background()

	c := cron.New()

	if err := c.AddFunc("*/10 * * * *", func() {
		redisGroup.HeadlthCheck(ctx)
	}); err != nil {
		log.Fatalf("Failed to add cron redis-health-check func: %v", err)
	}

	if err := c.AddFunc("*/20 * * * *", func() {
		dbGroup.HeadlthCheck(ctx)
	}); err != nil {
		log.Fatalf("Failed to add cron postgre-health-check func: %v", err)
	}

	c.Start()

	cr := implement.NewChannelReadRepository(dbGroup, redisGroup, syslogger)
	cw := implement.NewChannelWriteRepository(dbGroup, redisGroup, syslogger)

	mr := implement.NewMessageReadRepository(dbGroup, redisGroup, syslogger)
	mw := implement.NewMessageWriteRepository(dbGroup, redisGroup, syslogger)

	ur := implement.NewUserReadRepository(dbGroup, redisGroup, syslogger)
	uw := implement.NewUserWriteRepository(dbGroup, redisGroup, syslogger)

	tu := implement.NewTokenAuthRepository(redisGroup, syslogger)

	cp := implement.NewChannelP2PService()
	up := implement.NewUserP2PService()

	channel := usecase.NewChannelUsecase(&cr, &cw)
	p2p := usecase.NewP2PUsecase(&cp, &up)
	message := usecase.NewMessageUsecase(&mr, &mw)
	user := usecase.NewUserUsecase(&ur, &uw, &tu, secretKey)

	resp := response.NewJSONResponse(applogger)

	channelCon := controller.NewChannelController(channel, p2p, resp)
	channelMessageCon := controller.NewChannelMessageController(message, resp)
	channelUserCon := controller.NewChannelUserController(channel, p2p, resp)
	messageCon := controller.NewMessageController(message, resp)
	userCon := controller.NewUserController(user, p2p, resp)
	userChannelCon := controller.NewUserChannelController(channel, resp)
	userChannelMessageCon := controller.NewUserChannelMessageController(message, resp)

	midCORS := middlewareCORS.NewCORS(middlewareCORS.NewOption(conf))
	midJWT := middlewareJWT.NewJWT(resp, &tu, secretKey)
	midRole := middlewarePermission.NewPermission(resp, &tu, secretKey)
	midLog := middlewareLOGGER.NewLogger(applogger)
	// midRate := middlewareRate.NewRateLimiter(middlewareRate.NewOption(conf))
	midSecure := middlewareSecure.NewSecureHeader()

	userRoute := route.NewUserRouter(userCon, midJWT)
	userChannelRoute := route.NewUserChannelRouter(userChannelCon, midJWT)
	userChannelMessageRoute := route.NewUserChannelMessageRouter(userChannelMessageCon, midJWT)
	channelRoute := route.NewChannelRouter(channelCon)
	channelUserRoute := route.NewChannelUserRouter(channelUserCon, midJWT, midRole)
	channelMessageRoute := route.NewChannelMessageRouter(channelMessageCon)
	messageRoute := route.NewMessageRouter(messageCon, midJWT)

	server := bootstrap.NewServer(
		[]router.IRoute{
			userRoute,
			userChannelRoute,
			userChannelMessageRoute,
			channelRoute,
			channelUserRoute,
			channelMessageRoute,
			messageRoute,
		},
		[]middleware.IMiddleware{
			midCORS,
			midSecure,
			midLog,
		},
	)

	bootstrap.Run(server, bootstrap.NewServerOption(conf), c)

}
