// package producer
//
// import (
//
//	"ChatRoom001/global"
//	"ChatRoom001/model/reply"
//	"encoding/json"
//	"fmt"
//	"github.com/IBM/sarama"
//
// )
//
// // SendMsgToKafka 通过 Kafka 发送消息
//
//	func SendMsgToKafka(mID int64, msg reply.ParamMsgInfoWithRly) {
//		// 构建 Kafka Broker 地址
//		brokers := []string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)}
//
//		// 创建生产者配置
//		config := sarama.NewConfig()
//		config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
//		config.Producer.Retry.Max = 5                    // 重试次数
//		config.Producer.Return.Successes = true          // 成功交付的消息将在 success channel 返回
//
//		// 创建同步生产者
//		producer, err := sarama.NewSyncProducer(brokers, config)
//		if err != nil {
//			panic(fmt.Sprintf("创建 Kafka 生产者失败: %s", err))
//		}
//		defer producer.Close()
//
//		// 构建合法的 Kafka Topic 名称（替换 RocketMQ 中的冒号为下划线）
//		topic := fmt.Sprintf("accountID_%d", mID)
//
//		// 序列化消息内容
//		sendMsg, err := json.Marshal(msg)
//		if err != nil {
//			fmt.Println("序列化消息失败", err)
//			return
//		}
//
//		// 构造 Kafka 消息对象
//		kafkaMsg := &sarama.ProducerMessage{
//			Topic: topic,
//			Value: sarama.ByteEncoder(sendMsg),
//			// 可根据需要设置消息 Key
//			// Key: sarama.StringEncoder(fmt.Sprintf("%d", mID)),
//		}
//
//		// 发送消息到 Kafka
//		partition, offset, err := producer.SendMessage(kafkaMsg)
//		if err != nil {
//			fmt.Println("消息发送失败:", err)
//		} else {
//			fmt.Printf("消息发送成功 [分区:%d] [偏移量:%d]\n", partition, offset)
//		}
//	}
package producer

import (
	"ChatRoom001/global"
	"ChatRoom001/model/reply"
	"encoding/json"
	"fmt"
	"github.com/IBM/sarama"
)

// SendMsgToKafka 通过 Kafka 发送消息
func SendMsgToKafka(mID int64, msg reply.ParamMsgInfoWithRly) {
	// 构建 Kafka Broker 地址
	brokers := []string{fmt.Sprintf("%s:%d", global.PrivateSetting.RocketMQ.Addr, global.PrivateSetting.RocketMQ.Port)}

	// 创建生产者配置
	config := sarama.NewConfig()
	config.Producer.RequiredAcks = sarama.WaitForAll // 等待所有副本确认
	config.Producer.Retry.Max = 5                    // 重试次数
	config.Producer.Return.Successes = true          // 成功交付的消息将在 success channel 返回

	// 创建同步生产者
	producer, err := sarama.NewSyncProducer(brokers, config)
	if err != nil {
		panic(fmt.Sprintf("创建 Kafka 生产者失败: %s", err))
	}
	defer producer.Close()

	// 构建合法的 Kafka Topic 名称（替换 RocketMQ 中的冒号为下划线）
	topic := fmt.Sprintf("accountID_%d", mID)

	// 序列化消息内容
	sendMsg, err := json.Marshal(msg)
	if err != nil {
		fmt.Println("序列化消息失败", err)
		return
	}

	// 构造 Kafka 消息对象
	kafkaMsg := &sarama.ProducerMessage{
		Topic: topic,
		Value: sarama.ByteEncoder(sendMsg),
		// 可根据需要设置消息 Key
		// Key: sarama.StringEncoder(fmt.Sprintf("%d", mID)),
	}

	// 发送消息到 Kafka
	partition, offset, err := producer.SendMessage(kafkaMsg)
	if err != nil {
		fmt.Println("消息发送失败:", err)
	} else {
		fmt.Printf("消息发送成功 [分区:%d] [偏移量:%d]\n", partition, offset)
	}
}
