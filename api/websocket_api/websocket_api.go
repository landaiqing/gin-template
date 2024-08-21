package websocket_api

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"log"
	"net/http"
	"schisandra-cloud-album/common/result"
	"schisandra-cloud-album/global"
	"sync"
	"time"
)

var (
	// 消息通道
	msg = make(map[string]chan interface{})
	// websocket客户端链接池
	client = make(map[string]*websocket.Conn)
	// 互斥锁，防止程序对统一资源同时进行读写
	mux sync.Mutex
)

// NewSocketClient 创建websocket服务
// @Summary 创建websocket服务(gorilla)
// @Description 创建websocket服务
// @Tags websocket
// @Router /api/ws/socket [get]
func (WebsocketAPI) NewSocketClient(context *gin.Context) {
	id := context.Query("client_id")
	global.LOG.Println(id + "websocket链接")
	// 升级为websocket长链接
	WsHandler(context.Writer, context.Request, id)
}

// DeleteClient api:/deleteClient接口处理函数
func (WebsocketAPI) DeleteClient(context *gin.Context) {
	id := context.Query("client_id")
	// 关闭websocket链接
	conn, exist := getClient(id)
	if exist {
		err := conn.Close()
		if err != nil {
			return
		}
		deleteClient(id)
	} else {
		result.FailWithMessage("客户端不存在", context)
	}
	// 关闭其消息通道
	_, exist = getMsgChannel(id)
	if exist {
		deletemsgChannel(id)
	}
}

// SendMessageData 发送消息接口处理函数
func SendMessageData(clientId string, data interface{}) bool {
	m, exist := getMsgChannel(clientId)
	if !exist {
		log.Println("未找到该客户端的消息通道")
		return false
	}
	// 向消息通道发送消息
	select {
	case m <- data:
		global.LOG.Println("发送消息给客户端：" + clientId)
		return true
	default:
		global.LOG.Println("消息通道已满，消息发送失败")
	}
	return false
}

// websocket Upgrader
var wsupgrader = websocket.Upgrader{
	ReadBufferSize:   1024,
	WriteBufferSize:  1024,
	HandshakeTimeout: 5 * time.Second,
	// 取消ws跨域校验
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

// WsHandler 处理ws请求
func WsHandler(w http.ResponseWriter, r *http.Request, id string) {
	var conn *websocket.Conn
	var err error
	var exist bool
	// 创建一个定时器用于服务端心跳
	pingTicker := time.NewTicker(time.Second * 10)
	conn, err = wsupgrader.Upgrade(w, r, nil)
	if err != nil {
		global.LOG.Println(err)
		return
	}
	// 把与客户端的链接添加到客户端链接池中
	addClient(id, conn)

	// 获取该客户端的消息通道
	m, exist := getMsgChannel(id)
	if !exist {
		m = make(chan interface{})
		addMsgChannel(id, m)
	}
	// 设置客户端关闭ws链接回调函数
	conn.SetCloseHandler(func(code int, text string) error {
		deleteClient(id)
		return nil
	})
	defer conn.Close()
	for {
		_, _, err := conn.ReadMessage()
		if err != nil {
			global.LOG.Error(err)
			return
		}
		select {
		case content, _ := <-m:
			// 从消息通道接收消息，然后推送给前端
			err = conn.WriteJSON(content)
			if err != nil {
				err = conn.Close()
				if err != nil {
					return
				}
				deleteClient(id)
				break
			}
		case <-pingTicker.C:
			// 服务端心跳:每20秒ping一次客户端，查看其是否在线
			err := conn.SetWriteDeadline(time.Now().Add(time.Second * 20))
			if err != nil {
				return
			}
			err = conn.WriteMessage(websocket.PingMessage, []byte{})
			if err != nil {
				log.Println("send pong err:", err)
				err := conn.Close()
				if err != nil {
					return
				}
				deleteClient(id)
				return
			}
		}
	}

}

// 将客户端添加到客户端链接池
func addClient(id string, conn *websocket.Conn) {
	mux.Lock()
	client[id] = conn
	mux.Unlock()
}

// 获取指定客户端链接
func getClient(id string) (conn *websocket.Conn, exist bool) {
	mux.Lock()
	conn, exist = client[id]
	mux.Unlock()
	return
}

// 删除客户端链接
func deleteClient(id string) {
	mux.Lock()
	delete(client, id)
	log.Println(id + "websocket退出")
	mux.Unlock()
}

// 添加用户消息通道
func addMsgChannel(id string, m chan interface{}) {
	mux.Lock()
	msg[id] = m
	mux.Unlock()
}

// 获取指定用户消息通道
func getMsgChannel(id string) (m chan interface{}, exist bool) {
	mux.Lock()
	defer mux.Unlock()
	m, exist = msg[id]
	return
}

// 删除指定消息通道
func deletemsgChannel(id string) {
	mux.Lock()
	if m, ok := msg[id]; ok {
		close(m)
		delete(msg, id)
	}
	mux.Unlock()
}
