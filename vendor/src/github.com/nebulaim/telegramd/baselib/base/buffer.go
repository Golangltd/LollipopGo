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

package base

import (
	"encoding/binary"
	"fmt"
	"io"
	"unicode/utf8"
)

type Buffer struct {
	buf     []byte
	readPos int
}

func NewBuffer(data []byte) *Buffer {
	return &Buffer{buf: data}
}

func MakeBuffer(size, capacity int) *Buffer {
	return NewBuffer(make([]byte, size, capacity))
}

func (b *Buffer) Grows(n int) (i int) {
	i = len(b.buf)

	newLen := len(b.buf) + n
	if cap(b.buf) >= newLen {
		b.buf = b.buf[:newLen]
		return
	}

	data := make([]byte, newLen, cap(b.buf)/4+newLen)
	copy(data, b.buf)
	b.buf = data
	return
}

func (b *Buffer) Truncate(size int) {
	if len(b.buf) > size {
		b.buf = b.buf[:size]
		if b.readPos > size {
			b.readPos = size
		}
	}
}

func (b *Buffer) ResetUndelay(data []byte) {
	b.buf = data
	b.readPos = 0
}

func (b *Buffer) Reset() {
	b.Truncate(0)
}

func (b *Buffer) Bytes() []byte {
	return b.buf[b.readPos:]
}

func (b *Buffer) Len() int {
	return len(b.buf) - b.readPos
}

func (b *Buffer) Cap() int {
	return cap(b.buf)
}

func (b *Buffer) Read(p []byte) (int, error) {
	if len(p) == 0 {
		return 0, nil
	}
	if b.readPos >= len(b.buf) {
		return 0, io.EOF
	}
	n := copy(p, b.buf[b.readPos:])
	b.readPos += n
	return n, nil
}

func (b *Buffer) ReadByte() (byte, error) {
	if b.readPos >= len(b.buf) {
		return 0, io.EOF
	}
	r := b.buf[b.readPos]
	b.readPos += 1
	return r, nil
}

func (b *Buffer) ReadUint16(order binary.ByteOrder) (uint16, error) {
	if b.readPos >= len(b.buf)-2 {
		return 0, io.EOF
	}
	u := order.Uint16(b.buf[b.readPos:])
	b.readPos += 2
	return u, nil
}

func (b *Buffer) ReadUint32(order binary.ByteOrder) (uint32, error) {
	if b.readPos >= len(b.buf)-4 {
		return 0, io.EOF
	}
	u := order.Uint32(b.buf[b.readPos:])
	b.readPos += 4
	return u, nil
}

func (b *Buffer) ReadUint64(order binary.ByteOrder) (uint64, error) {
	if b.readPos >= len(b.buf)-8 {
		return 0, io.EOF
	}
	u := order.Uint64(b.buf[b.readPos:])
	b.readPos += 8
	return u, nil
}

func (b *Buffer) Skip(i int) int {
	pos := b.readPos + i
	if pos >= 0 && pos < len(b.buf) {
		b.readPos = pos
		return pos
	}

	return -1
}

func (b *Buffer) ReadAt(p []byte, off int64) (int, error) {
	if off < 0 {
		return 0, fmt.Errorf("binary.Buffer.ReadAt: negative offset")
	}

	if int(off) >= len(b.buf) {
		return 0, io.EOF
	}
	n := len(p)
	if n+int(off) > len(b.buf) {
		n = len(b.buf) - int(off)
	}
	copy(p, b.buf[off:])
	return n, nil
}

func (b *Buffer) ReadRune() (rune, int, error) {
	if b.readPos == len(b.buf) {
		return 0, 0, io.EOF
	}
	if c := b.buf[b.readPos]; c < utf8.RuneSelf {
		b.readPos += 1
		return rune(c), 1, nil
	}
	r, n := utf8.DecodeRune(b.buf[b.readPos:])
	b.readPos += n
	return r, n, nil
}

func (b *Buffer) ReadBytes(delim byte) ([]byte, error) {
	if b.readPos >= len(b.buf) {
		return nil, io.EOF
	}
	s := b.readPos
	for i := b.readPos; i < len(b.buf); i++ {
		if b.buf[i] == delim {
			b.readPos = i + 1
			return b.buf[s:b.readPos], nil
		}
	}
	return nil, io.EOF
}

func (b *Buffer) Write(p []byte) (int, error) {
	i := b.Grows(len(p))
	copy(b.buf[i:], p)
	return len(p), nil
}

//func (b *Buffer) WriteString(s string) (int, error) {
//	return b.Write(unsafe2.Bytes(s))
//}

func (b *Buffer) WriteByte(c byte) error {
	i := b.Grows(1)
	b.buf[i] = c
	return nil
}

func (b *Buffer) WriteUint16(u uint16, order binary.ByteOrder) error {
	i := b.Grows(2)
	order.PutUint16(b.buf[i:], u)
	return nil
}

func (b *Buffer) WriteUint32(u uint32, order binary.ByteOrder) error {
	i := b.Grows(4)
	order.PutUint32(b.buf[i:], u)
	return nil
}

func (b *Buffer) WriteUint64(u uint64, order binary.ByteOrder) error {
	i := b.Grows(8)
	order.PutUint64(b.buf[i:], u)
	return nil
}

func (b *Buffer) WriteRune(r rune) (int, error) {
	i := b.Grows(utf8.UTFMax)
	s := utf8.EncodeRune(b.buf[i:], r)
	n := utf8.UTFMax - s
	b.buf = b.buf[:len(b.buf)-n]
	return s, nil
}

func (b *Buffer) String() string {
	return string(b.Bytes())
}
