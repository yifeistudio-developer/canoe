package route

import (
	"canoe/internal/handler"
	. "canoe/internal/model"
	"canoe/internal/service"
	"github.com/kataras/iris/v12"
	"github.com/kataras/iris/v12/context"
	"gorm.io/gorm"
)

func SetupRoutes(app *iris.Application, db *gorm.DB) {
	party := app.Party("/canoe/api")

	party.Get("/", func(ctx iris.Context) {
		_ = ctx.JSON(Success(nil))
	})

	users := party.Party("/users")
	UserRoutes(users, &service.UserService{Db: db})

	sessions := party.Party("/sessions")
	SessionRoutes(sessions, &service.SessionService{Db: db})

	party.Get("/channel/{accessToken:string}", handler.HandleChannel())

	party.Get("/chat/{accessToken:string}", func(ctx *context.Context) {
		accessToken := ctx.Params().GetString("accessToken")
		server := handler.NewChatServer(accessToken)
		_, err := server.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil, nil)
		if err != nil {
			return
		}
	})

	//party.Post("/live", func(ctx iris.Context) {
	//	var upgrader = websocket.Upgrader{
	//		CheckOrigin: func(r *http.Request) bool {
	//			return true
	//		},
	//	}
	//	conn, err := upgrader.Upgrade(ctx.ResponseWriter(), ctx.Request(), nil)
	//	if err != nil {
	//		log.Println("Failed to upgrade connection:", err)
	//		return
	//	}
	//	defer func(conn *websocket.Conn) {
	//		err := conn.Close()
	//		if err != nil {
	//		}
	//	}(conn)
	//	handler.HandleWebRTConnection(conn)
	//})

}
