package handler

import (
	"github.com/commonpool/backend/mock"
	"github.com/commonpool/backend/pkg/test"
)

type resourceHandlerSuite struct {
	test.HandlerSuite
	Handler         *ResourceHandler
	ResourceService *mock.ResourceService
	GroupService    *mock.GroupService
	UserService     *mock.UserService
}

func (s *resourceHandlerSuite) SetupTest() {
	s.HandlerSuite.SetupTest()

	s.ResourceService = &mock.ResourceService{}
	s.GroupService = &mock.GroupService{}
	s.UserService = &mock.UserService{}

	s.Handler = &ResourceHandler{
		resourceService: s.ResourceService,
		groupService:    s.GroupService,
		userService:     s.UserService,
		authorization:   s.Authenticator,
	}

	group := s.Echo.Group("/api/v1")
	s.Handler.Register(group)

}
