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

package snowflake

import (
	"fmt"
	"log"
	"testing"
)

func TestID(t *testing.T) {
	id, err := NewIdWorker(0, 0, twepoch)
	if err != nil {
		fmt.Printf("NewIdWorker(0, 0) error(%v)\n", err)
		t.FailNow()
	}
	sid, err := id.NextId()
	if err != nil {
		fmt.Printf("id.NextId() error(%v)\n", err)
		t.FailNow()
	}
	log.Printf("snowflake id: %d\n", sid)
	sids, err := id.NextIds(10)
	if err != nil {
		fmt.Printf("id.NextId() error(%v)\n", err)
		t.FailNow()
	}
	fmt.Printf("snowflake ids: %v\n", sids)
}

func BenchmarkID(b *testing.B) {
	id, err := NewIdWorker(0, 0, twepoch)
	if err != nil {
		fmt.Printf("NewIdWorker(0, 0) error(%v)\n", err)
		b.FailNow()
	}
	for i := 0; i < b.N; i++ {
		if _, err := id.NextId(); err != nil {
			b.FailNow()
		}
	}
}
