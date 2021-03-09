package domain

import (
	"github.com/commonpool/backend/pkg/keys"
	"github.com/stretchr/testify/assert"
	"testing"
)

var groupInfo GroupInfo = GroupInfo{
	Name:        "name",
	Description: "description",
}

func TestGroup(t *testing.T) {

	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	assert.Len(t, group.changes, 2)
	assert.IsType(t, GroupCreated{}, group.changes[0])
	evt1 := group.changes[0].(GroupCreated)
	assert.Equal(t, groupInfo, evt1.GroupInfo)
	assert.Equal(t, 1, evt1.EventVersion)
	assert.Equal(t, GroupCreatedEvent, evt1.EventType)
	assert.Equal(t, owner, evt1.CreatedBy)

	assert.IsType(t, MembershipStatusChanged{}, group.changes[1])
	evt2 := group.changes[1].(MembershipStatusChanged)
	assert.Equal(t, MembershipStatusChangedEvent, evt2.EventType)
	assert.Equal(t, 1, evt2.EventVersion)
	assert.Equal(t, None, evt2.OldPermissions)
	assert.Equal(t, Owner, evt2.NewPermissions)
	assert.Nil(t, evt2.OldStatus)
	assert.Equal(t, ApprovedMembershipStatus, *evt2.NewStatus)
	assert.Equal(t, owner, evt2.MemberKey)
	assert.Equal(t, true, evt2.IsNewMembership)
	assert.Equal(t, owner, evt2.ChangedBy)

	if m, ok := group.GetMembership(owner); assert.True(t, ok) {
		assert.True(t, m.IsOwner())
		assert.True(t, m.IsAdmin())
		assert.True(t, m.IsMember())
	} else {
		return
	}

}

func TestUserInvitedByMember(t *testing.T) {

	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}
	group.MarkAsCommitted()

	user1 := keys.NewUserKey("user1")
	if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
		return
	}
	m, ok := group.GetMembership(user1)
	if !assert.True(t, ok) {
		return
	}

	assert.False(t, m.IsOwner())
	assert.False(t, m.IsAdmin())
	assert.False(t, m.IsMember())

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	assert.False(t, m.IsOwner())
	assert.False(t, m.IsAdmin())
	assert.True(t, m.IsMember())

}

func TestMemberAskedToJoin(t *testing.T) {

	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}
	group.MarkAsCommitted()

	user1 := keys.NewUserKey("user1")

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}
	if m, ok := group.GetMembership(user1); assert.True(t, ok) {
		assert.False(t, m.IsOwner())
		assert.False(t, m.IsAdmin())
		assert.False(t, m.IsMember())
	} else {
		return
	}

	if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
		return
	}

	if m, ok := group.GetMembership(user1); assert.True(t, ok) {
		assert.False(t, m.IsOwner())
		assert.False(t, m.IsAdmin())
		assert.True(t, m.IsMember())
	} else {
		return
	}

}

func TestJoinGroupIdempotent(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}
	group.MarkAsCommitted()

	user1 := keys.NewUserKey("user1")
	for i := 0; i < 5; i++ {
		if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
			return
		}
	}
	assert.Len(t, group.changes, 1)

	for i := 0; i < 5; i++ {
		if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
			return
		}
	}
	assert.Len(t, group.changes, 2)

	for i := 0; i < 5; i++ {
		if err := group.JoinGroup(owner, owner); !assert.NoError(t, err) {
			return
		}
	}
	assert.Len(t, group.changes, 2)
}

func TestJoinGroupShouldFailIfUserIsNotAdmin(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}
	user1 := keys.NewUserKey("user1")
	user2 := keys.NewUserKey("user2")

	err := group.JoinGroup(user2, user1)
	if !assert.Error(t, err) {
		return
	}

}

func TestOwnerAssignsPermission(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	user1 := keys.NewUserKey("user1")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(owner, user1, Admin); !assert.NoError(t, err) {
		return
	}

	m, _ := group.GetMembership(user1)
	assert.Equal(t, Admin, m.PermissionLevel)

}

func TestAdminAssignsPermission(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	admin := keys.NewUserKey("admin")
	user1 := keys.NewUserKey("user1")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, admin); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(admin, admin); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(owner, admin, Admin); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(admin, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(admin, user1, Admin); assert.NoError(t, err) {
		m, _ := group.GetMembership(user1)
		assert.Equal(t, Admin, m.PermissionLevel)
	} else {
		return
	}

	if err := group.AssignPermission(admin, user1, Member); assert.NoError(t, err) {
		m, _ := group.GetMembership(user1)
		assert.Equal(t, Member, m.PermissionLevel)
	} else {
		return
	}

}

func TestAdminGrantsOwnerPermissionsShouldFail(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	admin := keys.NewUserKey("admin")
	user1 := keys.NewUserKey("user1")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, admin); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(admin, admin); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(owner, admin, Admin); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(admin, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(admin, user1, Owner); !assert.Error(t, err) {
		return
	}

}

func TestMemberGrantsPermissionsShouldFail(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	user1 := keys.NewUserKey("user1")
	user2 := keys.NewUserKey("user2")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, user2); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user2, user2); !assert.NoError(t, err) {
		return
	}

	if err := group.AssignPermission(user1, user2, Admin); !assert.Error(t, err) {
		return
	}

	if err := group.AssignPermission(user1, user2, Owner); !assert.Error(t, err) {
		return
	}

}

func TestLeaveGroup(t *testing.T) {
	gk := keys.GenerateGroupKey()
	group := NewGroup(gk)
	owner := keys.NewUserKey("owner")
	user1 := keys.NewUserKey("user1")
	user2 := keys.NewUserKey("user2")

	if err := group.CreateGroup(owner, groupInfo); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user1, user1); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(owner, user2); !assert.NoError(t, err) {
		return
	}

	if err := group.JoinGroup(user2, user2); !assert.NoError(t, err) {
		return
	}

	if err := group.CancelMembership(owner, user1); !assert.NoError(t, err) {
		return
	}

	if _, ok := group.GetMembership(user1); !assert.False(t, ok) {
		return
	}

}
