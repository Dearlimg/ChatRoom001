package manager

import (
	socketio "github.com/googollee/go-socket.io"
	"sync"
	"time"
)

/*
每个账户各个客户端消息的发送
*/

const DefaultClientTimeout = time.Minute * 20

type ChatMap struct {
	// （GO 内置的 map 不是并发安全的，sync.Map 是并发安全的）
	m   sync.Map // k: accountID v: ConnMap（说明 accountID 可以不止有一个客户端设备）
	sID sync.Map // k: sID v: accountID，用于快速查找 某个连接 属于哪个 账户
}

type ConnMap struct {
	m sync.Map // k: sID v: ActiveConn（即每个 sID 对应一个活跃连接）
}

type ActiveConn struct {
	s          socketio.Conn // 连接对象，用于实际的连接操作（如发送/接收数据）
	activeTime time.Time     // 最后活动时间，用于判断连接是否超时
}

func NewChatMap() *ChatMap {
	return &ChatMap{m: sync.Map{}}
}

func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
	c.m.Store(s.ID(), accountID)
	cm, ok := c.m.Load(accountID)
	if !ok {
		cm := &ConnMap{}
		activeConn := &ActiveConn{}
		activeConn.s = s
		activeConn.activeTime = time.Now()
		cm.m.Store(s.ID(), activeConn)
		c.m.Store(accountID, cm)
		return
	}
	activeConn := &ActiveConn{}
	activeConn.s = s
	activeConn.activeTime = time.Now()
	cm.(*ConnMap).m.Store(s.ID(), activeConn)
}

func (c *ChatMap) Leave(s socketio.Conn) {
	accountID, ok := c.sID.LoadAndDelete(s.ID())
	if !ok {
		return
	}
	cm, ok := c.m.Load(accountID)
	if !ok {
		return
	}
	cm.(*ConnMap).m.Delete(s.ID())
	length := 0
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		length++
		return true
	})
	if length == 0 {
		c.m.Delete(accountID)
	}
}

// Send 给指定账号的全部设备推送消息
func (c *ChatMap) Send(accountID int64, event string, args ...interface{}) {
	cm, ok := c.m.Load(accountID)
	if !ok {
		return
	}
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		activeConn := value.(*ActiveConn)
		activeConn.activeTime = time.Now()
		activeConn.s.Emit(event, args...)
		return true
	})
}

type EachFunc socketio.EachFunc

func (c *ChatMap) ForEach(accountID int64, f EachFunc) {
	cm, ok := c.m.Load(accountID)
	if !ok {
		return
	}
	cm.(*ConnMap).m.Range(func(key, value any) bool {
		f(value.(*ActiveConn).s)
		return true
	})
}

func (c *ChatMap) CheckIsOnConnection(accountID int64) bool {
	_, ok := c.m.Load(accountID)
	return ok
}

func (c *ChatMap) HasSID(sID string) bool {
	_, ok := c.sID.Load(sID)
	return ok
}
