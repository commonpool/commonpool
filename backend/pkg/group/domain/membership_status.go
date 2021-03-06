package domain

import (
	"fmt"
	"strconv"
)

type MembershipStatus int

const (
	PendingGroupMembershipStatus MembershipStatus = iota
	PendingUserMembershipStatus
	ApprovedMembershipStatus
)

func AnyMembershipStatus() *MembershipStatus {
	return nil
}

func ParseMembershipStatus(str string) (MembershipStatus, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("cannot parse MembershipStatus: %s", err.Error())
	}
	if i < 0 || i > int(PendingUserMembershipStatus) {
		return 0, fmt.Errorf("cannot parse MembershipStatus")
	}
	return MembershipStatus(i), nil
}
