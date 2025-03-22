package consumer

import (
	"ChatRoom001/global"
	"ChatRoom001/model/chat"
	"ChatRoom001/model/reply"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
	"os"
	"os/signal"
	"strings"
	"syscall"
	"time"
)

// StartConsumer 启动 Kafka 消费者
func StartConsumer(accountID int64) {
	// 生成合法的 Kafka Topic 名称（替换 RocketMQ 中的冒号为下划线）
	topic := sanitizeTopicName(fmt.Sprintf("accountID_%d", accountID))

	// 创建消费者配置
	config := sarama.NewConfig()
	config.Consumer.Return.Errors = true
	config.Consumer.Offsets.AutoCommit.Enable = true              // 启用自动提交
	config.Consumer.Offsets.AutoCommit.Interval = 1 * time.Second // 自动提交间隔
	config.Consumer.Offsets.Initial = sarama.OffsetOldest         // 从最早的消息开始消费

	// 创建消费者客户端
	consumerClient, err := sarama.NewConsumer(
		[]string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)},
		config,
	)
	if err != nil {
		fmt.Printf("创建 Kafka 消费者客户端失败: %v\n", err)
		return
	}
	defer consumerClient.Close()

	// 获取主题分区列表
	partitions, err := consumerClient.Partitions(topic)
	if err != nil {
		fmt.Printf("获取主题分区失败: %v\n", err)
		return
	}

	// 创建退出信号通道
	sigchan := make(chan os.Signal, 1)
	signal.Notify(sigchan, syscall.SIGINT, syscall.SIGTERM)

	// 为每个分区创建消费者
	consumers := make([]sarama.PartitionConsumer, 0, len(partitions))
	for _, partition := range partitions {
		pc, err := consumerClient.ConsumePartition(topic, partition, sarama.OffsetOldest)
		if err != nil {
			fmt.Printf("创建分区消费者失败: %v\n", err)
			continue
		}
		defer pc.AsyncClose()
		consumers = append(consumers, pc)

		// 启动消费协程
		go func(pc sarama.PartitionConsumer) {
			for {
				select {
				case msg := <-pc.Messages():
					processMessage(accountID, msg)
				case err := <-pc.Errors():
					fmt.Printf("消费错误: %v\n", err)
				case <-sigchan:
					return
				}
			}
		}(pc)
	}

	fmt.Printf("消费者已启动，正在监听主题: %s\n", topic)
	<-sigchan // 阻塞直到接收到退出信号
	fmt.Println("正在关闭消费者...")
}

// 处理消息
func processMessage(accountID int64, msg *sarama.ConsumerMessage) {
	var message reply.ParamMsgInfoWithRly
	if err := json.Unmarshal(msg.Value, &message); err != nil {
		fmt.Printf("消息反序列化失败: %v\n", err)
		return
	}

	// 发送到全局 ChatMap
	global.ChatMap.Send(accountID, chat.ClientSendMsg, message)
	fmt.Printf("收到消息 [分区:%d 偏移量:%d] 内容: %+v\n",
		msg.Partition, msg.Offset, message)
}

// 清理非法字符生成合法 Topic 名称
func sanitizeTopicName(name string) string {
	return strings.ReplaceAll(name, ":", "_")
}
