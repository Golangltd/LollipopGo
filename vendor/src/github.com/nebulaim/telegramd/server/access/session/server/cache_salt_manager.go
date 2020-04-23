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

/**
 ## Android client source code, gen salts:
 ### handshake:
   - client received TL_server_DH_params_ok:
	```
	handshakeServerSalt = new TL_future_salt();
	handshakeServerSalt->valid_since = currentTime + timeDifference - 5;
	handshakeServerSalt->valid_until = handshakeServerSalt->valid_since + 30 * 60;
	for (int32_t a = 7; a >= 0; a--) {
		handshakeServerSalt->salt <<= 8;
		handshakeServerSalt->salt |= (authNewNonce->bytes[a] ^ authServerNonce->bytes[a]);
	}
	```

   - client received TL_dh_gen_ok:
	```
	std::unique_ptr<TL_future_salt> salt = std::unique_ptr<TL_future_salt>(handshakeServerSalt);
	addServerSalt(salt);
	handshakeServerSalt = nullptr;
	```

 ### received TL_new_session_created:
	```
	std::unique_ptr<TL_future_salt> salt = std::unique_ptr<TL_future_salt>(new TL_future_salt());
	salt->valid_until = salt->valid_since = getCurrentTime();
	salt->valid_until += 30 * 60;
	salt->salt = response->server_salt;
	datacenter->addServerSalt(salt);

	```

 ### send TL_get_future_salts request:
   - rpc request:
	```
    requestingSaltsForDc.push_back(datacenter->getDatacenterId());
    TL_get_future_salts *request = new TL_get_future_salts();
    request->num = 32;
    sendRequest(request, [&, datacenter](TLObject *response, TL_error *error, int32_t networkType) {
        std::vector<uint32_t>::iterator iter = std::find(requestingSaltsForDc.begin(), requestingSaltsForDc.end(), datacenter->getDatacenterId());
        if (iter != requestingSaltsForDc.end()) {
            requestingSaltsForDc.erase(iter);
        }
        if (error == nullptr) {
            TL_future_salts *res = (TL_future_salts *) response;
            datacenter->mergeServerSalts(res->salts);
            saveConfig();
        }
    }, nullptr, RequestFlagWithoutLogin | RequestFlagEnableUnauthorized, datacenter->getDatacenterId(), ConnectionTypeGeneric, true);

	```

  - rpc response:
	```
	TL_future_salts *response = (TL_future_salts *) message;
	int64_t requestMid = response->req_msg_id;
	for (requestsIter iter = runningRequests.begin(); iter != runningRequests.end(); iter++) {
		Request *request = iter->get();
		if (request->respondsToMessageId(requestMid)) {
			request->onComplete(response, nullptr, connection->currentNetworkType);
			request->completed = true;
			runningRequests.erase(iter);
			break;
		}
	}
	```

  - received TL_bad_server_salt:
	```
	datacenter->clearServerSalts();

	std::unique_ptr<TL_future_salt> salt = std::unique_ptr<TL_future_salt>(new TL_future_salt());
	salt->valid_until = salt->valid_since = getCurrentTime();
	salt->valid_until += 30 * 60;
	salt->salt = messageSalt;
	datacenter->addServerSalt(salt);
	```
*/

// salt cache
package server

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/cache"
	"github.com/nebulaim/telegramd/proto/mtproto"
)

const (
	kSaltTimeout = 30 * 60 // salt timeout
	// kCacheConfig = `{"conn":":"127.0.0.1:6039"}`
	kAdapterName     = "memory"
	kCacheConfig     = `{"interval":60}`
	kCacheSaltPrefix = "salts"
)

var cacheSalts *cacheSaltManager = nil

func init() {
	rand.Seed(time.Now().UnixNano())
	initCacheSaltsManager(kAdapterName, kCacheConfig)
}

type cacheSaltManager struct {
	cache   cache.Cache
	timeout time.Duration // salt timeout
}

