package model

import (
	"fmt"
	"strconv"
)

type MembershipStatus int

const (
	ApprovedMembershipStatus = iota
	PendingStatus
	PendingGroupMembershipStatus
	PendingUserMembershipStatus
)

func AnyMembershipStatus() *MembershipStatus {
	return nil
}

func ParseMembershipStatus(str string) (MembershipStatus, error) {
	i, err := strconv.Atoi(str)
	if err != nil {
		return 0, fmt.Errorf("cannot parse MembershipStatus: %s", err.Error())
	}
	if i < 0 || i > PendingUserMembershipStatus {
		return 0, fmt.Errorf("cannot parse MembershipStatus")
	}
	return MembershipStatus(i), nil
}
