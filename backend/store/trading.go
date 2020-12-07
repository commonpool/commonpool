package store

import (
	"context"
	"fmt"
	errs "github.com/commonpool/backend/errors"
	"github.com/commonpool/backend/graph"
	"github.com/commonpool/backend/group"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/resource"
	"github.com/commonpool/backend/trading"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"strconv"
	"strings"
	"time"
)

const (
	CompletedKey          = "completed"
	TypeKey               = "type"
	AmountKey             = "amount"
	FromApprovedKey       = "from_approved"
	ToApprovedKey         = "to_approved"
	CreditsTransferredKey = "credits_transferred"
	GivenKey              = "given"
	TakenKey              = "taken"
	ReturnedBackKey       = "returned_back"
	ReceivedBacKey        = "received_back"
	ReceivedKey           = "received"
	DurationKey           = "duration"
)

type TradingStore struct {
	graphDriver graph.GraphDriver
}

var _ trading.Store = TradingStore{}

func NewTradingStore(graphDriver graph.GraphDriver) *TradingStore {
	return &TradingStore{
		graphDriver: graphDriver,
	}
}

func (t TradingStore) MarkOfferItemsAsAccepted(ctx context.Context, approvedByGiver *model.OfferItemKeys, approvedByReceiver *model.OfferItemKeys) error {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if len(approvedByReceiver.Items) > 0 {
		result, err := session.Run(`
			MATCH 
			(offerItem:OfferItem)
			WHERE 
			offerItem.id in $ids
			SET offerItem += {`+ToApprovedKey+`: true}`,
			map[string]interface{}{
				"ids": approvedByReceiver.Strings(),
			})
		if err != nil {
			return err
		}
		if result.Err() != nil {
			return result.Err()
		}
	}

	if len(approvedByGiver.Items) > 0 {
		result, err := session.Run(`
			MATCH 
			(offerItem:OfferItem)
			WHERE 
			offerItem.id in $ids
			SET offerItem += {`+FromApprovedKey+`: true}`,
			map[string]interface{}{
				"ids": approvedByGiver.Strings(),
			})
		if err != nil {
			return err
		}
		if result.Err() != nil {
			return result.Err()
		}

	}

	return nil
}

func (t TradingStore) FindReceivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

		MATCH (offerItem {id:$id})
		WITH offerItem

		MATCH (offerItem)-[:To]->(user:User)
		RETURN user
		
		UNION

		MATCH (offerItem)-[:To]->(:Group)<-[membership:IsMemberOf]-(user:User)
		WHERE (membership.isAdmin = true) OR (membership.isManager = true) 

		WITH user

		return DISTINCT user.id as userId`,
		map[string]interface{}{
			"id": offerItemKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var userKeys []model.UserKey
	for result.Next() {
		record := result.Record()
		userIdField, _ := record.Get("userId")
		userId := userIdField.(string)
		userKey := model.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}
	return model.NewUserKeys(userKeys), nil
}

func (t TradingStore) FindGivingApproversForOfferItem(offerItemKey model.OfferItemKey) (*model.UserKeys, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

		MATCH (offerItem {id:$id})
		WITH offerItem

		MATCH (offerItem)<-[:From]-(:Resource)<-[:Manages]-(:Group)<-[membership:IsMemberOf]-(u:User)
		WHERE (membership.isAdmin = true) OR (membership.isManager = true) 
		RETURN user
		
		UNION 

		MATCH (offerItem)<-[:From]-(:Resource)<-[:Manages]-(user:User)
		RETURN user
		
		UNION		

		MATCH (offerItem)<-[:From]-(r:Resource)-[:CreatedBy]->(user:User)
		RETURN user

		UNION

		MATCH (offerItem)<-[:From]-(user:User)
		RETURN user
		
		UNION

		MATCH (offerItem)<-[:From]-(:Group)<-[membership:IsMemberOf]-(user:User)
		WHERE (membership.isAdmin = true) OR (membership.isManager = true) 
		RETURN user

		WITH user

		return DISTINCT user.id as userId`,
		map[string]interface{}{
			"id": offerItemKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var userKeys []model.UserKey
	for result.Next() {
		record := result.Record()
		userIdField, _ := record.Get("userId")
		userId := userIdField.(string)
		userKey := model.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}
	return model.NewUserKeys(userKeys), nil
}

