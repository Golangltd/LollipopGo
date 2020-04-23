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

package message

type MessagesFilterType int

//const (
//	kMessagesFilterEmpty MessagesFilterType = 0
//	kMessagesFilterPhotos MessagesFilterType = 1
//	kMessagesFilterVideo MessagesFilterType = 2
//	kMessagesFilterPhotoVideo MessagesFilterType = 3
//	kMessagesFilterDocument MessagesFilterType = 4
//	kMessagesFilterUrl MessagesFilterType = 5
//	kMessagesFilterGif MessagesFilterType = 6
//	kMessagesFilterVoice MessagesFilterType = 7
//	kMessagesFilterMusic MessagesFilterType = 8
//	kMessagesFilterChatPhotos MessagesFilterType = 9
//	kMessagesFilterPhoneCalls MessagesFilterType = 10
//	kMessagesFilterRoundVoice MessagesFilterType = 11
//	kMessagesFilterRoundVideo MessagesFilterType = 12
//	kMessagesFilterMyMentions MessagesFilterType = 13
//	kMessagesFilterGeo MessagesFilterType = 14
//	kMessagesFilterContacts MessagesFilterType = 15
//)
//
//inputMessagesFilterEmpty#57e2f66c = MessagesFilter;
//inputMessagesFilterPhotos#9609a51c = MessagesFilter;
//inputMessagesFilterVideo#9fc00e65 = MessagesFilter;
//inputMessagesFilterPhotoVideo#56e9f0e4 = MessagesFilter;
//inputMessagesFilterDocument#9eddf188 = MessagesFilter;
//inputMessagesFilterUrl#7ef0dd87 = MessagesFilter;
//inputMessagesFilterGif#ffc86587 = MessagesFilter;
//inputMessagesFilterVoice#50f5c392 = MessagesFilter;
//inputMessagesFilterMusic#3751b49e = MessagesFilter;
//inputMessagesFilterChatPhotos#3a20ecb8 = MessagesFilter;
//inputMessagesFilterPhoneCalls#80c99768 flags:# missed:flags.0?true = MessagesFilter;
//inputMessagesFilterRoundVoice#7a7c17a4 = MessagesFilter;
//inputMessagesFilterRoundVideo#b549da53 = MessagesFilter;
//inputMessagesFilterMyMentions#c1f8e69a = MessagesFilter;
//inputMessagesFilterGeo#e7026d0d = MessagesFilter;
//inputMessagesFilterContacts#e062db83 = MessagesFilter;
//
//# filter识别
//## documentAttributeAudio即包含voice也包含music，通过voice进行识别，true为voice，否则为music
//
//````
//documentAttributeAudio#9852f9c6 flags:# voice:flags.10?true duration:int title:flags.0?string performer:flags.1?string waveform:flags.2?bytes = DocumentAttribute;
//````
//
//##
//inputMediaEmpty#9664f57f = InputMedia;
//inputMediaUploadedPhoto#1e287d04 flags:# file:InputFile stickers:flags.0?Vector<InputDocument> ttl_seconds:flags.1?int = InputMedia;
//inputMediaPhoto#b3ba0635 flags:# id:InputPhoto ttl_seconds:flags.0?int = InputMedia;
//inputMediaGeoPoint#f9c44144 geo_point:InputGeoPoint = InputMedia;
//inputMediaContact#f8ab7dfb phone_number:string first_name:string last_name:string vcard:string = InputMedia;
//inputMediaUploadedDocument#5b38c6c1 flags:# nosound_video:flags.3?true file:InputFile thumb:flags.2?InputFile mime_type:string attributes:Vector<DocumentAttribute> stickers:flags.0?Vector<InputDocument> ttl_seconds:flags.1?int = InputMedia;
//inputMediaDocument#23ab23d2 flags:# id:InputDocument ttl_seconds:flags.0?int = InputMedia;
//inputMediaVenue#c13d1c11 geo_point:InputGeoPoint title:string address:string provider:string venue_id:string venue_type:string = InputMedia;
//inputMediaGifExternal#4843b0fd url:string q:string = InputMedia;
//inputMediaPhotoExternal#e5bbfe1a flags:# url:string ttl_seconds:flags.0?int = InputMedia;
//inputMediaDocumentExternal#fb52dc99 flags:# url:string ttl_seconds:flags.0?int = InputMedia;
//inputMediaGame#d33f43f3 id:InputGame = InputMedia;
//inputMediaInvoice#f4e096c3 flags:# title:string description:string photo:flags.0?InputWebDocument invoice:Invoice payload:bytes provider:string provider_data:DataJSON start_param:string = InputMedia;
//inputMediaGeoLive#7b1a118f geo_point:InputGeoPoint period:int = InputMedia;
//
//
//
//messageEntityUnknown#bb92ba95 offset:int length:int = MessageEntity;
//messageEntityMention#fa04579d offset:int length:int = MessageEntity;
//messageEntityHashtag#6f635b0d offset:int length:int = MessageEntity;
//messageEntityBotCommand#6cef8ac7 offset:int length:int = MessageEntity;
//messageEntityUrl#6ed02538 offset:int length:int = MessageEntity;
//messageEntityEmail#64e475c2 offset:int length:int = MessageEntity;
//messageEntityBold#bd610bc9 offset:int length:int = MessageEntity;
//messageEntityItalic#826f8b60 offset:int length:int = MessageEntity;
//messageEntityCode#28a20571 offset:int length:int = MessageEntity;
//messageEntityPre#73924be0 offset:int length:int language:string = MessageEntity;
//messageEntityTextUrl#76a6d327 offset:int length:int url:string = MessageEntity;
//messageEntityMentionName#352dca58 offset:int length:int user_id:int = MessageEntity;
//inputMessageEntityMentionName#208e68c9 offset:int length:int user_id:InputUser = MessageEntity;
//messageEntityPhone#9b69e34b offset:int length:int = MessageEntity;
//messageEntityCashtag#4c4e743f offset:int length:int = MessageEntity;
//
//
//
//documentAttributeImageSize#6c37c15c w:int h:int = DocumentAttribute;
//documentAttributeAnimated#11b58939 = DocumentAttribute;
//documentAttributeSticker#6319d612 flags:# mask:flags.1?true alt:string stickerset:InputStickerSet mask_coords:flags.0?MaskCoords = DocumentAttribute;
//documentAttributeVideo#ef02ce6 flags:# round_message:flags.0?true supports_streaming:flags.1?true duration:int w:int h:int = DocumentAttribute;
//documentAttributeAudio#9852f9c6 flags:# voice:flags.10?true duration:int title:flags.0?string performer:flags.1?string waveform:flags.2?bytes = DocumentAttribute;
//documentAttributeFilename#15590068 file_name:string = DocumentAttribute;
//documentAttributeHasStickers#9801d2f7 = DocumentAttribute;
//