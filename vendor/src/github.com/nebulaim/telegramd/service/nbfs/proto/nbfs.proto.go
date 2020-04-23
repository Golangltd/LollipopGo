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

package nbfs

type DocumentFileMetadata struct {
	FileId           int64
	DocumentId       int64
	AccessHash       int64
	DcId             int32
	FileSize         int32
	FilePath         string
	UploadedFileName string
	Ext              string
	Md5Hash          string
	MimeType         string
}

type PhotoFileMetadata struct {
	FileId    int64
	PhotoId   int64
	PhotoType int8
	SizeType  string
	DcId      int32
	VolumeId  int64
	LocalId   int32
	SecretId  int64
	Width     int32
	Height    int32
	FileSize  int32
	FilePath  string
	Ext       string
}
