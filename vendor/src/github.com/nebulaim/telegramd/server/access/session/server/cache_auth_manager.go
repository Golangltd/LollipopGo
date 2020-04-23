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

package server

import (
	"context"
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/base"
	"github.com/nebulaim/telegramd/baselib/cache"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/grpc_util/service_discovery"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

type CacheAuthInterface interface {
	GetAuthKey(int64) ([]byte, bool)
	GetUserID(int64) (int32, bool)
	GetLayer(int64) (int32, bool)
}

type cacheAuthValue struct {
	AuthKey []byte
	UserId  int32
	Layer   int32
}

// Impl cache.Value interface
func (cv *cacheAuthValue) Size() int {
	return 1
}

type cacheAuthManager struct {
	cache  *cache.LRUCache
	client mtproto.RPCSessionClient
}

var _cacheAuthManager *cacheAuthManager

func InitCacheAuthManager(cap int64, discovery *service_discovery.ServiceDiscoveryClientConfig) {
	conn, err := grpc_util.NewRPCClientByServiceDiscovery(discovery)
	if err != nil {
		glog.Error(err)
		panic(err)
	}

	_cacheAuthManager = &cacheAuthManager{
		cache:  cache.NewLRUCache(cap),
		client: mtproto.NewRPCSessionClient(conn),
	}
}

func (c *cacheAuthManager) GetAuthKey(authKeyId int64) ([]byte, bool) {
	var (
		cacheK = base.Int64ToString(authKeyId)
	)

	if v, ok := c.cache.Get(cacheK); !ok {
		keyInfo, err := c.client.SessionQueryAuthKey(context.Background(), &mtproto.TLSessionQueryAuthKey{AuthKeyId: authKeyId})
		if err != nil {
			glog.Error(err)
			return nil, false
		}
		//if r.Result != 0 {
		//	glog.Errorf("queryAuthKey err: {%v}", r)
		//	return nil, false
		//}
		c.cache.Set(cacheK, &cacheAuthValue{AuthKey: keyInfo.GetData2().GetAuthKey()})

		// TODO(@benqi): salt.
		return keyInfo.GetData2().GetAuthKey(), true
	} else {
		return v.(*cacheAuthValue).AuthKey, true
	}
}

func (c *cacheAuthManager) GetUserID(authKeyId int64) (int32, bool) {
	var (
		cacheK = base.Int64ToString(authKeyId)
	)

	if v, ok := c.cache.Peek(cacheK); !ok {
		glog.Error("not found authKeyId, bug???")
		return 0, false
	} else {
		cv, _ := v.(*cacheAuthValue)
		if cv.UserId == 0 {
			id, err := c.client.SessionGetUserId(context.Background(), &mtproto.TLSessionGetUserId{AuthKeyId: authKeyId})
			if err != nil {
				glog.Error(err)
				return 0, false
			}
			//if r.Result != 0 {
			//	glog.Errorf("queryAuthKey err: {%v}", r)
			//	return 0, false
			//}

			// update to cache
			cv.UserId = id.GetData2().GetV()
		}

		return cv.UserId, true
	}
}

func (c *cacheAuthManager) GetApiLayer(authKeyId int64) (int32, bool) {
	var (
		cacheK = base.Int64ToString(authKeyId)
	)

	if v, ok := c.cache.Peek(cacheK); !ok {
		glog.Error("not found authKeyId, bug???")
		return 0, false
	} else {
		cv, _ := v.(*cacheAuthValue)
		if cv.Layer == 0 {
			id, err := c.client.SessionGetLayer(context.Background(), &mtproto.TLSessionGetLayer{AuthKeyId: authKeyId})
			if err != nil {
				glog.Error(err)
				return 0, false
			}
			//if r.Result != 0 {
			//	glog.Errorf("queryAuthKey err: {%v}", r)
			//	return 0, false
			//}

			// update to cache
			cv.Layer = id.GetData2().GetV()
		}

		return cv.Layer, true
	}
}

func (c *cacheAuthManager) PutUserID(authKeyId int64, userId int32) {
	var (
		cacheK = base.Int64ToString(authKeyId)
	)

	if v, ok := c.cache.Peek(cacheK); ok {
		v.(*cacheAuthValue).UserId = userId
	} else {
		glog.Error("not found authKeyId, bug???")
	}
}

func getCacheUserID(authKeyId int64) int32 {
	if _cacheAuthManager == nil {
		panic("not init cacheAuthManager.")
	}

	userId, _ := _cacheAuthManager.GetUserID(authKeyId)
	return userId
}

func putCacheUserId(authKeyId int64, useId int32) {
	if _cacheAuthManager == nil {
		panic("not init cacheAuthManager.")
	}

	_cacheAuthManager.PutUserID(authKeyId, useId)
}

func getCacheAuthKey(authKeyId int64) []byte {
	if _cacheAuthManager == nil {
		panic("not init cacheAuthManager.")
	}

	key, _ := _cacheAuthManager.GetAuthKey(authKeyId)
	return key
}

func getCacheApiLayer(authKeyId int64) int32 {
	if _cacheAuthManager == nil {
		panic("not init cacheAuthManager.")
	}

	layer, _ := _cacheAuthManager.GetApiLayer(authKeyId)
	return layer
}

func uploadInitConnection(authKeyId int64, layer int32, ip string, initConnection *TLInitConnectionExt) error {
	session := &mtproto.TLClientSessionInfo { Data2: &mtproto.ClientSession_Data{
		AuthKeyId:      authKeyId,
		Ip:             ip,
		Layer:          layer,
		ApiId:          initConnection.ApiId,
		DeviceModel:    initConnection.DeviceMode,
		SystemVersion:  initConnection.SystemVersion,
		AppVersion:     initConnection.AppVersion,
		SystemLangCode: initConnection.SystemLangCode,
		LangPack:       initConnection.LangPack,
		LangCode:       initConnection.LangCode,
	}}

	request := &mtproto.TLSessionSetClientSessionInfo{
		Session: session.To_ClientSession(),
	}

	_, err := _cacheAuthManager.client.SessionSetClientSessionInfo(context.Background(), request)

	if err != nil {
		glog.Error(err)
	}

	return err
}