func (t TradingStore) FindApproversForOffer(offerKey model.OfferKey) (*trading.OfferApprovers, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		
		// Say hello to my little friend
	
		MATCH (o:Offer {id:$offerId})<-[:IsPartOf]-(oi:OfferItem)
		WITH o, oi

		// Retrieving people who can approve giving something
		
		// Matching admins of groups that can manage this resource
		// The admins of groups that can manage a resource are allowed to approve
		// giving the resource.

		OPTIONAL MATCH (oi)<-[:From]-(r:Resource)<-[:Manages]-(:Group)<-[m:IsMemberOf]-(u:User)
		WHERE (m.isAdmin = true) OR (m.isManager) 
		WITH o, oi, collect(u) as from_approvers, collect(r) as resources
		
		// Matching users that can manage this resource
		// Users that are allowed to manage a resource are allowed
		// to approve giving that resource

		OPTIONAL MATCH (oi)<-[:From]-(r:Resource)<-[:Manages]-(u:User)
		WITH o, oi, (from_approvers + collect(u)) as from_approvers, (collect(r) + resources) as resources
		
		// Matching users that created this resource
		// Users that created the resource are allowed to 
		// approve giving the resource

		OPTIONAL MATCH (oi)<-[:From]-(r:Resource)-[:CreatedBy]->(u:User)
		WITH o, oi, (collect(u) + from_approvers) as from_approvers, (collect(r) + resources) as resources
		
		// Matching all people who could approve receiving something
		
		// Matching users that would give time
		// The users who would be receiving time in an offer
		// are allowed to approve receiving time

		OPTIONAL MATCH (oi)<-[:From]-(u:User)
		WITH o, oi, (collect(u) + from_approvers) as from_approvers, resources
		
		// matching admins/managers of groups that would give time
		// Administrators/Managers of groups that would receive time are
		// allowed to approve the group receiving that time

		OPTIONAL MATCH (oi)<-[:From]-(g:Group)<-[m:IsMemberOf]-(u:User)
		WHERE m.isAdmin = true OR m.isManager = true
		WITH o, oi, apoc.coll.toSet(from_approvers) as from_approvers, resources
		
		// Matching people who would receive something
		// People who would receive something as part of an offer  
		// are allowed to approve receiving that thing

		OPTIONAL MATCH (oi)-[:To]->(u:User)
		WITH o, oi, from_approvers, collect(u) as to_approvers, resources
		
		// Matching admins/managers of groups who would receice something
		// Administrators of groups that would receive something 
		// as part of an offer are allowed to approve the group receiving that thing

		OPTIONAL MATCH (oi)-[:To]->(:Group)<-[m:IsMemberOf]-(u:User)
		WHERE (m.isAdmin = true) or (m.isManager = true)
		WITH o, oi, from_approvers, (collect(u) + to_approvers) as to_approvers, resources
		
		RETURN o, oi, from_approvers, to_approvers, apoc.coll.toSet(resources) as resources`,
		map[string]interface{}{
			"offerId": offerKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	fromUserApproversMap := map[model.UserKey][]model.OfferItemKey{}
	toUserApproversMap := map[model.UserKey][]model.OfferItemKey{}
	fromItemApproversMap := map[model.OfferItemKey][]model.UserKey{}
	toItemApproversMap := map[model.OfferItemKey][]model.UserKey{}

	for result.Next() {

		record := result.Record()
		offerItemField, _ := record.Get("oi")
		offerItemNode := offerItemField.(neo4j.Node)
		offerItemId := offerItemNode.Props()["id"].(string)
		offerItemKey, err := model.ParseOfferItemKey(offerItemId)
		if err != nil {
			return nil, err
		}

		fromApproversField, _ := record.Get("from_approvers")
		fromApproversSlice := fromApproversField.([]interface{})

		toApproversField, _ := record.Get("to_approvers")
		toApproversSlice := toApproversField.([]interface{})

		for _, fromApproverIntf := range fromApproversSlice {
			fromApproverNode := fromApproverIntf.(neo4j.Node)
			fromApproverId := fromApproverNode.Props()["id"].(string)
			fromApproverKey := model.NewUserKey(fromApproverId)
			fromUserApproversMap[fromApproverKey] = append(fromUserApproversMap[fromApproverKey], offerItemKey)
			fromItemApproversMap[offerItemKey] = append(fromItemApproversMap[offerItemKey], fromApproverKey)
		}

		for _, toApproverIntf := range toApproversSlice {
			toApproverNode := toApproverIntf.(neo4j.Node)
			toApproverId := toApproverNode.Props()["id"].(string)
			toApproverKey := model.NewUserKey(toApproverId)
			toUserApproversMap[toApproverKey] = append(fromUserApproversMap[toApproverKey], offerItemKey)
			toItemApproversMap[offerItemKey] = append(fromItemApproversMap[offerItemKey], toApproverKey)
		}

	}

	userFromApprovers := map[model.UserKey]*model.OfferItemKeys{}
	for userKey, offerItemKeys := range fromUserApproversMap {
		userFromApprovers[userKey] = model.NewOfferItemKeys(offerItemKeys)
	}
	userToApprovers := map[model.UserKey]*model.OfferItemKeys{}
	for userKey, offerItemKeys := range toUserApproversMap {
		userToApprovers[userKey] = model.NewOfferItemKeys(offerItemKeys)
	}
	itemFromApprovers := map[model.OfferItemKey]*model.UserKeys{}
	for offerItemKey, userKeys := range fromItemApproversMap {
		itemFromApprovers[offerItemKey] = model.NewUserKeys(userKeys)
	}
	itemToApprovers := map[model.OfferItemKey]*model.UserKeys{}
	for offerItemKey, userKeys := range toItemApproversMap {
		itemToApprovers[offerItemKey] = model.NewUserKeys(userKeys)
	}

	return &trading.OfferApprovers{
		OfferItemsUsersCanGive:    userFromApprovers,
		OfferItemsUsersCanReceive: userToApprovers,
		UsersAbleToGiveItem:       itemFromApprovers,
		UsersAbleToReceiveItem:    itemToApprovers,
	}, nil

}

func (t TradingStore) SaveOffer(offer *trading.Offer, offerItems *trading.OfferItems) error {

	/**

	Expected resulting query :

	MATCH

	(user1:User {id:'user1'}),
	(user2:User {id:'user2'}),
	(group1:Group {id:'group1'}),
	(resource1:Resource {id:'resource1'})

	CREATE

	(offer:Offer {id:'offer1'}),

	// When the offer item is about credits being transferred from one person to another
	// then we create an offer item node, and 2 relationships between the users
	//
	// example: (receivingUser)<-[:To]-(offerItem)-[:From]-(givingUser)

	(offerItem1:OfferItem {id:'offerItem1', type:'credits_transfer', seconds:6000})-[:IsPartOf]-(offer),
	(user1)<-[:From]-(offerItem1),
	(user2)<-[:To]-(offerItem1),

	// When the offer item is about a resource being given or borrowed,
	// We create an offer item with 2 relationships, between the receiving
	// user and the resource being given
	//
	// example: (receivingUser)<-[:To]-(offerItem)-(resource)

	(offerItem2:OfferItem {id:'offerItem2', type:'resource_transfer'}-[:IsPartOf]-(offer),
	(resource1)<-[:From]-(offerItem3)
	(user2)<-[:To]-(offerItem2),

	(offerItem3:OfferItem {id:'offerItem3', type:'borrow_resource', starting:'2020-04-15 08:00:00', duration:'5d'}-[:IsPartOf]-(offer),
	(resource1)<-[:From]-(offerItem3)
	(user2)<-[:To]-(offerItem2),

	// When the offer item is about services being given, we create an offer item node
	// and a relationship between the receiving user and the 'service' resource

	(offerItem4:OfferItem {id:'offerItem4', type:'provide_service', duration:'10h'}-[:IsPartOf]-(offer),
	(resource1)<-[:From]-(offerItem3)
	(user2)<-[:To]-(offerItem2),

	// If the recipient or the giver of a time credits / resources is a group, the receiving or giving end
	// is a group node

	(offerItem5:OfferItem {id:'offerItem5', type:'credits_transfer', seconds:6000})-[:IsPartOf]-(offer),
	(group1)<-[:From]-(offerItem1),
	(user3)<-[:To]-(offerItem1),

	*/

	var matchClauses []string
	var createClauses []string
	var params = map[string]interface{}{}

	var userRefMap = map[model.UserKey]string{}
	var groupRefMap = map[model.GroupKey]string{}
	var resourceRefMap = map[model.ResourceKey]string{}

	for i, userKey := range offerItems.GetUserKeys().Items {
		userStr := "user" + strconv.Itoa(i+1)
		userRefMap[userKey] = userStr
		userIdParamName := userStr + "_id"
		params[userIdParamName] = userKey.String()
		matchClauses = append(matchClauses, "("+userStr+":User {id:$"+userIdParamName+"})")
	}
	for i, groupKey := range offerItems.GetGroupKeys().Items {
		groupStr := "group" + strconv.Itoa(i+1)
		groupRefMap[groupKey] = groupStr
		groupIdParamName := groupStr + "_id"
		params[groupIdParamName] = groupKey.String()
		matchClauses = append(matchClauses, "("+groupStr+":Group {id:"+groupIdParamName+"})")
	}
	for i, resourceKey := range offerItems.GetResourceKeys().Items {
		resourceStr := "resource" + strconv.Itoa(i+1)
		resourceRefMap[resourceKey] = resourceStr
		resourceIdParamName := resourceStr + "_id"
		params[resourceIdParamName] = resourceKey.String()
		matchClauses = append(matchClauses, "("+resourceStr+":Resource {id:$"+resourceIdParamName+"})")
	}

	matchClauses = append(matchClauses, "(createdBy:User{id:$created_by_id})")
	params["created_by_id"] = offer.CreatedByKey.String()

	// Add the CREATE (offer:Offer {id:$offer_id}) clause
	createClauses = append(createClauses, `
		(offer:Offer 
			{
				id         :$offer_id,
				status     :$status,
				created_at :$created_at
			}
		)-[:CreatedBy]->(createdBy)`)

	params["offer_id"] = offer.GetKey().String()
	params["created_at"] = offer.CreatedAt
	statusStr, err := statusToString(offer.Status)
	if err != nil {
		return err
	}
	params["status"] = statusStr

	// Loop through each item
	for i, offerItem := range offerItems.Items {

		// The unique identifier for that offer item in the query
		offerItemRef := "offerItem" + strconv.Itoa(i+1)

		// add the type-specific params to the param map
		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*trading.CreditTransferItem)
			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem 
				{
					id:                        $`+offerItemRef+`_id,
					`+TypeKey+`:               $`+offerItemRef+`_type,
					`+AmountKey+`:             $`+offerItemRef+`_amount,
					`+FromApprovedKey+`:       $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:         $`+offerItemRef+`_to_approved,
					`+CreditsTransferredKey+`: $`+offerItemRef+`_credits_transferred,
					`+CompletedKey+`:          $`+offerItemRef+`_completed
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = offerItem.GetKey().String()
			params[offerItemRef+"_type"] = offerItem.Type()
			params[offerItemRef+"_amount"] = int(creditTransfer.Amount.Seconds())
			params[offerItemRef+"_from_approved"] = creditTransfer.GiverAccepted
			params[offerItemRef+"_to_approved"] = creditTransfer.ReceiverAccepted
			params[offerItemRef+"_credits_transferred"] = creditTransfer.CreditsTransferred
			params[offerItemRef+"_completed"] = creditTransfer.IsCompleted()

			fromStr := ""
			toStr := ""

			// Creating the TO part in (<TO>)<-[:To]-(<OFFER_ITEM>)
			if creditTransfer.From.IsForGroup() {
				fromStr = "(" + groupRefMap[creditTransfer.From.GetGroupKey()] + ")"
			} else {
				fromStr = "(" + userRefMap[creditTransfer.From.GetUserKey()] + ")"
			}

			if creditTransfer.To.IsForGroup() {
				toStr = "(" + groupRefMap[creditTransfer.To.GetGroupKey()] + ")"
			} else {
				toStr = "(" + userRefMap[creditTransfer.To.GetUserKey()] + ")"
			}

			createClauses = append(createClauses, "("+offerItemRef+")-[:To]->"+toStr)
			createClauses = append(createClauses, "("+offerItemRef+")<-[:From]-"+fromStr)

		} else if offerItem.IsServiceProviding() {

			provideService := offerItem.(*trading.ProvideServiceItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id
					`+TypeKey+`:         $`+offerItemRef+`_type
					`+DurationKey+`:     $`+offerItemRef+`_duration
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved
					`+GivenKey+`:        $`+offerItemRef+`_given
					`+ReceivedKey+`:     $`+offerItemRef+`_received
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = provideService.Key.String()
			params[offerItemRef+"_type"] = string(trading.ProvideService)
			params[offerItemRef+"_duration"] = int(provideService.Duration.Seconds())
			params[offerItemRef+"_from_approved"] = provideService.GiverAccepted
			params[offerItemRef+"_to_approved"] = provideService.ReceiverAccepted
			params[offerItemRef+"_given"] = provideService.ServiceGivenConfirmation
			params[offerItemRef+"_received"] = provideService.ServiceReceivedConfirmation

			fromStr := ""
			toStr := ""

			if provideService.To.IsForGroup() {
				toStr = "(" + groupRefMap[provideService.To.GetGroupKey()] + ")"
			} else {
				toStr = "(" + userRefMap[provideService.To.GetUserKey()] + ")"
			}

			// Creating the TO part in (<TO>)<-[:To]-(<OFFER_ITEM>)
			fromStr = "(" + resourceRefMap[provideService.ResourceKey] + ")"

			createClauses = append(createClauses, "("+offerItemRef+")-[:To]->"+toStr)
			createClauses = append(createClauses, "("+offerItemRef+")<-[:From]-"+fromStr)

		} else if offerItem.IsBorrowingResource() {

			borrowResource := offerItem.(*trading.BorrowResourceItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id
					`+TypeKey+`:         $`+offerItemRef+`_type
					`+DurationKey+`:     $`+offerItemRef+`_duration
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved
					`+GivenKey+`:        $`+offerItemRef+`_given
					`+TakenKey+`:        $`+offerItemRef+`_taken
					`+ReturnedBackKey+`: $`+offerItemRef+`_returned_back
					`+ReceivedBacKey+`:  $`+offerItemRef+`_received_back,
					`+CompletedKey+`:    $`+offerItemRef+`_completed
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = borrowResource.Key.String()
			params[offerItemRef+"_type"] = string(trading.BorrowResource)
			params[offerItemRef+"_duration"] = int(borrowResource.Duration.Seconds())
			params[offerItemRef+"_from_approved"] = borrowResource.GiverAccepted
			params[offerItemRef+"_to_approved"] = borrowResource.ReceiverAccepted
			params[offerItemRef+"_given"] = borrowResource.ItemGiven
			params[offerItemRef+"_taken"] = borrowResource.ItemTaken
			params[offerItemRef+"_returned_back"] = borrowResource.ItemReturnedBack
			params[offerItemRef+"_received_back"] = borrowResource.ItemReceivedBack
			params[offerItemRef+"_completed"] = borrowResource.IsCompleted()

			fromStr := ""
			toStr := ""

			if borrowResource.To.IsForGroup() {
				toStr = "(" + groupRefMap[borrowResource.To.GetGroupKey()] + ")"
			} else {
				toStr = "(" + userRefMap[borrowResource.To.GetUserKey()] + ")"
			}

			// Creating the TO part in (<TO>)<-[:To]-(<OFFER_ITEM>)
			fromStr = "(" + resourceRefMap[borrowResource.ResourceKey] + ")"

			createClauses = append(createClauses, "("+offerItemRef+")-[:To]->"+toStr)
			createClauses = append(createClauses, "("+offerItemRef+")<-[:From]-"+fromStr)

		} else if offerItem.IsResourceTransfer() {

			resourceTransfer := offerItem.(*trading.ResourceTransferItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id,
					`+TypeKey+`:         $`+offerItemRef+`_type,
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved,
					`+GivenKey+`:        $`+offerItemRef+`_given,
					`+ReceivedKey+`:     $`+offerItemRef+`_received,
					`+CompletedKey+`:    $`+offerItemRef+`_completed
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = resourceTransfer.Key.String()
			params[offerItemRef+"_type"] = string(trading.ResourceTransfer)
			params[offerItemRef+"_from_approved"] = resourceTransfer.GiverAccepted
			params[offerItemRef+"_to_approved"] = resourceTransfer.ReceiverAccepted
			params[offerItemRef+"_given"] = resourceTransfer.ItemGiven
			params[offerItemRef+"_received"] = resourceTransfer.ItemReceived
			params[offerItemRef+"_completed"] = resourceTransfer.IsCompleted()

			fromStr := ""
			toStr := ""

			if resourceTransfer.To.IsForGroup() {
				toStr = "(" + groupRefMap[resourceTransfer.To.GetGroupKey()] + ")"
			} else {
				toStr = "(" + userRefMap[resourceTransfer.To.GetUserKey()] + ")"
			}

			// Creating the TO part in (<TO>)<-[:To]-(<OFFER_ITEM>)
			fromStr = "(" + resourceRefMap[resourceTransfer.ResourceKey] + ")"

			createClauses = append(createClauses, "("+offerItemRef+")-[:To]->"+toStr)
			createClauses = append(createClauses, "("+offerItemRef+")<-[:From]-"+fromStr)

		}

	}

	// Build the cypher

	cypher := "MATCH\n"
	cypher = cypher + strings.Join(matchClauses, ",\n") + "\n"
	cypher = cypher + "CREATE\n"
	cypher = cypher + strings.Join(createClauses, ",\n") + "\n"

	// Execute that baby

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}

	result, err := session.Run(cypher, params)
	if err != nil {
		return err
	}

	if result.Err() != nil {
		return result.Err()
	}

	return nil

}

func statusToString(offerStatus trading.OfferStatus) (string, error) {
	if offerStatus == trading.CompletedOffer {
		return "completed", nil
	} else if offerStatus == trading.DeclinedOffer {
		return "declined", nil
	} else if offerStatus == trading.AcceptedOffer {
		return "accepted", nil
	} else if offerStatus == trading.CanceledOffer {
		return "canceled", nil
	} else if offerStatus == trading.ExpiredOffer {
		return "expired", nil
	} else if offerStatus == trading.PendingOffer {
		return "pending", nil
	} else {
		return "", fmt.Errorf("unknown offer status type")
	}
}

func stringToStatus(offerStatus string) (trading.OfferStatus, error) {
	if offerStatus == "completed" {
		return trading.CompletedOffer, nil
	} else if offerStatus == "declined" {
		return trading.DeclinedOffer, nil
	} else if offerStatus == "accepted" {
		return trading.AcceptedOffer, nil
	} else if offerStatus == "canceled" {
		return trading.CanceledOffer, nil
	} else if offerStatus == "expired" {
		return trading.ExpiredOffer, nil
	} else if offerStatus == "pending" {
		return trading.PendingOffer, nil
	} else {
		return trading.CanceledOffer, fmt.Errorf("unknown offer status type")
	}
}

func (t TradingStore) SaveOfferStatus(key model.OfferKey, status trading.OfferStatus) error {
	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	statusString, err := statusToString(status)
	if err != nil {
		return err
	}
	result, err := session.Run(`MATCH (o:Offer {id:$id}) SET o += {status: $status} return o`, map[string]interface{}{
		"id":     key.String(),
		"status": statusString,
	})
	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	if !result.Next() {
		return errs.ErrOfferNotFound
	}
	return nil
}

func (t TradingStore) GetOfferItem(ctx context.Context, key model.OfferItemKey) (trading.OfferItem2, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (offer:Offer)<-[:IsPartOf]-(o:OfferItem {id:$id})-[:To]->(to)
		OPTIONAL MATCH (o)<-[:From]-(from)
		RETURN offer.id as offerId, o, from, to`,
		map[string]interface{}{
			"id": key.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	if !result.Next() {
		return nil, errs.ErrOfferItemNotFound
	}

	offerIdField, _ := result.Record().Get("offerId")
	offerId := offerIdField.(string)
	offerKey, err := model.ParseOfferKey(offerId)
	if err != nil {
		return nil, err
	}

	offerField, _ := result.Record().Get("o")
	offerNode := offerField.(neo4j.Node)

	fromField, _ := result.Record().Get("from")
	fromNode := fromField.(neo4j.Node)

	toField, _ := result.Record().Get("to")
	toNode := toField.(neo4j.Node)

	return MapOfferItem(offerKey, offerNode, fromNode, toNode)

}

func MapOfferItemTarget(node neo4j.Node) (*trading.OfferItemTarget, error) {
	if node == nil {
		return nil, fmt.Errorf("node is nil")
	}
	isGroup := IsGroupNode(node)
	isUser := !isGroup && IsUserNode(node)
	if !isGroup && !isUser {
		return nil, fmt.Errorf("target is neither user nor group")
	}

	if isGroup {
		groupKey, err := group.ParseGroupKey(node.Props()["id"].(string))
		if err != nil {
			return nil, err
		}
		return &trading.OfferItemTarget{
			GroupKey: &groupKey,
			Type:     trading.GroupTarget,
		}, nil
	}
	userKey := model.NewUserKey(node.Props()["id"].(string))
	return &trading.OfferItemTarget{
		UserKey: &userKey,
		Type:    trading.UserTarget,
	}, nil
}

func MapOfferItem(offerKey model.OfferKey, offerItemNode neo4j.Node, fromNode neo4j.Node, toNode neo4j.Node) (trading.OfferItem2, error) {

	offerItemType := offerItemNode.Props()["type"].(string)
	var fromResource *resource.Resource
	var fromTarget *trading.OfferItemTarget
	var toTarget *trading.OfferItemTarget
	var err error

	if fromNode != nil {
		if IsGroupNode(fromNode) || IsUserNode(fromNode) {
			fromTarget, err = MapOfferItemTarget(fromNode)
			if err != nil {
				return nil, err
			}
		} else if IsResourceNode(fromNode) {
			fromResource, err = MapResourceNode(fromNode)
			if err != nil {
				return nil, err
			}
		}
	}
	if toNode != nil {
		if IsGroupNode(toNode) || IsUserNode(toNode) {
			toTarget, err = MapOfferItemTarget(toNode)
			if err != nil {
				return nil, err
			}
		}

	}

	offerItemId := offerItemNode.Props()["id"].(string)
	offerItemKey, err := model.ParseOfferItemKey(offerItemId)
	if err != nil {
		return nil, err
	}

	offerItemBase := trading.OfferItemBase{
		Type:             trading.OfferItemType2(offerItemType),
		Key:              offerItemKey,
		OfferKey:         offerKey,
		To:               toTarget,
		ReceiverAccepted: offerItemNode.Props()[ToApprovedKey].(bool),
		GiverAccepted:    offerItemNode.Props()[FromApprovedKey].(bool),
	}

	if offerItemType == string(trading.CreditTransfer) {

		if toTarget == nil {
			return nil, fmt.Errorf("result should have a 'To' user/group")
		}

		if fromTarget == nil {
			return nil, fmt.Errorf("result should have a 'From' user/group")
		}

		amount, ok := offerItemNode.Props()["amount"]
		if !ok {
			return nil, fmt.Errorf("result should have an 'amount' prop")
		}

		return &trading.CreditTransferItem{
			OfferItemBase:      offerItemBase,
			From:               fromTarget,
			Amount:             time.Duration(int64(time.Second) * amount.(int64)),
			CreditsTransferred: offerItemNode.Props()[CreditsTransferredKey].(bool),
		}, nil

	} else if offerItemType == string(trading.ProvideService) {

		if fromResource == nil {
			return nil, fmt.Errorf("result should have a 'From' resource")
		}

		if toTarget == nil {
			return nil, fmt.Errorf("result should have a 'To' user/group")
		}

		duration, ok := offerItemNode.Props()["duration"]
		if !ok {
			return nil, fmt.Errorf("result should have a 'duration' prop")
		}

		return &trading.ProvideServiceItem{
			OfferItemBase:               offerItemBase,
			ResourceKey:                 fromResource.Key,
			Duration:                    time.Duration(int64(time.Second) * duration.(int64)),
			ServiceReceivedConfirmation: offerItemNode.Props()[ReceivedKey].(bool),
			ServiceGivenConfirmation:    offerItemNode.Props()[GivenKey].(bool),
		}, nil

	} else if offerItemType == string(trading.BorrowResource) {

		if fromResource == nil {
			return nil, fmt.Errorf("result should have a 'From' resource")
		}

		if toTarget == nil {
			return nil, fmt.Errorf("result should have a 'To' user/group")
		}

		duration, ok := offerItemNode.Props()["duration"]
		if !ok {
			return nil, fmt.Errorf("result should have a 'duration' prop")
		}

		return &trading.BorrowResourceItem{
			OfferItemBase:    offerItemBase,
			ResourceKey:      fromResource.Key,
			Duration:         time.Duration(int64(time.Second) * duration.(int64)),
			ItemTaken:        offerItemNode.Props()[TakenKey].(bool),
			ItemGiven:        offerItemNode.Props()[GivenKey].(bool),
			ItemReturnedBack: offerItemNode.Props()[ReturnedBackKey].(bool),
			ItemReceivedBack: offerItemNode.Props()[ReceivedBacKey].(bool),
		}, nil

	} else if offerItemType == string(trading.ResourceTransfer) {

		if fromResource == nil {
			return nil, fmt.Errorf("result should have a 'From' resource")
		}

		if toTarget == nil {
			return nil, fmt.Errorf("result should have a 'To' user/group")
		}

		return &trading.ResourceTransferItem{
			OfferItemBase: offerItemBase,
			ResourceKey:   fromResource.Key,
			ItemReceived:  offerItemNode.Props()[ReceivedKey].(bool),
			ItemGiven:     offerItemNode.Props()[GivenKey].(bool),
		}, nil

	} else {
		return nil, fmt.Errorf("unexpected offer item type: %s", offerItemType)
	}

}

func (t TradingStore) UpdateOfferItem(ctx context.Context, offerItem trading.OfferItem2) error {
	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var result neo4j.Result

	if offerItem.IsServiceProviding() {

		service := offerItem.(*trading.ProvideServiceItem)
		result, err = session.Run(`
		MATCH (o:OfferItem {id:$id})
		SET o += {
			`+FromApprovedKey+`: $`+FromApprovedKey+`,
			`+ToApprovedKey+`:   $`+ToApprovedKey+`,
			`+GivenKey+`:        $`+GivenKey+`,
			`+ReceivedKey+`:     $`+ReceivedKey+`
		}
		RETURN o`,
			map[string]interface{}{
				FromApprovedKey: service.GiverAccepted,
				ToApprovedKey:   service.ReceiverAccepted,
				GivenKey:        service.ServiceGivenConfirmation,
				ReceivedKey:     service.ServiceReceivedConfirmation,
				CompletedKey:    service.IsCompleted(),
			})

	} else if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(trading.CreditTransferItem)
		result, err = session.Run(`
			MATCH (o:OfferItem {id:$id})
			SET o += {
				`+FromApprovedKey+`:       $`+FromApprovedKey+`,
				`+ToApprovedKey+`:         $`+ToApprovedKey+`,
				`+CreditsTransferredKey+`: $`+CreditsTransferredKey+`,
				`+CompletedKey+`:          $`+CompletedKey+`
			}
			RETURN o`,
			map[string]interface{}{
				FromApprovedKey:       creditTransfer.GiverAccepted,
				ToApprovedKey:         creditTransfer.ReceiverAccepted,
				CreditsTransferredKey: creditTransfer.CreditsTransferred,
				CompletedKey:          creditTransfer.IsCompleted(),
			})

	} else if offerItem.IsBorrowingResource() {
		resourceBorrow := offerItem.(trading.BorrowResourceItem)
		result, err = session.Run(`
			MATCH (o:OfferItem {id:$id})
			SET o += {
					`+FromApprovedKey+`: $`+FromApprovedKey+`,
					`+ToApprovedKey+`:   $`+ToApprovedKey+`,
					`+GivenKey+`:        $`+GivenKey+`,
					`+TakenKey+`:        $`+TakenKey+`,
					`+ReturnedBackKey+`: $`+ReturnedBackKey+`,
					`+ReceivedBacKey+`:  $`+ReceivedBacKey+`,
					`+CompletedKey+`:    $`+CompletedKey+`
			}
			RETURN o`,
			map[string]interface{}{
				FromApprovedKey: resourceBorrow.GiverAccepted,
				ToApprovedKey:   resourceBorrow.ReceiverAccepted,
				GivenKey:        resourceBorrow.ItemGiven,
				TakenKey:        resourceBorrow.ItemTaken,
				ReturnedBackKey: resourceBorrow.ItemReturnedBack,
				ReceivedBacKey:  resourceBorrow.ItemReceivedBack,
				CompletedKey:    resourceBorrow.IsCompleted(),
			})

	} else if offerItem.IsResourceTransfer() {

		service := offerItem.(trading.ResourceTransferItem)
		result, err = session.Run(`
		MATCH (o:OfferItem {id:$id})
		SET o += {
			`+FromApprovedKey+`: $`+FromApprovedKey+`,
			`+ToApprovedKey+`:   $`+ToApprovedKey+`,
			`+GivenKey+`:        $`+GivenKey+`,
			`+ReceivedKey+`:     $`+ReceivedKey+`
		}
		RETURN o`,
			map[string]interface{}{
				FromApprovedKey: service.GiverAccepted,
				ToApprovedKey:   service.ReceiverAccepted,
				GivenKey:        service.ItemGiven,
				ReceivedKey:     service.ItemReceived,
				CompletedKey:    service.IsCompleted(),
			})

	}

	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	if !result.Next() {
		return errs.ErrOfferItemNotFound
	}
	return nil
}

