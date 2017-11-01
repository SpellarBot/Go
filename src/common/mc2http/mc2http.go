package mc2http

import (
	"fmt"
	"runtime"
	"log"
	"io/ioutil"

	"github.com/hongst/rend/common"
	"github.com/hongst/rend/handlers"
	"github.com/hongst/rend/orcas"
	"github.com/hongst/rend/server"
)

type MC2HTTPServer struct {
	ListenPort int
	Callback   CallbackHookFunc
	Logger     func(string)
}

func (m *MC2HTTPServer) Init() {
	if m.Logger == nil {
		m.Logger = func(msg string) {
			fmt.Printf(msg)
		}
	}

	log.SetOutput(ioutil.Discard) // 去除不必要的输出，需要输出的话，统一放到log中
}

func (m *MC2HTTPServer) Run() {
	if m.ListenPort == 0 {
		m.Logger("[MC2HTTPServer] Error! Must Set ListenPorts!")
		return
	}
	largs := server.ListenArgs{
		Type: server.ListenTCP,
		Port: m.ListenPort,
	}
	go server.ListenAndServe(
		largs,
		server.Default,
		orcas.L1Only,
		mc2httpHandlerNew(m.Callback),
		handlers.NilHandler,
	)
	runtime.Goexit()
}

type CallbackHookFunc func(request string) (string, error)

type mc2httpHandler struct {
	Callback CallbackHookFunc
}

func mc2httpHandlerNew(callback CallbackHookFunc) handlers.HandlerConst {
	singleton := &mc2httpHandler{
		Callback: callback,
	}
	return func() (handlers.Handler, error) {
		return singleton, nil
	}
}

func (h *mc2httpHandler) Get(cmd common.GetRequest) (<-chan common.GetResponse, <-chan error) {
	dataOut := make(chan common.GetResponse)
	errorOut := make(chan error)
	go realHandleGet(h, cmd, dataOut, errorOut)
	return dataOut, errorOut
}

func realHandleGet(h *mc2httpHandler, cmd common.GetRequest, dataOut chan common.GetResponse, errorOut chan error) {
	defer close(errorOut)
	defer close(dataOut)

	for idx, key := range cmd.Keys {

		skey := string(key)

		// PHP扩展会使用stats命令来尝试连接
		if skey == "stats" {
			dataOut <- common.GetResponse{
				Miss:   false,
				Quiet:  cmd.Quiet[idx],
				Opaque: cmd.Opaques[idx],
				Flags:  0,
				Key:    []byte("hello"),
				Data:   []byte("world"),
			}
			break
		}

		// 回调
		data, err := h.Callback(skey)
		if err != nil {
			errorOut <- common.ErrInternal
		}

		// 返回
		dataOut <- common.GetResponse{
			Miss:   false,
			Quiet:  cmd.Quiet[idx],
			Opaque: cmd.Opaques[idx],
			Flags:  0,
			Key:    []byte("zhh"), // 替换掉RESP中的key，因为超过250个字节连接会被reset
			Data:   []byte(data),
		}

	}
}

func (h *mc2httpHandler) Close() error {
	return nil
}

func (h *mc2httpHandler) Set(cmd common.SetRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) Delete(cmd common.DeleteRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) Add(cmd common.SetRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) Replace(cmd common.SetRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) Append(cmd common.SetRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) Prepend(cmd common.SetRequest) error {
	return common.ErrUnknownCmd
}

func (h *mc2httpHandler) GetE(cmd common.GetRequest) (<-chan common.GetEResponse, <-chan error) {
	errchan := make(chan error, 1)
	errchan <- common.ErrUnknownCmd
	return nil, errchan
}

func (h *mc2httpHandler) GAT(cmd common.GATRequest) (common.GetResponse, error) {
	return common.GetResponse{}, common.ErrUnknownCmd
}

func (h *mc2httpHandler) Touch(cmd common.TouchRequest) error {
	return common.ErrUnknownCmd
}