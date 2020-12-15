package service

import (
	"fmt"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
)

func (c ChatService) GetGroupLink(groupKey groupmodel.GroupKey) string {
	return fmt.Sprintf("<commonpool-group id='%s'><commonpool-group>", groupKey.String())
}
