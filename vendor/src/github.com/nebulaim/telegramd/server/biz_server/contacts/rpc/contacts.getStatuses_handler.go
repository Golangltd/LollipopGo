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

// contacts.getStatuses#c4a353ee = Vector<ContactStatus>;
func (s *ContactsServiceImpl) ContactsGetStatuses(ctx context.Context, request *mtproto.TLContactsGetStatuses) (*mtproto.Vector_ContactStatus, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("contacts.getStatuses#c4a353ee - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	contactLogic := s.ContactModel.MakeContactLogic(md.UserId)
	cList := contactLogic.GetContactList()

	statusList := &mtproto.Vector_ContactStatus{
		Datas: make([]*mtproto.ContactStatus, 0, len(cList)),
	}

	for _, c := range cList {
		contactStatus := &mtproto.TLContactStatus{Data2: &mtproto.ContactStatus_Data{
			UserId: c.ContactUserId,
			Status: s.UserModel.GetUserStatus(c.ContactUserId),
		}}
		statusList.Datas = append(statusList.Datas, contactStatus.To_ContactStatus())
	}

	glog.Infof("contacts.getStatuses#c4a353ee - reply: ", logger.JsonDebugData(statusList))
	return statusList, nil
}
