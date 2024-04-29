package QNQ

import (
	"context"
	"errors"
	"flag"
	"fmt"
	"google.golang.org/grpc"
	"google.golang.org/grpc/credentials/insecure"
	"google.golang.org/grpc/reflection"
	"log/slog"
	"net"
	"reflect"
	"strings"
)

var timeFormat = "2006/01/02 15:04:05"

var port = flag.Int("port", 6615, "rpc port")

var needAgentMap = make(map[string]bool)

var loggerIntercept = func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (any, error) {
	slog.Info("RPC Call ", "method", info.FullMethod)
	defer slog.Info("RPC Done ", "method", info.FullMethod)
	var res any
	var err error
	b, ip := checkAgent(req, info)
	if b && ip != localhost {
		res, err = agentRPC(ip, info.FullMethod, ctx, req)
	} else {
		res, err = handler(ctx, req)
	}
	return res, err
}

func agentRPC(ip, fullMethod string, ctx context.Context, req any) (any, error) {
	slog.Info("agent intercept", "target ip", ip)
	var res any
	var err error
	info := strings.Split(fullMethod[1:], "/")
	if len(info) != 2 {
		err = errors.New("fullMethod err")
		return res, err
	}
	cli, err := getRPClient(ip, info[0])
	if err != nil {
		return res, err
	}

	reflectRes, err := reflectMethod(cli, info[1], ctx, req)
	if err != nil {
		return res, err
	}
	if len(reflectRes) != 2 {
		err = errors.New("reflect res err")
		return res, err
	}

	defer func() {
		if r := recover(); r != nil {
			slog.Error("cast res err", "info", r)
		}
	}()

	res = reflectRes[0]
	if reflect.ValueOf(reflectRes[1]).Kind() != reflect.Invalid {
		err = reflectRes[1].(error)
	}
	return res, err

}

func interceptChain(intercepts ...grpc.UnaryServerInterceptor) grpc.UnaryServerInterceptor {
	l := len(intercepts)
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		chain := func(currentInter grpc.UnaryServerInterceptor, currentHandler grpc.UnaryHandler) grpc.UnaryHandler {
			return func(ctx context.Context, req any) (any, error) {
				return currentInter(ctx, req, info, currentHandler)
			}
		}
		chainHandler := handler
		for i := l - 1; i >= 0; i-- {
			chainHandler = chain(intercepts[i], chainHandler)
		}
		return chainHandler(ctx, req)
	}
}

func Start() {
	//_, err := credentials.NewServerTLSFromFile("./ca/server_cert.pem", "./ca/server_key.pem")
	//if err != nil {
	//	slog.Error("create tls listen err: ", err.Error())
	//}
	flag.Parse()
	lis, err := net.Listen("tcp", fmt.Sprintf(":%d", *port))
	if err != nil {
		slog.Error("listen rpc port err", "err", err.Error())
	}
	server := grpc.NewServer(grpc.UnaryInterceptor(loggerIntercept))

	// local
	RegisterLocalSyncServer(server, &LocalSync{})
	// remote
	RegisterRemoteSyncServer(server, &RemoteSync{})
	// sys probe
	RegisterSysProbeServer(server, &SysProbe{})
	// task center
	RegisterTaskCenterServer(server, &TaskCenter{})

	reflection.Register(server)

	listenGRPC(server, &lis)

}

func listenGRPC(s *grpc.Server, lis *net.Listener) {
	go func() {
		for k, _ := range s.GetServiceInfo() {
			slog.Info("start listen : ", "server name", k)
		}
		err := s.Serve(*lis)
		if err != nil {
			slog.Error("start listen rpc err: " + err.Error())
		}
	}()
}

/*
*
key0: ip, key1: client
*/
var clientMap = make(map[string]RPClient)

type RPClient struct {
	cliMap map[string]any
	conn   *grpc.ClientConn
	Des    string
}

func addRPClient(ip string, port int) error {
	flag.Parse()
	addr := flag.String("addr", fmt.Sprintf("%s:%d", ip, port), "the address to connect to")
	conn, err := grpc.Dial(*addr, grpc.WithTransportCredentials(insecure.NewCredentials()))
	if err != nil {
		slog.Error("add qnq target err", "err", err.Error())
		return err
	}
	m := make(map[string]any)
	m["LocalSync"] = NewLocalSyncClient(conn)
	m["RemoteSync"] = NewRemoteSyncClient(conn)
	m["SysProbe"] = NewSysProbeClient(conn)
	m["TaskCenter"] = NewTaskCenterClient(conn)
	clientMap[ip] = RPClient{
		cliMap: m,
		conn:   conn,
	}
	slog.Info("add RPC client : " + fmt.Sprintf("%s:%d", ip, port))
	return nil
}

func getRPClient(ip, serverName string) (any, error) {
	if _, ok := clientMap[ip]; ok {
		if server, ok := clientMap[ip].cliMap[serverName]; ok {
			return server, nil
		} else {
			return nil, errors.New("server not register")
		}
	} else {
		return nil, errors.New("target not register")
	}
}

func deleteRPClient(ip string) error {
	var err error
	if v, ok := clientMap[ip]; ok {
		err = v.conn.Close()
		delete(clientMap, ip)
		slog.Info("delete rpc client", "ip", ip)
	} else {
		err = errors.New("rpc client no connected, ip: " + ip)
	}
	return err
}

func loadErrToResult(res *Result, err error) {
	if err == nil || res == nil {
		return
	}
	res.Code = ERR_CODE
	res.Message = err.Error()
}

func checkAgent(req any, info *grpc.UnaryServerInfo) (bool, string) {
	var ip string
	method := info.FullMethod
	b, val := hasField(req, "TargetAddr", localhost)
	ip = fmt.Sprint(val)
	if _, ok := needAgentMap[method]; !ok {
		needAgentMap[method] = b
	}
	return needAgentMap[method], ip
}
