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

package rpc

import (
    "github.com/golang/glog"
    "github.com/nebulaim/telegramd/proto/mtproto"
    "golang.org/x/net/context"
    "github.com/nebulaim/telegramd/baselib/logger"
    "time"
)

// sync.getDifference flags:# auth_key_id:long user_id:int pts:int pts_total_limit:flags.0?int date:int qts:int = updates.Difference;
func (s *SyncServiceImpl) SyncGetDifference(ctx context.Context, request *mtproto.TLSyncGetDifference) (*mtproto.Updates_Difference, error) {
    glog.Infof("sync.getDifference - request: %s", logger.JsonDebugData(request))

    var (
        lastPts      = request.GetPts()
        otherUpdates []*mtproto.Update
        messages     []*mtproto.Message
        userList     []*mtproto.User
        chatList     []*mtproto.Chat
        difference   *mtproto.Updates_Difference
    )

    updateList := s.UpdateModel.GetUpdateListByGtPts(request.GetUserId(), request.GetPts())
    if len(updateList) == 0 {
        difference2 := &mtproto.TLUpdatesDifferenceEmpty{Data2: &mtproto.Updates_Difference_Data{
            Date: int32(time.Now().Unix()),
            Seq:  0,
        }}
        difference = difference2.To_Updates_Difference()
    } else {
        for _, update := range updateList {
            switch update.GetConstructor() {
            case mtproto.TLConstructor_CRC32_updateNewMessage:
                newMessage := update.To_UpdateNewMessage()
                messages = append(messages, newMessage.GetMessage())
                // otherUpdates = append(otherUpdates, update)

            case mtproto.TLConstructor_CRC32_updateReadHistoryOutbox:
                readHistoryOutbox := update.To_UpdateReadHistoryOutbox()
                readHistoryOutbox.SetPtsCount(0)
                otherUpdates = append(otherUpdates, readHistoryOutbox.To_Update())
            case mtproto.TLConstructor_CRC32_updateReadHistoryInbox:
                readHistoryInbox := update.To_UpdateReadHistoryInbox()
                readHistoryInbox.SetPtsCount(0)
                otherUpdates = append(otherUpdates, readHistoryInbox.To_Update())
            default:
                continue
            }
            if update.Data2.GetPts() > lastPts {
                lastPts = update.Data2.GetPts()
            }
        }
        if lastPts <= request.GetPts() {
            lastPts = 0
        }
        state := &mtproto.TLUpdatesState{Data2: &mtproto.Updates_State_Data{
            Pts:         lastPts,
            Date:        int32(time.Now().Unix()),
            UnreadCount: 0,
            // Seq:         int32(model.GetSequenceModel().CurrentSeqId(base2.Int32ToString(md.UserId))),
            Seq: 0,
        }}

        difference2 := &mtproto.TLUpdatesDifference{Data2: &mtproto.Updates_Difference_Data{
            NewMessages:  messages,
            OtherUpdates: otherUpdates,
            Users:        userList,
            Chats:        chatList,
            State:        state.To_Updates_State(),
        }}
        difference = difference2.To_Updates_Difference()
    }

    glog.Infof("sync.getDifference - reply: %s", logger.JsonDebugData(difference))
    return difference, nil
}
