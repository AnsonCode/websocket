/**
* Created by GoLand.
* User: link1st
* Date: 2019-08-03
* Time: 16:43
 */

package grpcserver

import (
	"context"
	"fmt"
	"log"
	"net"

	"github.com/golang/protobuf/proto"
	"github.com/spf13/viper"
	"google.golang.org/grpc"
	"turing.com/push/common"
	"turing.com/push/grpc/protobuf"
	"turing.com/push/socket"
	// "turing.com/push/servers/websocket"
)

type server struct {
}

func setErr(rsp proto.Message, code uint32, message string) {

	message = common.GetErrorMessage(code, message)
	switch v := rsp.(type) {
	case *protobuf.QueryUsersOnlineRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.BroadcastToRoomRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.SendMsgAllRsp:
		v.RetCode = code
		v.ErrMsg = message
	case *protobuf.GetUserListRsp:
		v.RetCode = code
		v.ErrMsg = message
	default:

	}

}

// 查询用户是否在线
func (s *server) QueryUsersOnline(c context.Context, req *protobuf.QueryUsersOnlineReq) (rsp *protobuf.QueryUsersOnlineRsp, err error) {

	fmt.Println("grpc_request 查询用户是否在线", req.String())

	rsp = &protobuf.QueryUsersOnlineRsp{}

	// online := websocket.CheckUserOnline(req.GetAppId(), req.GetUserId())

	setErr(req, common.OK, "")
	rsp.Online = true

	return rsp, nil
}

// 给本机用户发消息
func (s *server) BroadcastToRoom(c context.Context, req *protobuf.BroadcastToRoomReq) (rsp *protobuf.BroadcastToRoomRsp, err error) {

	fmt.Println("grpc_request 给本机用户发消息", req.String())
	rsp = &protobuf.BroadcastToRoomRsp{}

	if ok := socket.SocketIOServer.BroadcastToRoom(req.Namespace, req.Room, req.Event, req.Msg); ok {
		setErr(rsp, common.OK, "")
	} else {
		setErr(rsp, common.ServerError, "BroadcastToRoom error")
	}

	fmt.Println("grpc_response 给本机用户发消息", rsp.String())
	return
}

// 获取本机用户列表
func (s *server) GetUserList(c context.Context, req *protobuf.GetUserListReq) (rsp *protobuf.GetUserListRsp, err error) {

	fmt.Println("grpc_request 获取本机用户列表", req.String())

	// appId := req.GetAppId()
	rsp = &protobuf.GetUserListRsp{}

	// 本机
	// userList := websocket.GetUserList(appId)

	setErr(rsp, common.OK, "")
	// rsp.UserId = userList

	fmt.Println("grpc_response 获取本机用户列表:", rsp.String())

	return
}

// rpc server
// link::https://github.com/grpc/grpc-go/blob/master/examples/helloworld/greeter_server/main.go
func Init() {

	rpcPort := viper.GetString("app.rpcPort")
	fmt.Println("rpc server 启动", rpcPort)

	lis, err := net.Listen("tcp", ":"+rpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}
	s := grpc.NewServer()
	protobuf.RegisterAccServerServer(s, &server{})
	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
