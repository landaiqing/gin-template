package websocket_api

import (
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/lxzan/gws"
	"net/http"
	"schisandra-cloud-album/global"
	"time"
)

const (
	PingInterval         = 5 * time.Second  // 客户端心跳间隔
	HeartbeatWaitTimeout = 10 * time.Second // 心跳等待超时时间
)

var Handler = NewWebSocket()

func (WebsocketAPI) NewGWSServer(c *gin.Context) {

	upgrader := gws.NewUpgrader(Handler, &gws.ServerOption{
		HandshakeTimeout: 5 * time.Second, // 握手超时时间
		ReadBufferSize:   1024,            // 读缓冲区大小
		ParallelEnabled:  true,            // 开启并行消息处理
		Recovery:         gws.Recovery,    // 开启异常恢复
		CheckUtf8Enabled: true,            // 开启UTF8校验
		PermessageDeflate: gws.PermessageDeflate{
			Enabled: true, // 开启压缩
		},
		Authorize: func(r *http.Request, session gws.SessionStorage) bool {
			var clientId = r.URL.Query().Get("client_id")
			if clientId == "" {
				return false
			}
			session.Store("client_id", clientId)
			return true
		},
	})
	socket, err := upgrader.Upgrade(c.Writer, c.Request)
	if err != nil {
		return
	}
	go func() {
		socket.ReadLoop() // 此处阻塞会使请求上下文不能顺利被GC
	}()
}
func MustLoad[T any](session gws.SessionStorage, key string) (v T) {
	if value, exist := session.Load(key); exist {
		v = value.(T)
	}
	return
}

func NewWebSocket() *WebSocket {
	return &WebSocket{
		sessions: gws.NewConcurrentMap[string, *gws.Conn](64, 128),
	}
}

type WebSocket struct {
	gws.BuiltinEventHandler
	sessions *gws.ConcurrentMap[string, *gws.Conn] // 使用内置的ConcurrentMap存储连接, 可以减少锁冲突
}

func (c *WebSocket) OnOpen(socket *gws.Conn) {
	name := MustLoad[string](socket.Session(), "client_id")
	if conn, ok := c.sessions.Load(name); ok {
		conn.WriteClose(1000, []byte("connection is replaced"))
	}
	c.sessions.Store(name, socket)
	global.LOG.Printf("%s connected\n", name)
}

func (c *WebSocket) OnClose(socket *gws.Conn, err error) {
	name := MustLoad[string](socket.Session(), "client_id")
	sharding := c.sessions.GetSharding(name)
	c.sessions.Delete(name)
	sharding.Lock()
	defer sharding.Unlock()

	global.LOG.Printf("onerror, name=%s, msg=%s\n", name, err.Error())
}

func (c *WebSocket) OnPing(socket *gws.Conn, payload []byte) {
	_ = socket.SetDeadline(time.Now().Add(PingInterval + HeartbeatWaitTimeout))
	_ = socket.WritePong(payload)
}

func (c *WebSocket) OnPong(socket *gws.Conn, payload []byte) {}

func (c *WebSocket) OnMessage(socket *gws.Conn, message *gws.Message) {
	defer message.Close()
	name := MustLoad[string](socket.Session(), "client_id")
	if conn, ok := c.sessions.Load(name); ok {
		_ = conn.WriteMessage(gws.OpcodeText, message.Bytes())
	}
}

// SendMessageToClient 向指定客户端发送消息
func (c *WebSocket) SendMessageToClient(clientId string, message []byte) error {
	conn, ok := c.sessions.Load(clientId)
	if ok {
		return conn.WriteMessage(gws.OpcodeText, message)
	}
	return fmt.Errorf("client %s not found", clientId)
}
