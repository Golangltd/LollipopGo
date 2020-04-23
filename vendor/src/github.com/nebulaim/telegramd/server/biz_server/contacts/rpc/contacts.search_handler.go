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

package rpc

import (
	"github.com/golang/glog"
	"github.com/nebulaim/telegramd/baselib/grpc_util"
	"github.com/nebulaim/telegramd/baselib/logger"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"golang.org/x/net/context"
)

// contacts.search#11f812d8 q:string limit:int = contacts.Found;
func (s *ContactsServiceImpl) ContactsSearch(ctx context.Context, request *mtproto.TLContactsSearch) (*mtproto.Contacts_Found, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.search#11f812d8 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	// Check query string and limit
	if len(request.Q) < 5 || request.Limit < 1 {
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err, ": query or limit invalid")
		return nil, err
	}

	contactLogic := s.ContactModel.MakeContactLogic(md.UserId)
	idList := contactLogic.SearchContacts(request.Q, request.Limit)

	// results
	results := make([]*mtproto.Peer, 0, len(idList))
	for _, id := range idList {
		peer := &mtproto.TLPeerUser{Data2: &mtproto.Peer_Data{
			UserId: id,
		}}
		results = append(results, peer.To_Peer())
	}

	// users
	users := s.UserModel.GetUsersBySelfAndIDList(md.UserId, idList)

	found := &mtproto.TLContactsFound{Data2: &mtproto.Contacts_Found_Data{
		Results: results,
		Users:   users,
	}}

	glog.Infof("contacts.search#11f812d8 - reply: ", logger.JsonDebugData(found))
	return found.To_Contacts_Found(), nil
}
