package service

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

func (c ChatService) GetGroupLink(groupKey model.GroupKey) string {
	return fmt.Sprintf("<commonpool-group id='%s'><commonpool-group>", groupKey.String())
}
