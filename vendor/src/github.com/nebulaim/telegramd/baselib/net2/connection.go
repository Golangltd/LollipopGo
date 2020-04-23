/*
 *  Copyright (c) 2017, https://github.com/nebulaim
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

package net2

import (
//"net"
//"sync/atomic"
//"fmt"
//"sync"
)

type Connection interface {
	GetConnID() uint64
	IsClosed() bool
	Close() error
	Codec() Codec
	Receive() (interface{}, error)
	Send(msg interface{}) error
}

type closeCallback interface {
	// func(Connection)
	OnConnectionClosed(Connection)
}

//type Connection2 struct {
//	connID        uint64
//	name          string			// name
//	conn          net.Conn
//	codec         Codec
//	sendChan      chan interface{}
//	recvMutex     sync.Mutex
//	sendMutex     sync.RWMutex
//	closeFlag     int32
//	closeChan     chan int
//	closeMutex    sync.Mutex
//	closeCallback closeCallback
//	Context       interface{}
//}
//
//func NewConnection2(name string, conn net.Conn, sendChanSize int, codec Codec, cb closeCallback) *Connection2 {
//	//if globalConnectionId >= 0xfffffffffffffff {
//	//	atomic.StoreUint64(&globalConnectionId, 0)
//	//}
//	conn2 := &Connection2{
//		name:          name,
//		conn:          conn,
//		codec:         codec,
//		closeChan:     make(chan int),
//		connID:        atomic.AddUint64(&globalConnectionId, 1),
//		closeCallback: cb,
//	}
//
//	if sendChanSize > 0 {
//		conn2.sendChan = make(chan interface{}, sendChanSize)
//		go conn2.sendLoop()
//	}
//	return conn2
//}
//
//func (c *Connection2) String() string {
//	return fmt.Sprintf("{connID: %d, name: %s, lAddr: %s, rAddr: %s}", c.connID, c.name, c.conn.LocalAddr(), c.conn.RemoteAddr())
//}
//
//func (c *Connection2) LoadAddr() net.Addr {
//	return c.conn.LocalAddr()
//}
//
//func (c *Connection2) RemoteAddr() net.Addr {
//	return c.conn.RemoteAddr()
//}
//
//func (c *Connection2) Name() string {
//	return c.name
//}
//
//func (c *Connection2) GetConnID() uint64 {
//	return c.connID
//}
//
//func (c *Connection2) GetNetConn() net.Conn {
//	return c.conn
//}
//
//func (c *Connection2) IsClosed() bool {
//	return atomic.LoadInt32(&c.closeFlag) == 1
//}
//
//func (c *Connection2) Close() error {
//	if atomic.CompareAndSwapInt32(&c.closeFlag, 0, 1) {
//		if c.closeCallback != nil {
//			c.closeCallback.OnConnectionClosed(c)
//		}
//
//		close(c.closeChan)
//
//		if c.sendChan != nil {
//			c.sendMutex.Lock()
//			close(c.sendChan)
//			if clear, ok := c.codec.(ClearSendChan); ok {
//				clear.ClearSendChan(c.sendChan)
//			}
//			c.sendMutex.Unlock()
//		}
//
//		err := c.codec.Close()
//		return err
//	}
//	return ConnectionClosedError
//}
//
//func (c *Connection2) Codec() Codec {
//	return c.codec
//}
//
//func (c *Connection2) Receive() (interface{}, error) {
//	c.recvMutex.Lock()
//	defer c.recvMutex.Unlock()
//
//	msg, err := c.codec.Receive()
//	if err != nil {
//		c.Close()
//	}
//	return msg, err
//}
//
//func (c *Connection2) sendLoop() {
//	defer c.Close()
//	for {
//		select {
//		case msg, ok := <-c.sendChan:
//			if !ok || c.codec.Send(msg) != nil {
//				return
//			}
//		case <-c.closeChan:
//			return
//		}
//	}
//}
//
//func (c *Connection2) Send(msg interface{}) error {
//	if c.sendChan == nil {
//		if c.IsClosed() {
//			return ConnectionClosedError
//		}
//
//		c.sendMutex.Lock()
//		defer c.sendMutex.Unlock()
//
//		err := c.codec.Send(msg)
//		if err != nil {
//			c.Close()
//		}
//		return err
//	}
//
//	c.sendMutex.RLock()
//	if c.IsClosed() {
//		c.sendMutex.RUnlock()
//		return ConnectionClosedError
//	}
//
//	select {
//	case c.sendChan <- msg:
//		c.sendMutex.RUnlock()
//		return nil
//	default:
//		c.sendMutex.RUnlock()
//		c.Close()
//		return ConnectionBlockedError
//	}
//}
