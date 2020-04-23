/*
 *  Copyright (c) 2018, https://github.com/nebulaim
 *  All rights reserved.
 *
 * Licensed under the Apache License, Version 2.0 (the "License");
 * you may not use this file except in compliance with the License.
 * You may obtain a copy of the License at
 *
 *   http://www.apache.org/licenses/LICENSE-2.0
 *
 * Unless required by applicable law or agreed to in writing, software
 * distributed under the License is distributed on an "AS IS" BASIS,
 * WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
 * See the License for the specific language governing permissions and
 * limitations under the License.
 */

package nsq_client

/*
import (
	"errors"
	"time"
	"github.com/nsqio/go-nsq"
	"github.com/golang/glog"
)


// 使用者可以使用快速的配置进行事件的消费
// 配置
// 生产者
// [producer]
// nsqd地址
// nsqd=127.0.0.1:4151
// 消费者
// [consumer]
// nsqd连接地址
// nsqd=127.0.0.1:4151
// nsqlookupd连接地址
// nsqlookupd=127.0.0.1:4161,127.0.0.2:4161,127.0.0.3:4161
// max_flight=100
// concurrent=20
// channel=chan1
// max_retry=5

type ConsumerConfig struct {
	NsqdAddr string
	NsqlookupdAddr string
	MaxFlight int
	Concurrent int
	Channel int
	MaxRetry int
}

// consumer 消费者结构体
type consumer struct {
	isInit      bool
	// Debug       bool
	channelName string
	concurrent  int
	maxInFlight int
	maxAttempt  uint16
	//addr 连接地址
	nsqdAddr       []string
	nsqLookupdAddr []string
	// 各个topic的worker
	topics map[string]*topicInfo
}

// topicInfo topic 信息结构体
type topicInfo struct {
	topic         string
	maxInFlight   int
	concurrentNum int
	config        *nsq.Config
	handler       nsq.HandlerFunc
	failHandler   FailMessageFunc
	consumer      *nsq.Consumer
}

// 失败消息处理函数类型
type FailMessageFunc func(message FailMessage) (err error)

func (f FailMessageFunc) HandleFailMessage(message FailMessage) (err error) {
	err = f(message)
	return
}

// 失败消息处理接口,继承了该接口的接口都会调用该接口
type FailMessageHandler interface {
	HandleFailMessage(message FailMessage) (err error)
}

type FailMessage struct {
	Body      []byte
	Attempt   uint16
	Timestamp int64
	MessageID string
	FailMsg   string
}

// Connect 连接
func (t *topicInfo) Connect(channelName string, nsqdAddr []string, nsqlookupdAddr []string) {
	if len(nsqdAddr) == 0 && len(nsqlookupdAddr) == 0 {
		glog.Warningln("nsqd和nsqlookupd地址皆为空，跳过连接,topic:", t.topic)
		return
	}
	var (
		retryNum     = 0
		sleepSeconds = 0
		err          error
	)
	t.consumer, err = nsq.NewConsumer(t.topic, channelName, t.config)
	if err != nil {
		glog.Errorf("新建nsq consumer失败，err:%s,topic:%s,channel:%s", err.Error(), t.topic, channelName)
		return
	}
	t.consumer.ChangeMaxInFlight(t.maxInFlight)
	t.consumer.AddConcurrentHandlers(nsq.Handler(t.handler), t.concurrentNum)
	// 不断进行重连，直到连接成功
	for {
		if len(nsqlookupdAddr) > 0 {
			if len(nsqlookupdAddr) == 1 {
				err = t.consumer.ConnectToNSQLookupd(nsqlookupdAddr[0])
			} else {
				err = t.consumer.ConnectToNSQLookupds(nsqlookupdAddr)
			}
		} else {
			if len(nsqdAddr) == 1 {
				err = t.consumer.ConnectToNSQD(nsqdAddr[0])
			} else {
				err = t.consumer.ConnectToNSQDs(nsqdAddr)
			}
		}
		if err != nil {
			glog.Warningf("连接nsqlookupd(addr:%v)/nsqd(addr:%v)失败,err:%s", nsqlookupdAddr, nsqdAddr, err.Error())
			retryNum++
			sleepSeconds = 5
			if retryNum%6 == 0 {
				sleepSeconds = 30
			}
			time.Sleep(time.Duration(sleepSeconds) * time.Second)
			continue
		}

		glog.Infof("连接nsqlookupd(addr:%v)/nsqd(%v)成功", nsqlookupdAddr, nsqdAddr)
		break
	}

	err = nil
	return
}

// newConsumer 新建消费者
func newConsumer() consumer {
	return consumer{
		nsqdAddr:       make([]string, 0),
		nsqLookupdAddr: make([]string, 0),
		topics:         make(map[string]*topicInfo),
	}
}

// AddHandler 添加handler
func (c *consumer) AddHandler(topic string, handler nsq.HandlerFunc) {
	var (
		t  = &topicInfo{}
		ok bool
	)
	if t, ok = c.topics[topic]; !ok {
		t = &topicInfo{}
		t.concurrentNum = c.concurrent
		t.maxInFlight = c.maxInFlight
		t.config = nsq.NewConfig()
		t.config.MaxAttempts = c.maxAttempt
	}

	t.topic = topic
	// 自定义 handler
	t.handler = func(nm *nsq.Message) (err error) {
		err = handler(nm)
		if err != nil && Consumer.topics[topic].config.MaxAttempts > 0 && Consumer.topics[topic].config.MaxAttempts == nm.Attempts && Consumer.topics[topic].failHandler != nil {
			messageID := make([]byte, 0)
			for _, v := range nm.ID {
				messageID = append(messageID, v)
			}
			Consumer.topics[topic].failHandler(FailMessage{
				MessageID: string(messageID),
				Body:      nm.Body,
				Timestamp: nm.Timestamp,
				FailMsg:   err.Error(),
			})
			err = nil
		}
		return
	}
	c.topics[topic] = t
}

func (c *consumer) AddFailHandler(topic string, handler FailMessageFunc) {
	var (
		t  = &topicInfo{}
		ok bool
	)
	if t, ok = c.topics[topic]; !ok {
		t = &topicInfo{}
		t.concurrentNum = c.concurrent
		t.maxInFlight = c.maxInFlight
		t.config = nsq.NewConfig()
		t.config.MaxAttempts = c.maxAttempt
	}

	t.topic = topic
	t.failHandler = handler
	c.topics[topic] = t
}

// SetAddr 设置consumer地址
func (c *consumer) SetNsqlookupdAddr(node, addr string) {
	exist := false
	for _, v := range c.nsqLookupdAddr {
		if v == addr {
			exist = true
			break
		}
	}
	if !exist {
		c.nsqLookupdAddr = append(c.nsqLookupdAddr, addr)
	}
}

// SetMultiNsqlookupdAddr 设置多个consumer地址
func (c *consumer) SetMultiNsqlookupdAddr(node string, addrArr []string) {
	for _, v := range addrArr {
		c.SetNsqlookupdAddr(node, v)
	}
}

// SetNsqdAddr
func (c *consumer) SetNsqdAddr(node, addr string) {
	exist := false
	for _, v := range c.nsqdAddr {
		if v == addr {
			exist = true
			break
		}
	}
	if !exist {
		c.nsqdAddr = append(c.nsqdAddr, addr)
	}
}

// SetMultiNsqdAddr
func (c *consumer) SetMultiNsqdAddr(node string, addrArr []string) {
	for _, v := range addrArr {
		c.SetNsqdAddr(node, v)
	}
}

// StopAll 停止
func (c *consumer) Stop() {
	for k := range c.topics {
		c.topics[k].consumer.Stop()
	}
}

// Run 运行
func (c *consumer) Run() (err error) {
	if !c.isInit {
		err = errors.New("consumer not init")
		return
	}
	if len(c.nsqdAddr) == 0 && len(c.nsqLookupdAddr) == 0 {
		err = errors.New("nsqd addr or nsqlookupd address required")
		return
	}
	for _, topicInfo := range c.topics {
		topicInfo.config.MaxAttempts = c.maxAttempt
		topicInfo.config.MaxInFlight = c.maxInFlight
		go topicInfo.Connect(c.channelName, c.nsqdAddr, c.nsqLookupdAddr)
	}

	return
}

// Init 初始化
func (c *consumer) Initialize(config *ConsumerConfig) (err error) {
	//if nsqdAddr, ok := configSection["nsqd"]; ok {
	//	Consumer.nsqdAddr = strings.Split(nsqdAddr, ",")
	//}
	//if nsqlookupdArr, ok := configSection["nsqlookupd"]; ok {
	//	Consumer.nsqLookupdAddr = strings.Split(nsqlookupdArr, ",")
	//}
	//if maxInFlight, _ := strconv.Atoi(configSection["maxInFlight"]); maxInFlight > 0 {
	//	Consumer.maxInFlight = maxInFlight
	//}
	//if concurrent, _ := strconv.Atoi(configSection["concurrent"]); concurrent > 0 {
	//	Consumer.concurrent = concurrent
	//}
	//if channelName, _ := configSection["channelName"]; channelName != "" {
	//	Consumer.channelName = channelName
	//}
	//if Consumer.channelName == "" {
	//	err = errors.New("config channelName not found")
	//	return
	//}
	//if maxAttempt, _ := strconv.Atoi(configSection["max_attempt"]); maxAttempt > 0 {
	//	Consumer.maxAttempt = uint16(maxAttempt)
	//}
	//
	//if Consumer.maxInFlight < 1 {
	//	Consumer.maxInFlight = 1
	//}
	//if Consumer.concurrent < 1 {
	//	Consumer.concurrent = 1
	//}
	//
	//if Consumer.maxInFlight < Consumer.concurrent {
	//	err = errors.New("max_in_flight should exceed than concurrent")
	//	return
	//}

	Consumer.isInit = true
	return
}

var (
	Consumer = newConsumer()
)
*/
