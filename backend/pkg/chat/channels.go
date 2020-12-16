package chat

type Channels struct {
	Items []Channel
}

func NewChannels(channels []Channel) Channels {
	return Channels{
		Items: channels,
	}
}
