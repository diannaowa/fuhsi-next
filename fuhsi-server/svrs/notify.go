// Copyright 2019 - now The https://github.com/fuhsicloud/fuhsi-next Authors
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package svrs

import (
	"github.com/fuhsicloud/fuhsi-next/fuhsi-server/daos"
	. "github.com/fuhsicloud/fuhsi-next/fuhsi-server/entities"
	"github.com/fuhsicloud/fuhsi-next/fuhsi-server/libs"
	"github.com/fuhsicloud/fuhsi-next/fuhsi-server/libs/logger"
)

var (
	DefaultNotifySvr = new(NotifySvr)

	// messages waiting to deal
	notifyQueue = make(chan *NotifyItem, 1000)
)

type (
	NotifySvr struct{}

	// notify item in chan queue
	NotifyItem struct {
		Message string   `json:"message"`
		Users   []string `json:"users"`
	}
)

func (t *NotifySvr) CommitNotifyByUid(uid int64, message string) error {
	user := new(UserEntity)
	err := daos.DefaultUserDao.GetById(uid, user)
	if err != nil {
		return err
	}

	t.CommitNotify(&NotifyItem{
		Message: message,
		Users:   []string{user.Username},
	})
	return nil
}

func (t *NotifySvr) CommitNotify(notify *NotifyItem) {
	logger.Debugf("Received notification: %v", notify)
	notifyQueue <- notify
}

// consume from notification queue
func (t *NotifySvr) Consume() {
	logger.Debugf("Notification consumer booted.")
	for {
		notify := <-notifyQueue
		logger.Debugf("Consume notify: %s", libs.JsonStr(notify))
	}

	// @TODO notify to email/IM
}
