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

package main

import (
	"flag"
	"github.com/nebulaim/telegramd/baselib/app"
	"github.com/nebulaim/telegramd/server/access/session/server"
)

/*
  // Subscriber-related

  def subscriber(consumer: ActorRef): Receive = {
    case OnNext(cmd: SubscribeCommand) ⇒
      cmd match {
        case SubscribeToOnline(userIds) ⇒
          consumer ! UpdatesConsumerMessage.SubscribeToUserPresences(userIds.toSet)
        case SubscribeFromOnline(userIds) ⇒
          consumer ! UpdatesConsumerMessage.UnsubscribeFromUserPresences(userIds.toSet)
        case SubscribeToGroupOnline(groupIds) ⇒
          consumer ! UpdatesConsumerMessage.SubscribeToGroupPresences(groupIds.toSet)
        case SubscribeFromGroupOnline(groupIds) ⇒
          consumer ! UpdatesConsumerMessage.UnsubscribeFromGroupPresences(groupIds.toSet)
        case SubscribeToSeq(_) ⇒
          consumer ! UpdatesConsumerMessage.SubscribeToSeq
        case SubscribeToWeak(Some(group)) ⇒
          consumer ! UpdatesConsumerMessage.SubscribeToWeak(Some(group))
        case SubscribeToWeak(None) ⇒
          log.error("Subscribe to weak is done implicitly on UpdatesConsumer start")
      }
    case OnComplete ⇒
      context.stop(self)
    case OnError(cause) ⇒
      log.error(cause, "Error in upstream")
  }
*/

func init() {
	flag.Set("alsologtostderr", "true")
	flag.Set("log_dir", "false")
}

func main() {
	flag.Parse()

	instance := server.NewSessionServer()
	app.DoMainAppInstance(instance)
}
