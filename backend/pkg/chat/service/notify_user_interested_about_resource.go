package service

// NotifyUserInterestedAboutResource will create a channel between two users if it doesn't exist,
// and will send a message to the owner of the resource notifying them that someone is interested.
// func (c ChatService) NotifyUserInterestedAboutResource(ctx context.Context, request *chat.NotifyUserInterestedAboutResource) (*chat.NotifyUserInterestedAboutResourceResponse, error) {
//
// 	loggedInUser, err := auth.GetLoggedInUser(ctx)
// 	if err != nil {
// 		return nil, err
// 	}
// 	loggedInUserKey := loggedInUser.GetUserKey()
//
// 	getResource, err := c.rs.GetByKey(ctx, resource3.NewGetResourceByKeyQuery(request.ResourceKey))
// 	if err != nil {
// 		return nil, err
// 	}
// 	resource := getResource.Resource
// 	resourceOwnerKey := resource.GetOwnerKey()
//
// 	// make sure auth user is not resource owner
// 	// doesn't make sense for one to inquire about his own stuff
// 	if resourceOwnerKey == loggedInUserKey {
// 		err := exceptions.ErrCannotInquireAboutOwnResource()
// 		return nil, err
// 	}
//
// 	userKeys := usermodel.NewUserKeys([]usermodel.MemberKey{loggedInUserKey, resourceOwnerKey})
//
// 	_, err = c.SendConversationMessage(ctx, chat.NewSendConversationMessage(
// 		loggedInUserKey,
// 		loggedInUser.Username,
// 		userKeys,
// 		request.Message,
// 		[]chatmodel.Block{
// 			*chatmodel.NewHeaderBlock(chatmodel.NewMarkdownObject("Someone is interested in your stuff!"), nil),
// 			*chatmodel.NewContextBlock([]chatmodel.BlockElement{
// 				chatmodel.NewMarkdownObject(
// 					fmt.Sprintf("%s is interested by your post %s.",
// 						c.GetUserLink(loggedInUserKey),
// 						c.GetResourceLink(request.ResourceKey),
// 					),
// 				),
// 			}, nil),
// 		},
// 		[]chatmodel.Attachment{},
// 		&resourceOwnerKey,
// 	))
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	sentPublicMessage, err := c.SendConversationMessage(ctx, chat.NewSendConversationMessage(
// 		loggedInUserKey,
// 		loggedInUser.Username,
// 		userKeys,
// 		request.Message,
// 		[]chatmodel.Block{},
// 		[]chatmodel.Attachment{},
// 		nil,
// 	))
//
// 	if err != nil {
// 		return nil, err
// 	}
//
// 	return &chat.NotifyUserInterestedAboutResourceResponse{
// 		ChannelKey: sentPublicMessage.Message.ChannelKey,
// 	}, nil
//
// }
