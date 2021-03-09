package domain

import (
	"github.com/commonpool/backend/pkg/keys"
)

//
// ########### Mock Offer Permission Getter ###############
//

type MockApprover struct {
	Approve bool
}

func (m MockApprover) Can(userKey keys.UserKey, offerItem OfferItem, direction ApprovalDirection) bool {
	return m.Approve
}

var _ OfferPermissionGetter = &MockApprover{}

var ApproveAllMatrix = MockApprover{Approve: true}
var DenyAllMatrix = MockApprover{Approve: false}
