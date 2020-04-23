package dialog

import (
	"github.com/nebulaim/telegramd/biz/dal/dataobject"
	"github.com/nebulaim/telegramd/proto/mtproto"
	"github.com/nebulaim/telegramd/biz/base"
	base2 "github.com/nebulaim/telegramd/baselib/base"
	"encoding/json"
	"github.com/nebulaim/telegramd/biz/core"
)

func dialogDOToDialog(dialogDO *dataobject.UserDialogsDO) *mtproto.TLDialog {
	dialog := mtproto.NewTLDialog()

	dialog.SetPinned(dialogDO.IsPinned == 1)
	dialog.SetPeer(base.ToPeerByTypeAndID(dialogDO.PeerType, dialogDO.PeerId))
	//if dialogDO.PeerType == base.PEER_CHANNEL {
	//	// TODO(@benqi): only channel has pts
	//	// dialog.SetPts(messageBoxsDO.Pts)
	//	// peerChannlIdList = append(peerChannlIdList, dialogDO.PeerId)
	//	dialog.SetPts(dialogDO.Pts)
	//}

	dialog.SetTopMessage(dialogDO.TopMessage)
	dialog.SetReadInboxMaxId(dialogDO.ReadInboxMaxId)
	dialog.SetReadOutboxMaxId(dialogDO.ReadOutboxMaxId)
	dialog.SetUnreadCount(dialogDO.UnreadCount)
	dialog.SetUnreadMentionsCount(dialogDO.UnreadMentionsCount)

	// TODO(@benqi): draft message.
	if dialogDO.DraftType == 2 {
		draft := &mtproto.DraftMessage{}
		err := json.Unmarshal([]byte(dialogDO.DraftMessageData), &draft)
		if err == nil {
			dialog.SetDraft(draft)
		}
	} else {
		dialog.SetDraft(mtproto.NewTLDraftMessageEmpty().To_DraftMessage())
	}

	// NotifySettings
	peerNotifySettings := mtproto.NewTLPeerNotifySettings()
	if dialogDO.ShowPreviews == 1 {
		peerNotifySettings.SetShowPreviews(mtproto.ToBool(true))
	}
	if dialogDO.Silent == 1 {
		peerNotifySettings.SetSilent(mtproto.ToBool(true))
	}
	if dialogDO.MuteUntil == 0 {
		peerNotifySettings.SetMuteUntil(1)
	} else {
		peerNotifySettings.SetMuteUntil(dialogDO.MuteUntil)
	}
	if dialogDO.Sound == "" {
		peerNotifySettings.SetSound("default")
	} else {
		peerNotifySettings.SetSound(dialogDO.Sound)
	}
	dialog.SetNotifySettings(peerNotifySettings.To_PeerNotifySettings())
	return dialog
}

func (m *DialogModel) dialogDOListToDialogList(dialogDOList []dataobject.UserDialogsDO) (dialogs []*mtproto.Dialog) {
	// draftIdList := make([]int32, 0)
	channelIdList := make([]int32, 0, len(dialogDOList))
	for i := 0; i < len(dialogDOList); i++ {
		//if dialogDO.DraftId > 0 {
		//	draftIdList = append(draftIdList, dialogDO.DraftId)
		//}
		dialogDO := &dialogDOList[i]
		dialog := dialogDOToDialog(dialogDO)
		if dialogDO.PeerType == base.PEER_CHANNEL {
			channelIdList = append(channelIdList, dialogDO.PeerId)
		}
		dialogs = append(dialogs, dialog.To_Dialog())
	}

	topMessageMap := m.channelCallback.GetTopMessageListByIdList(channelIdList)
	for _, dialog := range dialogs {
		if dialog.Data2.Peer.GetConstructor() == mtproto.TLConstructor_CRC32_peerChannel {
			dialog.Data2.TopMessage = topMessageMap[dialog.Data2.Peer.Data2.ChannelId]
			dialog.Data2.Pts = int32(core.NextChannelPtsId(dialog.Data2.Peer.Data2.ChannelId))
		}
	}

	// TODO(@benqi): fetch draft message list
	return
}

func (m *DialogModel) GetDialogsByOffsetId(userId int32, isPinned bool, offsetId int32, limit int32) (dialogs []*mtproto.Dialog) {
	dialogDOList := m.dao.UserDialogsDAO.SelectByPinnedAndOffset(
		userId, base2.BoolToInt8(isPinned), offsetId, limit)
	dialogs = m.dialogDOListToDialogList(dialogDOList)
	return
}

func (m *DialogModel) GetPeersDialogs(selfId int32, peers []*base.PeerUtil) (dialogs []*mtproto.Dialog) {
	channelIdList := make([]int32, 0, len(peers))

	for _, peer := range peers {
		// peerUtil := base.FromInputPeer2(selfId, peer)
		dialogDO := m.dao.UserDialogsDAO.SelectByPeer(selfId, int8(peer.PeerType), peer.PeerId)
		if dialogDO != nil {
			dialogs = append(dialogs, dialogDOToDialog(dialogDO).To_Dialog())
			if dialogDO.PeerType == base.PEER_CHANNEL {
				channelIdList = append(channelIdList, dialogDO.PeerId)
			}
		}
	}

	topMessageMap := m.channelCallback.GetTopMessageListByIdList(channelIdList)
	for _, dialog := range dialogs {
		if dialog.Data2.Peer.GetConstructor() == mtproto.TLConstructor_CRC32_peerChannel {
			dialog.Data2.TopMessage = topMessageMap[dialog.Data2.Peer.Data2.ChannelId]
			dialog.Data2.Pts = int32(core.NextChannelPtsId(dialog.Data2.Peer.Data2.ChannelId))
		}
	}
	return
}

func (m *DialogModel) GetPinnedDialogs(userId int32) (dialogs []*mtproto.Dialog) {
	dialogDOList := m.dao.UserDialogsDAO.SelectPinnedDialogs(userId)
	dialogs = m.dialogDOListToDialogList(dialogDOList)
	return
}
