package service

import (
	"fmt"
	"github.com/commonpool/backend/model"
)

// GetResourceLink Gets the markdown representing the link to a resource
func (c ChatService) GetResourceLink(resource model.ResourceKey) string {
	return fmt.Sprintf("<commonpool-resource id='%s'><commonpool-resource>", resource.String())
}
