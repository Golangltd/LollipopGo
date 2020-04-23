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
	"github.com/nebulaim/telegramd/service/document/client"
	"golang.org/x/net/context"
	"time"
	"github.com/nebulaim/telegramd/biz/base"
)

// users.getFullUser#ca30a5b1 id:InputUser = UserFull;
func (s *UsersServiceImpl) UsersGetFullUser(ctx context.Context, request *mtproto.TLUsersGetFullUser) (*mtproto.UserFull, error) {
	md := grpc_util.RpcMetadataFromIncoming(ctx)
	glog.Infof("users.getFullUser#ca30a5b1 - metadata: %s, request: %s", logger.JsonDebugData(md), logger.JsonDebugData(request))

	var (
		user *mtproto.User
		peer *base.PeerUtil
	)

	id := request.GetId()
	switch id.GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUserSelf:
		peer = &base.PeerUtil{
			PeerType: base.PEER_USER,
			PeerId:   md.UserId,
		}
	case mtproto.TLConstructor_CRC32_inputUser:
		peer = &base.PeerUtil{
			PeerType:   base.PEER_USER,
			PeerId:     id.GetData2().GetUserId(),
			AccessHash: id.GetData2().GetAccessHash(),
		}
	default:

	}

	fullUser := mtproto.NewTLUserFull()
	fullUser.SetPhoneCallsAvailable(true)
	fullUser.SetPhoneCallsPrivate(false)
	fullUser.SetAbout("@Benqi")
	fullUser.SetCommonChatsCount(0)

	switch request.GetId().GetConstructor() {
	case mtproto.TLConstructor_CRC32_inputUserSelf:
		user = s.UserModel.GetUserById(md.UserId, md.UserId).To_User()
		fullUser.SetUser(user)
		// Link
		link := &mtproto.TLContactsLink{Data2: &mtproto.Contacts_Link_Data{
			MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
			ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
			User:        user,
		}}
		fullUser.SetLink(link.To_Contacts_Link())
	case mtproto.TLConstructor_CRC32_inputUser:
		inputUser := request.GetId().To_InputUser()
		user = s.UserModel.GetUserById(md.UserId, inputUser.GetUserId()).To_User()
		fullUser.SetUser(user)

		// Link
		link := &mtproto.TLContactsLink{Data2: &mtproto.Contacts_Link_Data{
			MyLink:      mtproto.NewTLContactLinkContact().To_ContactLink(),
			ForeignLink: mtproto.NewTLContactLinkContact().To_ContactLink(),
			User:        user,
		}}
		fullUser.SetLink(link.To_Contacts_Link())
	case mtproto.TLConstructor_CRC32_inputUserEmpty:
		// TODO(@benqi): BAD_REQUEST: 400
		err := mtproto.NewRpcError2(mtproto.TLRpcErrorCodes_BAD_REQUEST)
		glog.Error(err)
		return nil, err
	}

	// NotifySettings
	peerNotifySettings := s.AccountModel.GetNotifySettings(md.UserId, peer)
	fullUser.SetNotifySettings(peerNotifySettings)

	photoId := user.GetData2().GetPhoto().GetData2().GetPhotoId()
	sizes, _ := document_client.GetPhotoSizeList(photoId)
	photo := &mtproto.TLPhoto{Data2: &mtproto.Photo_Data{
		Id:          photoId,
		HasStickers: false,
		AccessHash:  photoId, // photo2.GetFileAccessHash(file.GetData2().GetId(), file.GetData2().GetParts()),
		Date:        int32(time.Now().Unix()),
		Sizes:       sizes,
	}}
	fullUser.SetProfilePhoto(photo.To_Photo())

	glog.Infof("users.getFullUser#ca30a5b1 - reply: %s", logger.JsonDebugData(fullUser))
	return fullUser.To_UserFull(), nil
}
