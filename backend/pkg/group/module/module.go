package module

import (
	"github.com/commonpool/backend/pkg/graph"
	"github.com/commonpool/backend/pkg/group"
	"github.com/commonpool/backend/pkg/group/domain"
	grouphandler "github.com/commonpool/backend/pkg/group/handler"
	"github.com/commonpool/backend/pkg/mq"
)

type GroupModule struct {
	Service group.Service
	Handler *grouphandler.Handler
	Store   group.Store
	Repo    domain.GroupRepository
}

func NewGroupModule(amqpClient mq.Client, driver graph.Driver) *GroupModule {
	// store := store2.NewGroupStore(driver)
	// service := service2.NewGroupService(store, mq)
	return &GroupModule{}
}

func (m *GroupModule) Register() {

}
