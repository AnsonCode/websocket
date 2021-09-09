package socket

import (
	"fmt"
	"log"
	"net/http"

	"github.com/gin-gonic/gin"
	socketio "github.com/googollee/go-socket.io"
)

// "turing.com/push/grpc/grpcclient"

var SocketIOServer *socketio.Server

func InitSocketManager() {
	if SocketIOServer != nil {
		return
	}
	server := socketio.NewServer(nil)
	// server.to
	// _, err := server.Adapter(&socketio.RedisAdapterOptions{
	// 	Host:   "127.0.0.1",
	// 	Port:   "6379",
	// 	Prefix: "socket.io",
	// })
	// address     = "127.0.0.1:6379"
	// Addr:         viper.GetString("redis.addr"),
	// 	Password:     viper.GetString("redis.password"),
	// 	DB:           viper.GetInt("redis.DB"),
	// 	PoolSize:     viper.GetInt("redis.poolSize"),
	// 	MinIdleConns: viper.GetInt("redis.minIdleConns"),

	// if err != nil {
	// 	log.Fatal("error:", err)
	// }

	server.OnConnect("/", func(s socketio.Conn) error {
		// s.Join("airswitchRoom")
		// s.JoinRoom("airswitchRoom")
		// either with send()
		s.Emit("greetings", "Hey")
		// server.BroadcastToNamespace("", "test", "sss")
		server.BroadcastToNamespace("ant", "test2", "sss2")
		server.BroadcastToRoom("", "room", "test", "sss")
		// socket.emit('greetings', 'Hey!', { 'ms': 'jane' }, Buffer.from([4, 3, 3, 1]));

		user := s.URL().User
		fmt.Printf("user,%v", user)
		//TODO:这里需要 ？
		s.Join("/sys/notice") //所有的用户都能收到
		// s.URL().RawQuery
		s.SetContext("") //这里可以附加别的参数
		fmt.Println("connected:", s.ID())
		return nil
	})
	server.OnEvent("/", "join", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		fmt.Println("join:", msg)
		s.Join(msg)
		s.Emit("reply", "join success "+msg)
		// return "join " + msg
	})
	server.OnEvent("/", "leave", func(s socketio.Conn, msg string) {
		s.SetContext(msg)
		fmt.Println("leave:", msg)
		s.Leave(msg)
		s.Emit("reply", "leave success "+msg)
		// return "join " + msg
	})
	server.OnEvent("/", "leaveAll", func(s socketio.Conn) {
		fmt.Println("leaveAll:")
		s.LeaveAll()
		s.Emit("reply", "leaveAll success ")
		// return "join " + msg
	})

	server.OnEvent("/sys", "notice", func(s socketio.Conn, msg string) {
		fmt.Println("notice:", msg)
		s.Emit("reply", "have "+msg)
	})

	server.OnEvent("/sys", "bye", func(s socketio.Conn) string {
		last := s.Context().(string)
		s.Emit("bye", last)
		s.Close()
		return last
	})

	server.OnEvent("/", "chat", func(s socketio.Conn, msg string) string {
		s.SetContext(msg)
		fmt.Println("recv:", msg)

		return "recv " + msg
	})

	server.OnError("/", func(s socketio.Conn, e error) {
		fmt.Println("meet error:", e)
	})

	server.OnDisconnect("/", func(s socketio.Conn, reason string) {
		fmt.Println("closed", reason)
	})

	// server.BroadcastToRoom()
	// server.BroadcastToNamespace()
	// server.JoinRoom()
	// server.BroadcastToRoom()
	// server.BroadcastToRoom("", "bcast", "event:name", msg)

	go func() {
		if err := server.Serve(); err != nil {
			log.Fatalf("socketio listen error: %s\n", err)
		}
	}()
	defer server.Close()
	SocketIOServer = server

	router := gin.New()
	router.Use(GinMiddleware("*"))
	router.GET("/socket.io/*any", gin.WrapH(server))
	router.POST("/socket.io/*any", gin.WrapH(server))
	router.StaticFS("/public", http.Dir("./views"))

	if err := router.Run(":8000"); err != nil {
		log.Fatal("failed run app: ", err)
	}

	// http.Handle("/socket.io/", server)
	// http.Handle("/", http.FileServer(http.Dir("./views")))
	// log.Println("Serving at localhost:8000...")
	// log.Fatal(http.ListenAndServe(":8000", nil))
}

// func GetClient() (c *socketio.Server) {

// 	return SocketIOServer
// }

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