func initCacheSaltsManager(name, config string) error {
	if config == "" {
		config = kCacheConfig
	}

	c, err := cache.NewCache(name, kCacheConfig)
	if err != nil {
		glog.Error(err)
		return err
	}

	cacheSalts = &cacheSaltManager{cache: c, timeout: kSaltTimeout}
	return nil
}

func genCacheSaltKey(id int64) string {
	return fmt.Sprintf("%s_%d", kCacheSaltPrefix, id)
}

func GetOrInsertSaltList(keyId int64, size int) ([]*mtproto.TLFutureSalt, error) {
	var (
		salts = make([]*mtproto.TLFutureSalt, size)

		date           = int32(time.Now().Unix())
		lastValidUntil = date
		// ok = false
		saltsData []*mtproto.FutureSalt_Data
	)

	v := cacheSalts.cache.Get(genCacheSaltKey(keyId))
	if v != nil {
		if saltList, ok := v.([]*mtproto.FutureSalt_Data); ok {
			for _, salt := range saltList {
				if salt.ValidUntil >= date {
					saltsData = append(saltsData, salt)
					if lastValidUntil < salt.ValidUntil {
						lastValidUntil = salt.ValidUntil
					}
				}
			}
		}
	}

	left := size - len(saltsData)
	if left > 0 {
		for i := 0; i < size; i++ {
			salt := &mtproto.FutureSalt_Data{
				ValidSince: lastValidUntil,
				ValidUntil: lastValidUntil + kSaltTimeout,
				Salt:       rand.Int63(),
			}
			saltsData = append(saltsData, salt)
			lastValidUntil += kSaltTimeout
		}
	}

	for i := 0; i < size; i++ {
		salt := &mtproto.TLFutureSalt{
			Data2: saltsData[i],
		}
		salts[i] = salt
	}

	if left > 0 {
		err := cacheSalts.cache.Put(genCacheSaltKey(keyId), saltsData, time.Duration(len(saltsData))*kSaltTimeout*time.Second)
		if err != nil {
			glog.Error(err)
			return nil, err
		}
	}
	return salts, nil
}

func GetOrInsertSalt(keyId int64) (int64, error) {
	var date = int32(time.Now().Unix())

	var salt int64 = 0
	v := cacheSalts.cache.Get(genCacheSaltKey(keyId))
	if v != nil {
		if saltList, ok := v.([]*mtproto.FutureSalt_Data); ok {
			for _, v := range saltList {
				if v.ValidSince <= date && v.ValidUntil > date {
					salt = v.Salt
					break
				}
			}
		}
	}

	if salt == 0 {
		salt = rand.Int63()
		saltData := &mtproto.FutureSalt_Data{
			ValidSince: date,
			ValidUntil: date + kSaltTimeout + kSaltTimeout,
			Salt:       salt,
		}

		// Put cache.
		err := cacheSalts.cache.Put(genCacheSaltKey(keyId), []*mtproto.FutureSalt_Data{saltData}, kSaltTimeout*time.Second)
		if err != nil {
			glog.Error(err)
			return 0, err
		}
	}

	return salt, nil
}

// https://core.telegram.org/mtproto/description#server-salt

// Server Salt
//
// A (random) 64-bit number periodically (say, every 24 hours) changed
// (separately for each session) at the request of the server.
// All subsequent messages must contain the new salt
// (although, messages with the old salt are still accepted for a further 300 seconds).
// Required to protect against replay attacks and certain tricks
// associated with adjusting the client clock to a moment in the distant future.
//
func CheckBySalt(keyId, salt int64) bool {
	var date = int32(time.Now().Unix())

	// var salt int64 = 0
	v := cacheSalts.cache.Get(genCacheSaltKey(keyId))
	if v == nil {
		return false
	}

	if saltList, ok := v.([]*mtproto.FutureSalt_Data); ok {
		for _, v := range saltList {
			// old salt are still accepted for a further 300 seconds
			if v.ValidSince <= date && v.ValidUntil+kSaltTimeout > date && salt == v.Salt {
				return true
			}
		}
	}

	return false
}