func (t TradingStore) ConfirmItemGiven(ctx context.Context, key model.OfferItemKey) error {
	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	result, err := session.Run(`
		MATCH (o:OfferItem {id:$id})
		SET o += {given: true}
		RETURN o`,
		map[string]interface{}{
			"id": key.String(),
		})

	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}

	if !result.Next() {
		return errs.ErrOfferItemNotFound
	}
	return nil
}

func (t TradingStore) GetOffer(key model.OfferKey) (*trading.Offer, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (offer:Offer {id:$id})-[:CreatedBy]->(createdBy:User) return offer, createdBy.id as createdById`,
		map[string]interface{}{
			"id": key.String(),
		})
	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}
	if !result.Next() {
		return nil, errs.ErrOfferNotFound
	}

	createdByField, _ := result.Record().Get("createdById")
	createdById := createdByField.(string)
	createdByKey := model.NewUserKey(createdById)

	offerField, _ := result.Record().Get("offer")
	offerNode := offerField.(neo4j.Node)

	return MapOfferNode(offerNode, createdByKey)

}

func MapOfferNode(node neo4j.Node, createdByKey model.UserKey) (*trading.Offer, error) {

	offerId := node.Props()["id"].(string)
	offerKey, err := model.ParseOfferKey(offerId)
	if err != nil {
		return nil, err
	}

	status, err := stringToStatus(node.Props()["status"].(string))
	if err != nil {
		return nil, err
	}

	return &trading.Offer{
		Key:          offerKey,
		CreatedByKey: createdByKey,
		Status:       status,
		CreatedAt:    node.Props()["created_at"].(time.Time),
	}, nil

}

func (t TradingStore) GetOfferItemsForOffer(key model.OfferKey) (*trading.OfferItems, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		MATCH (offer:Offer {id:$id})<-[:IsPartOf]-(offerItem:OfferItem)-[:To]->(to)
		OPTIONAL MATCH (offerItem)<-[:From]-(from)
		RETURN offer, offer.id as offerId, offerItem, from, to`,
		map[string]interface{}{
			"id": key.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var offerItems []trading.OfferItem2
	for result.Next() {

		offerIdField, _ := result.Record().Get("offerId")
		offerId := offerIdField.(string)
		offerKey, err := model.ParseOfferKey(offerId)
		if err != nil {
			return nil, err
		}

		offerItemField, _ := result.Record().Get("offerItem")
		offerItemNode := offerItemField.(neo4j.Node)

		fromField, _ := result.Record().Get("from")
		fromNode := fromField.(neo4j.Node)

		toField, _ := result.Record().Get("to")
		toNode := toField.(neo4j.Node)

		offerItem, err := MapOfferItem(offerKey, offerItemNode, fromNode, toNode)
		if err != nil {
			return nil, err
		}

		offerItems = append(offerItems, offerItem)

	}

	return trading.NewOfferItems(offerItems), nil

}

func (t TradingStore) GetOffersForUser(userKey model.UserKey) (*trading.GetOffersResult, error) {

	/**
	ResourceKey *model.ResourceKey
	Status      *OfferStatus
	UserKeys    []model.UserKey
	*/

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

		MATCH (user:User {id:$userId})

		CALL {
		
			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)<-[:From]-(:Resource)<-[:Manages]-(:Group)<-[membership:IsMemberOf]-(user)
			WHERE membership.isAdmin or membership.isManager
			RETURN offer

			UNION

			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)<-[:From]-(:Resource)<-[:Manages|CreatedBy]-(user)
			RETURN offer

			UNION
		
			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)<-[:From]-(user)
			RETURN offer

			UNION

			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)<-[:From]-(:Group)<-[membership:IsMemberOf]-(u:User)
			WHERE membership.isAdmin or membership.isManager
			RETURN offer

			UNION

			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)-[:To]->(user)
			RETURN offer

			UNION

			MATCH (offer:Offer)<-[:IsPartOf]-(:OfferItem)-[:To]->(:Group)<-[membership:IsMemberOf]-(u:User)
			WHERE membership.isAdmin or membership.isManager
			RETURN offer

		}
		
		WITH user, offer
		MATCH (createdBy:User)<-[:CreatedBy]-(offer)<-[:IsPartOf]-(offerItem:OfferItem)-[:To]->(to)

		WITH user, offer, createdBy, to
		OPTIONAL MATCH (from)-[:From]->(offerItem)
		
		RETURN createdBy.id as createdById, offer, collect({offerItem: offerItem, from:from, to:to}) as offerItems`,
		map[string]interface{}{
			"userId": userKey.String(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var resultItems []*trading.GetOffersResultItem

	for result.Next() {

		record := result.Record()
		offerField, _ := record.Get("offer")
		offerNode := offerField.(neo4j.Node)

		createdByField, _ := record.Get("createdById")
		createdById := createdByField.(string)
		createdByKey := model.NewUserKey(createdById)

		offer, err := MapOfferNode(offerNode, createdByKey)
		if err != nil {
			return nil, err
		}

		offerItemsContainerField, _ := record.Get("offerItems")
		offerItemsContainerSlice := offerItemsContainerField.([]interface{})

		var offerItems []trading.OfferItem2
		for _, offerItemContainerIntf := range offerItemsContainerSlice {

			offerItemContainer := offerItemContainerIntf.(map[string]interface{})
			offerItemField := offerItemContainer["offerItem"]
			offerItemNode := offerItemField.(neo4j.Node)
			fromNode := offerItemContainer["from"].(neo4j.Node)
			toNode := offerItemContainer["to"].(neo4j.Node)

			offerItem, err := MapOfferItem(offer.GetKey(), offerItemNode, fromNode, toNode)
			if err != nil {
				return nil, err
			}
			offerItems = append(offerItems, offerItem)
		}

		resultItems = append(resultItems, &trading.GetOffersResultItem{
			Offer:      offer,
			OfferItems: trading.NewOfferItems(offerItems),
		})

	}

	return &trading.GetOffersResult{
		Items: resultItems,
	}, nil

}

func (t TradingStore) GetTradingHistory(ctx context.Context, ids *model.UserKeys) ([]trading.HistoryEntry, error) {
	panic("implement me")
}
