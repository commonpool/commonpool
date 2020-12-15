package model

type ChannelSubscriptions struct {
	Items []ChannelSubscription
}

func NewChannelSubscriptions(items []ChannelSubscription) *ChannelSubscriptions {
	return &ChannelSubscriptions{
		Items: items,
	}
}
