package service

import (
	"fmt"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
)

// GetResourceLink Gets the markdown representing the link to a resource
func (c ChatService) GetResourceLink(resource resourcemodel.ResourceKey) string {
	return fmt.Sprintf("<commonpool-resource id='%s'><commonpool-resource>", resource.String())
}
