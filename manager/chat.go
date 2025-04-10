package manager

import (
	"fmt"
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

	mu sync.Mutex
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

// Link 添加设备
func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
	c.mu.Lock()
	defer c.mu.Unlock()

	fmt.Printf("\033[32mmanager chat Link: %v accountID: %d sid: %s\033[0m\n", time.Now(), accountID, s.ID())
	c.sID.Store(s.ID(), accountID)

	// 清理旧连接（新增逻辑）
	if cm, ok := c.m.Load(accountID); ok {
		// 遍历并关闭所有旧连接
		cm.(*ConnMap).m.Range(func(key, value any) bool {
			oldConn := value.(*ActiveConn)
			oldConn.s.Close() // 关闭旧连接的 Socket.IO 连接
			cm.(*ConnMap).m.Delete(key)
			fmt.Printf("\033[33m[Close Old Connection] accountID: %d sid: %s\033[0m\n", accountID, key.(string))
			return true
		})
	}

	// 创建或更新 ConnMap
	cm, ok := c.m.Load(accountID)
	if !ok {
		newCm := &ConnMap{}
		activeConn := &ActiveConn{
			s:          s,
			activeTime: time.Now(),
		}
		newCm.m.Store(s.ID(), activeConn)
		c.m.Store(accountID, newCm)
		fmt.Printf("\033[32m[New ConnMap] accountID: %d sid: %s\033[0m\n", accountID, s.ID())
		return
	}
	activeConn := &ActiveConn{
		s:          s,
		activeTime: time.Now(),
	}
	cm.(*ConnMap).m.Store(s.ID(), activeConn)
	fmt.Printf("\033[32m[Back Connection] accountID: %d sid: %s\033[0m\n", accountID, s.ID())
}

//func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
//	c.mu.Lock()         // 加锁
//	defer c.mu.Unlock() // 确保解锁
//
//	fmt.Printf("\033[32mmanager chat Link: %v accountID: %d sid: %s\033[0m\n", time.Now(), accountID, s.ID())
//	c.sID.Store(s.ID(), accountID)
//
//	// 以下操作在锁保护下，保证原子性
//	cm, ok := c.m.Load(accountID)
//	if !ok {
//		// 没有 ConnMap，初始化并存储
//		newCm := &ConnMap{}
//		activeConn := &ActiveConn{
//			s:          s,
//			activeTime: time.Now(),
//		}
//		newCm.m.Store(s.ID(), activeConn)
//		c.m.Store(accountID, newCm)
//		fmt.Printf("\033[32m[New ConnMap] accountID: %d sid: %s\033[0m\n", accountID, s.ID())
//		return
//	}
//
//	// 已有 ConnMap，直接添加连接
//	activeConn := &ActiveConn{
//		s:          s,
//		activeTime: time.Now(),
//	}
//	cm.(*ConnMap).m.Store(s.ID(), activeConn)
//	fmt.Printf("\033[32m[Back Connection] accountID: %d sid: %s\033[0m\n", accountID, s.ID())
//}

//func (c *ChatMap) Link(s socketio.Conn, accountID int64) {
//	fmt.Println("\032[32mmanager chat Link:", time.Now(), "accountID:", accountID, s.ID(), "\032")
//	c.sID.Store(s.ID(), accountID) // 存入 SID 和 accountID 的对应关系
//	cm, ok := c.m.Load(accountID)
//	if !ok { // 没有找到对应的 ConnMap 对象，创建一个新的
//		cm := &ConnMap{}
//		activeConn := &ActiveConn{}
//		activeConn.s = s
//		activeConn.activeTime = time.Now()
//		cm.m.Store(s.ID(), activeConn) // 将新连接存储在 ConnMap 中
//		c.m.Store(accountID, cm)       // 将 ConnMap 存储在 c.m 中，以 accountID 为键
//		fmt.Println("\033[32222mmanager chat Link:", time.Now(), "accountID:", accountID, s.ID(), activeConn, "\033")
//		return
//	}
//	activeConn := &ActiveConn{}
//	activeConn.s = s
//	activeConn.activeTime = time.Now()
//	cm.(*ConnMap).m.Store(s.ID(), activeConn) // 将新的连接存储在 ConnMap 对象中
//	//fmt.Println("\033[32mmanager chat Link:", time.Now(), "accountID:", accountID, s.ID(), activeConn, "\033")
//}

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
		fmt.Println("\033[31mmanager chat Send:", time.Now(), "event:", event, "args:", args, "your_args\033[0m")
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

// SendMany 给指定多个账号的全部设备推送消息
// 参数：账号列表，事件名，要发送的数据
func (c *ChatMap) SendMany(accountIDs []int64, event string, args ...interface{}) {
	for _, accountID := range accountIDs {
		cm, ok := c.m.Load(accountID)
		if !ok { // 不存在该 accountID
			return
		}
		cm.(*ConnMap).m.Range(func(key, value interface{}) bool { // 遍历所有键值对
			activeConn := value.(*ActiveConn)
			activeConn.activeTime = time.Now() // 每次有消息发送，就重新计时
			activeConn.s.Emit(event, args...)  // 向指定客户端发送信息
			return true
		})
	}
}
