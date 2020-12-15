package store

import (
	"context"
	"fmt"
	"github.com/commonpool/backend/model"
	"github.com/commonpool/backend/pkg/exceptions"
	graph2 "github.com/commonpool/backend/pkg/graph"
	groupmodel "github.com/commonpool/backend/pkg/group/model"
	groupstore "github.com/commonpool/backend/pkg/group/store"
	resourcemodel "github.com/commonpool/backend/pkg/resource/model"
	resourcestore "github.com/commonpool/backend/pkg/resource/store"
	sharedstore "github.com/commonpool/backend/pkg/shared/store"
	"github.com/commonpool/backend/pkg/trading"
	tradingmodel "github.com/commonpool/backend/pkg/trading/model"
	usermodel "github.com/commonpool/backend/pkg/user/model"
	store2 "github.com/commonpool/backend/pkg/user/store"
	"github.com/neo4j/neo4j-go-driver/neo4j"
	"strconv"
	"strings"
	"time"
)

const (
	CompletedKey          = "completed"
	CreatedAtKey          = "created_at"
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
	AcceptedAtKey         = "accepted_at"
	DeclinedAtKey         = "declined_at"
	CanceledAtKey         = "canceled_at"
	CompletedAtKey        = "completed_at"
	UpdatedAtKey          = "updated_at"
	StatusKey             = "status"
)

type TradingStore struct {
	graphDriver graph2.Driver
}

var _ trading.Store = TradingStore{}

func NewTradingStore(graphDriver graph2.Driver) *TradingStore {
	return &TradingStore{
		graphDriver: graphDriver,
	}
}

func (t TradingStore) MarkOfferItemsAsAccepted(
	ctx context.Context,
	approvedBy usermodel.UserKey,
	approvedByGiver *tradingmodel.OfferItemKeys,
	approvedByReceiver *tradingmodel.OfferItemKeys) error {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	if len(approvedByReceiver.Items) > 0 {
		result, err := session.Run(`
			MATCH (approver:User {id:$userId})
			WITH approver
			MATCH 
			(offerItem:OfferItem)
			WHERE 
			offerItem.id in $ids
			SET offerItem += {`+ToApprovedKey+`: true}
			MERGE (approver)-[:ApprovedReceiving]->(offerItem)
			`,
			map[string]interface{}{
				"ids":    approvedByReceiver.Strings(),
				"userId": approvedBy.String(),
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
			MATCH (approver:User {id:$userId})
			WITH approver
			MATCH 
			(offerItem:OfferItem)
			WHERE 
			offerItem.id in $ids
			SET offerItem += {`+FromApprovedKey+`: true}
			MERGE (approver)-[:ApprovedGiving]->(offerItem)
			`,
			map[string]interface{}{
				"ids":    approvedByGiver.Strings(),
				"userId": approvedBy.String(),
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

func (t TradingStore) FindReceivingApproversForOfferItem(offerItemKey tradingmodel.OfferItemKey) (*usermodel.UserKeys, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

		MATCH (offerItem {id:$id})
		WITH offerItem

		CALL {

			WITH offerItem
			MATCH (offerItem)-[:To]->(user:User)
			RETURN user

			UNION
	
			WITH offerItem
			MATCH (offerItem)-[:To]->(:Group)<-[membership:IsMemberOf]-(user:User)
			WHERE (membership.isAdmin = true) OR (membership.isManager = true) 
			RETURN user
		}

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

	var userKeys []usermodel.UserKey
	for result.Next() {
		record := result.Record()
		userIdField, _ := record.Get("userId")
		userId := userIdField.(string)
		userKey := usermodel.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}
	return usermodel.NewUserKeys(userKeys), nil
}

func (t TradingStore) FindGivingApproversForOfferItem(offerItemKey tradingmodel.OfferItemKey) (*usermodel.UserKeys, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

		MATCH (offerItem {id:$id})
		WITH offerItem

		CALL {
		
			WITH offerItem
			MATCH (offerItem)<-[:From]-(:Resource)<-[:Manages]-(:Group)<-[membership:IsMemberOf]-(user:User)
			WHERE (membership.isAdmin = true) OR (membership.isManager = true) 
			RETURN user
			
			UNION 
	
			WITH offerItem
			MATCH (offerItem)<-[:From]-(:Resource)<-[:Manages]-(user:User)
			RETURN user
			
			UNION		
	
			WITH offerItem
			MATCH (offerItem)<-[:From]-(r:Resource)-[:CreatedBy]->(user:User)
			RETURN user
	
			UNION
	
			WITH offerItem
			MATCH (offerItem)<-[:From]-(user:User)
			RETURN user
			
			UNION
	
			WITH offerItem
			MATCH (offerItem)<-[:From]-(:Group)<-[membership:IsMemberOf]-(user:User)
			WHERE (membership.isAdmin = true) OR (membership.isManager = true) 
			RETURN user
		}

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

	var userKeys []usermodel.UserKey
	for result.Next() {
		record := result.Record()
		userIdField, _ := record.Get("userId")
		userId := userIdField.(string)
		userKey := usermodel.NewUserKey(userId)
		userKeys = append(userKeys, userKey)
	}
	return usermodel.NewUserKeys(userKeys), nil
}

func (t TradingStore) FindApproversForCandidateOffer(offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems) (*usermodel.UserKeys, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	resourceKeys := offerItems.GetResourceKeys()
	userKeys := offerItems.GetUserKeys()
	groupKeys := offerItems.GetGroupKeys()

	result, err := session.Run(`

		CALL {
			MATCH (group:Group)<-[membership:IsMemberOf]-(user:User)
			WHERE membership.isAdmin and group.id in $groupIds
			RETURN user
	
			UNION
	
			MATCH (user:User)
			WHERE user.id in $userIds
			RETURN user
	
			UNION
	
			MATCH (resource:Resource)<-[:CreatedBy]-(user:User)
			WHERE resource.id in $resourceIds
			RETURN user
		}
		
		RETURN DISTINCT user.id as userId`,
		map[string]interface{}{
			"groupIds":    groupKeys.Strings(),
			"userIds":     userKeys.Strings(),
			"resourceIds": resourceKeys.Strings(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var uks []usermodel.UserKey
	for result.Next() {

		userIdField, _ := result.Record().Get("userId")
		userIdStr := userIdField.(string)
		userKey := usermodel.NewUserKey(userIdStr)
		uks = append(uks, userKey)
	}

	return usermodel.NewUserKeys(uks), nil
}

func (t TradingStore) FindApproversForOffers(offerKeys *tradingmodel.OfferKeys) (*tradingmodel.OffersApprovers, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`

			match (offer:Offer)
			where offer.id in $offerIds
			with offer
			match (offerItem:OfferItem)-[:IsPartOf]->(offer)
			with offer, offerItem
			match (offerItem)-[toRel:To]->(to)
			with offer, offerItem, to, toRel
			match (from)-[fromRel:From]->(offerItem)
			with offer, offerItem, to, toRel, from, fromRel
			call {
				with from, fromRel, offerItem
				match (from)-[:CreatedBy]->(user)
				where from:Resource and user:User
				return user as fromApprover
				
				union
				with from, fromRel, offerItem
				match (from)-[fromRel]->(offerItem)
				where from:User
				return from as fromApprover
				
				union
				with from, fromRel, offerItem
				match (from)-[fromRel]->(offerItem),(user:User)-[membership:IsMemberOf]->(from)
				where membership.isAdmin
				return user as fromApprover
				
			}
			
			with offer, offerItem, to, toRel, from, fromRel, collect(distinct fromApprover) as fromApprovers
			
			call {
			
			  with to, toRel, offerItem
			  match (to)<-[toRel]-(offerItem)
			  where to:User
			  return to as toApprover
			  
			  union
			  
			  with to, toRel, offerItem
			  match (user:User)-[membership:IsMemberOf]->(to)<-[toRel]-(offerItem)
			  where to:Group and membership.isAdmin
			  return user as toApprover
			
			}

			with offer, offerItem, from, to, fromApprovers, collect(distinct toApprover) as toApprovers
			return offer, offer.id as offerId, collect({offerItem: offerItem, from: from, to: to, fromApprovers: fromApprovers, toApprovers: toApprovers}) as offerItems
`,
		map[string]interface{}{
			"offerIds": offerKeys.Strings(),
		})

	if err != nil {
		return nil, err
	}
	if result.Err() != nil {
		return nil, result.Err()
	}

	var approversForOffer []*tradingmodel.OfferApprovers

	for result.Next() {

		fromUserApproversMap := map[usermodel.UserKey][]tradingmodel.OfferItemKey{}
		toUserApproversMap := map[usermodel.UserKey][]tradingmodel.OfferItemKey{}
		fromItemApproversMap := map[tradingmodel.OfferItemKey][]usermodel.UserKey{}
		toItemApproversMap := map[tradingmodel.OfferItemKey][]usermodel.UserKey{}

		record := result.Record()

		offerIdField, _ := record.Get("offerId")
		offerId := offerIdField.(string)
		offerKey, err := tradingmodel.ParseOfferKey(offerId)
		if err != nil {
			return nil, err
		}

		structField, _ := record.Get("offerItems")
		structIntfs := structField.([]interface{})

		for _, structIntf := range structIntfs {

			fieldMap := structIntf.(map[string]interface{})

			offerItemField := fieldMap["offerItem"]
			offerItemNode := offerItemField.(neo4j.Node)
			offerItemId := offerItemNode.Props()["id"].(string)
			offerItemKey, err := tradingmodel.ParseOfferItemKey(offerItemId)
			if err != nil {
				return nil, err
			}

			fromApproversField, _ := fieldMap["fromApprovers"]
			fromApproversSlice := fromApproversField.([]interface{})

			toApproversField, _ := fieldMap["toApprovers"]
			toApproversSlice := toApproversField.([]interface{})

			for _, fromApproverIntf := range fromApproversSlice {
				fromApproverNode := fromApproverIntf.(neo4j.Node)
				fromApproverId := fromApproverNode.Props()["id"].(string)
				fromApproverKey := usermodel.NewUserKey(fromApproverId)
				fromUserApproversMap[fromApproverKey] = append(fromUserApproversMap[fromApproverKey], offerItemKey)
				fromItemApproversMap[offerItemKey] = append(fromItemApproversMap[offerItemKey], fromApproverKey)
			}

			for _, toApproverIntf := range toApproversSlice {
				toApproverNode := toApproverIntf.(neo4j.Node)
				toApproverId := toApproverNode.Props()["id"].(string)
				toApproverKey := usermodel.NewUserKey(toApproverId)
				toUserApproversMap[toApproverKey] = append(toUserApproversMap[toApproverKey], offerItemKey)
				toItemApproversMap[offerItemKey] = append(toItemApproversMap[offerItemKey], toApproverKey)
			}

		}

		userFromApprovers := map[usermodel.UserKey][]tradingmodel.OfferItemKey{}
		for userKey, offerItemKeys := range fromUserApproversMap {
			if _, ok := userFromApprovers[userKey]; !ok {
				userFromApprovers[userKey] = []tradingmodel.OfferItemKey{}
			}
			for _, offerItemKey := range offerItemKeys {
				userFromApprovers[userKey] = append(userFromApprovers[userKey], offerItemKey)
			}
		}
		userToApprovers := map[usermodel.UserKey][]tradingmodel.OfferItemKey{}
		for userKey, offerItemKeys := range toUserApproversMap {
			if _, ok := userToApprovers[userKey]; !ok {
				userToApprovers[userKey] = []tradingmodel.OfferItemKey{}
			}
			for _, offerItemKey := range offerItemKeys {
				userToApprovers[userKey] = append(userToApprovers[userKey], offerItemKey)
			}
		}
		itemFromApprovers := map[tradingmodel.OfferItemKey][]usermodel.UserKey{}
		for offerItemKey, userKeys := range fromItemApproversMap {
			if _, ok := itemFromApprovers[offerItemKey]; !ok {
				itemFromApprovers[offerItemKey] = []usermodel.UserKey{}
			}
			for _, userKey := range userKeys {
				itemFromApprovers[offerItemKey] = append(itemFromApprovers[offerItemKey], userKey)
			}
		}
		itemToApprovers := map[tradingmodel.OfferItemKey][]usermodel.UserKey{}
		for offerItemKey, userKeys := range toItemApproversMap {
			if _, ok := itemToApprovers[offerItemKey]; !ok {
				itemToApprovers[offerItemKey] = []usermodel.UserKey{}
			}
			for _, userKey := range userKeys {
				itemToApprovers[offerItemKey] = append(itemToApprovers[offerItemKey], userKey)
			}
		}

		userFromApproversMap := map[usermodel.UserKey]*tradingmodel.OfferItemKeys{}
		for userKey, offerItemKeys := range userFromApprovers {
			userFromApproversMap[userKey] = tradingmodel.NewOfferItemKeys(offerItemKeys)
		}

		userToApproversMap := map[usermodel.UserKey]*tradingmodel.OfferItemKeys{}
		for userKey, offerItemKeys := range userToApprovers {
			userToApproversMap[userKey] = tradingmodel.NewOfferItemKeys(offerItemKeys)
		}

		itemFromApproversMap := map[tradingmodel.OfferItemKey]*usermodel.UserKeys{}
		for offerItemKey, userKeys := range itemFromApprovers {
			itemFromApproversMap[offerItemKey] = usermodel.NewUserKeys(userKeys)
		}

		itemToApproversMap := map[tradingmodel.OfferItemKey]*usermodel.UserKeys{}
		for offerItemKey, userKeys := range itemToApprovers {
			itemToApproversMap[offerItemKey] = usermodel.NewUserKeys(userKeys)
		}

		offerApprovers := &tradingmodel.OfferApprovers{
			OfferKey:                  offerKey,
			OfferItemsUsersCanGive:    userFromApproversMap,
			OfferItemsUsersCanReceive: userToApproversMap,
			UsersAbleToGiveItem:       itemFromApproversMap,
			UsersAbleToReceiveItem:    itemToApproversMap,
		}

		approversForOffer = append(approversForOffer, offerApprovers)

	}

	offersApprovers := tradingmodel.NewOffersApprovers(approversForOffer)

	return offersApprovers, nil

}

func (t TradingStore) FindApproversForOffer(offerKey tradingmodel.OfferKey) (*tradingmodel.OfferApprovers, error) {

	approvers, err := t.FindApproversForOffers(tradingmodel.NewOfferKeys([]tradingmodel.OfferKey{offerKey}))
	if err != nil {
		return nil, err
	}

	approversForOffer, err := approvers.GetApproversForOffer(offerKey)
	if err != nil {
		return nil, err
	}

	return approversForOffer, nil

}

func (t TradingStore) SaveOffer(offer *tradingmodel.Offer, offerItems *tradingmodel.OfferItems) error {

	var matchClauses []string
	var createClauses []string
	var params = map[string]interface{}{}

	var userRefMap = map[usermodel.UserKey]string{}
	var groupRefMap = map[groupmodel.GroupKey]string{}
	var resourceRefMap = map[resourcemodel.ResourceKey]string{}

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
		matchClauses = append(matchClauses, "("+groupStr+":Group {id:$"+groupIdParamName+"})")
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

	matchClauses = append(matchClauses, "(group:Group{id:$group_id})")
	params["group_id"] = offer.GroupKey.String()

	now := time.Now().UTC()

	params["offer_id"] = offer.GetKey().String()
	params[CompletedKey] = offerItems.AllUserActionsCompleted()

	var completedAt *time.Time = nil
	if offerItems.AllUserActionsCompleted() {
		completedAt = &now
	}

	var acceptedAt *time.Time = nil
	if offerItems.AllPartiesAccepted() {
		acceptedAt = &now
	}

	params[CreatedAtKey] = now
	params[CompletedAtKey] = completedAt
	params[AcceptedAtKey] = acceptedAt
	statusStr, err := statusToString(offer.Status)

	if err != nil {
		return err
	}
	params[StatusKey] = statusStr

	// Add the CREATE (offer:Offer {id:$offer_id}) clause
	createClauses = append(createClauses, `
		(offer:Offer 
			{
				id                 : $offer_id,
				`+StatusKey+`      : $status,
				`+CreatedAtKey+`   : $created_at,
				`+UpdatedAtKey+`   : $created_at,
				`+CompletedKey+`   : $completed,
				`+CompletedAtKey+` : $completed_at,
				`+DeclinedAtKey+`  : null,
				`+CanceledAtKey+`  : null,
				`+AcceptedAtKey+`  : $accepted_at
			}
		)-[:CreatedBy]->(createdBy),(offer)-[:In]->(group)`)

	// Loop through each item
	for i, offerItem := range offerItems.Items {

		// The unique identifier for that offer item in the query
		offerItemRef := "offerItem" + strconv.Itoa(i+1)

		// add the type-specific params to the param map
		if offerItem.IsCreditTransfer() {

			creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)
			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem 
				{
					id:                        $`+offerItemRef+`_id,
					`+TypeKey+`:               $`+offerItemRef+`_type,
					`+AmountKey+`:             $`+offerItemRef+`_amount,
					`+FromApprovedKey+`:       $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:         $`+offerItemRef+`_to_approved,
					`+CreditsTransferredKey+`: $`+offerItemRef+`_credits_transferred,
					`+CompletedKey+`:          $`+offerItemRef+`_completed,
					`+CreatedAtKey+`:          $`+offerItemRef+`_created_at,
					`+UpdatedAtKey+`:          $`+offerItemRef+`_updated_at
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = offerItem.GetKey().String()
			params[offerItemRef+"_type"] = offerItem.Type()
			params[offerItemRef+"_amount"] = int(creditTransfer.Amount.Seconds())
			params[offerItemRef+"_from_approved"] = creditTransfer.GiverAccepted
			params[offerItemRef+"_to_approved"] = creditTransfer.ReceiverAccepted
			params[offerItemRef+"_credits_transferred"] = creditTransfer.CreditsTransferred
			params[offerItemRef+"_completed"] = creditTransfer.IsCompleted()
			params[offerItemRef+"_created_at"] = now
			params[offerItemRef+"_updated_at"] = now

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

			provideService := offerItem.(*tradingmodel.ProvideServiceItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id,
					`+TypeKey+`:         $`+offerItemRef+`_type,
					`+DurationKey+`:     $`+offerItemRef+`_duration,
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved,
					`+GivenKey+`:        $`+offerItemRef+`_given,
					`+ReceivedKey+`:     $`+offerItemRef+`_received,
					`+CreatedAtKey+`:    $`+offerItemRef+`_created_at,
					`+UpdatedAtKey+`:    $`+offerItemRef+`_updated_at
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = provideService.Key.String()
			params[offerItemRef+"_type"] = string(tradingmodel.ProvideService)
			params[offerItemRef+"_duration"] = int(provideService.Duration.Seconds())
			params[offerItemRef+"_from_approved"] = provideService.GiverAccepted
			params[offerItemRef+"_to_approved"] = provideService.ReceiverAccepted
			params[offerItemRef+"_given"] = provideService.ServiceGivenConfirmation
			params[offerItemRef+"_received"] = provideService.ServiceReceivedConfirmation
			params[offerItemRef+"_created_at"] = now
			params[offerItemRef+"_updated_at"] = now

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

			borrowResource := offerItem.(*tradingmodel.BorrowResourceItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id,
					`+TypeKey+`:         $`+offerItemRef+`_type,
					`+DurationKey+`:     $`+offerItemRef+`_duration,
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved,
					`+GivenKey+`:        $`+offerItemRef+`_given,
					`+TakenKey+`:        $`+offerItemRef+`_taken,
					`+ReturnedBackKey+`: $`+offerItemRef+`_returned_back,
					`+ReceivedBacKey+`:  $`+offerItemRef+`_received_back,
					`+CompletedKey+`:    $`+offerItemRef+`_completed,
					`+CreatedAtKey+`:    $`+offerItemRef+`_created_at,
					`+UpdatedAtKey+`:    $`+offerItemRef+`_updated_at
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = borrowResource.Key.String()
			params[offerItemRef+"_type"] = string(tradingmodel.BorrowResource)
			params[offerItemRef+"_duration"] = int(borrowResource.Duration.Seconds())
			params[offerItemRef+"_from_approved"] = borrowResource.GiverAccepted
			params[offerItemRef+"_to_approved"] = borrowResource.ReceiverAccepted
			params[offerItemRef+"_given"] = borrowResource.ItemGiven
			params[offerItemRef+"_taken"] = borrowResource.ItemTaken
			params[offerItemRef+"_returned_back"] = borrowResource.ItemReturnedBack
			params[offerItemRef+"_received_back"] = borrowResource.ItemReceivedBack
			params[offerItemRef+"_completed"] = borrowResource.IsCompleted()
			params[offerItemRef+"_created_at"] = now
			params[offerItemRef+"_updated_at"] = now

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

			resourceTransfer := offerItem.(*tradingmodel.ResourceTransferItem)

			createClauses = append(createClauses, "("+offerItemRef+`:OfferItem {
					id:                  $`+offerItemRef+`_id,
					`+TypeKey+`:         $`+offerItemRef+`_type,
					`+FromApprovedKey+`: $`+offerItemRef+`_from_approved,
					`+ToApprovedKey+`:   $`+offerItemRef+`_to_approved,
					`+GivenKey+`:        $`+offerItemRef+`_given,
					`+ReceivedKey+`:     $`+offerItemRef+`_received,
					`+CompletedKey+`:    $`+offerItemRef+`_completed,
					`+CreatedAtKey+`:    $`+offerItemRef+`_created_at,
					`+UpdatedAtKey+`:    $`+offerItemRef+`_updated_at
				})-[:IsPartOf]->(offer)`)

			params[offerItemRef+"_id"] = resourceTransfer.Key.String()
			params[offerItemRef+"_type"] = string(tradingmodel.ResourceTransfer)
			params[offerItemRef+"_from_approved"] = resourceTransfer.GiverAccepted
			params[offerItemRef+"_to_approved"] = resourceTransfer.ReceiverAccepted
			params[offerItemRef+"_given"] = resourceTransfer.ItemGiven
			params[offerItemRef+"_received"] = resourceTransfer.ItemReceived
			params[offerItemRef+"_completed"] = resourceTransfer.IsCompleted()
			params[offerItemRef+"_created_at"] = now
			params[offerItemRef+"_updated_at"] = now

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

func statusToString(offerStatus tradingmodel.OfferStatus) (string, error) {
	if offerStatus == tradingmodel.CompletedOffer {
		return "completed", nil
	} else if offerStatus == tradingmodel.DeclinedOffer {
		return "declined", nil
	} else if offerStatus == tradingmodel.AcceptedOffer {
		return "accepted", nil
	} else if offerStatus == tradingmodel.CanceledOffer {
		return "canceled", nil
	} else if offerStatus == tradingmodel.ExpiredOffer {
		return "expired", nil
	} else if offerStatus == tradingmodel.PendingOffer {
		return "pending", nil
	} else {
		return "", fmt.Errorf("unknown offer status type")
	}
}

func stringToStatus(offerStatus string) (tradingmodel.OfferStatus, error) {
	if offerStatus == "completed" {
		return tradingmodel.CompletedOffer, nil
	} else if offerStatus == "declined" {
		return tradingmodel.DeclinedOffer, nil
	} else if offerStatus == "accepted" {
		return tradingmodel.AcceptedOffer, nil
	} else if offerStatus == "canceled" {
		return tradingmodel.CanceledOffer, nil
	} else if offerStatus == "expired" {
		return tradingmodel.ExpiredOffer, nil
	} else if offerStatus == "pending" {
		return tradingmodel.PendingOffer, nil
	} else {
		return tradingmodel.CanceledOffer, fmt.Errorf("unknown offer status type")
	}
}

func (t TradingStore) UpdateOfferStatus(key tradingmodel.OfferKey, status tradingmodel.OfferStatus) error {
	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()
	statusString, err := statusToString(status)
	if err != nil {
		return err
	}

	now := time.Now().UTC()

	params := map[string]interface{}{
		"id":         key.String(),
		StatusKey:    statusString,
		UpdatedAtKey: now,
	}

	cypherUpdates := []string{
		StatusKey + ": $status",
		UpdatedAtKey + ": $updated_at",
	}

	if status == tradingmodel.AcceptedOffer {
		cypherUpdates = append(cypherUpdates, AcceptedAtKey+": $accepted_at")
		params[AcceptedAtKey] = now
	} else if status == tradingmodel.DeclinedOffer {
		cypherUpdates = append(cypherUpdates, DeclinedAtKey+": $declined_at")
		params[DeclinedAtKey] = now
	} else if status == tradingmodel.CanceledOffer {
		cypherUpdates = append(cypherUpdates, CanceledAtKey+": $canceled_at")
		params[CanceledAtKey] = now
	} else if status == tradingmodel.CompletedOffer {
		cypherUpdates = append(cypherUpdates, CompletedAtKey+": $completed_at")
		params[CompletedAtKey] = now
	}

	result, err := session.Run(`MATCH (o:Offer {id:$id}) SET o += {`+strings.Join(cypherUpdates, ",")+`} return o`, params)
	if err != nil {
		return err
	}
	if result.Err() != nil {
		return result.Err()
	}
	if !result.Next() {
		return exceptions.ErrOfferNotFound
	}
	return nil
}

func (t TradingStore) GetOfferItem(ctx context.Context, key tradingmodel.OfferItemKey) (tradingmodel.OfferItem, error) {

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
		return nil, exceptions.ErrOfferItemNotFound
	}

	offerIdField, _ := result.Record().Get("offerId")
	offerId := offerIdField.(string)
	offerKey, err := tradingmodel.ParseOfferKey(offerId)
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

func MapOfferItem(offerKey tradingmodel.OfferKey, offerItemNode neo4j.Node, fromNode neo4j.Node, toNode neo4j.Node) (tradingmodel.OfferItem, error) {

	offerItemType := offerItemNode.Props()["type"].(string)
	var fromResource *resourcemodel.Resource
	var fromTarget *model.Target
	var toTarget *model.Target
	var err error

	if fromNode != nil {
		if groupstore.IsGroupNode(fromNode) || store2.IsUserNode(fromNode) {
			fromTarget, err = sharedstore.MapOfferItemTarget(fromNode)
			if err != nil {
				return nil, err
			}
		} else if resourcestore.IsResourceNode(fromNode) {
			fromResource, err = resourcestore.MapResourceNode(fromNode)
			if err != nil {
				return nil, err
			}
		}
	}
	if toNode != nil {
		if groupstore.IsGroupNode(toNode) || store2.IsUserNode(toNode) {
			toTarget, err = sharedstore.MapOfferItemTarget(toNode)
			if err != nil {
				return nil, err
			}
		}

	}

	offerItemId := offerItemNode.Props()["id"].(string)
	offerItemKey, err := tradingmodel.ParseOfferItemKey(offerItemId)
	if err != nil {
		return nil, err
	}

	offerItemBase := tradingmodel.OfferItemBase{
		Type:             tradingmodel.OfferItemType(offerItemType),
		Key:              offerItemKey,
		OfferKey:         offerKey,
		To:               toTarget,
		ReceiverAccepted: offerItemNode.Props()[ToApprovedKey].(bool),
		GiverAccepted:    offerItemNode.Props()[FromApprovedKey].(bool),
		CreatedAt:        offerItemNode.Props()[CreatedAtKey].(time.Time),
		UpdatedAt:        offerItemNode.Props()[UpdatedAtKey].(time.Time),
	}

	if offerItemType == string(tradingmodel.CreditTransfer) {

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

		return &tradingmodel.CreditTransferItem{
			OfferItemBase:      offerItemBase,
			From:               fromTarget,
			Amount:             time.Duration(int64(time.Second) * amount.(int64)),
			CreditsTransferred: offerItemNode.Props()[CreditsTransferredKey].(bool),
		}, nil

	} else if offerItemType == string(tradingmodel.ProvideService) {

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

		return &tradingmodel.ProvideServiceItem{
			OfferItemBase:               offerItemBase,
			ResourceKey:                 fromResource.Key,
			Duration:                    time.Duration(int64(time.Second) * duration.(int64)),
			ServiceReceivedConfirmation: offerItemNode.Props()[ReceivedKey].(bool),
			ServiceGivenConfirmation:    offerItemNode.Props()[GivenKey].(bool),
		}, nil

	} else if offerItemType == string(tradingmodel.BorrowResource) {

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

		return &tradingmodel.BorrowResourceItem{
			OfferItemBase:    offerItemBase,
			ResourceKey:      fromResource.Key,
			Duration:         time.Duration(int64(time.Second) * duration.(int64)),
			ItemTaken:        offerItemNode.Props()[TakenKey].(bool),
			ItemGiven:        offerItemNode.Props()[GivenKey].(bool),
			ItemReturnedBack: offerItemNode.Props()[ReturnedBackKey].(bool),
			ItemReceivedBack: offerItemNode.Props()[ReceivedBacKey].(bool),
		}, nil

	} else if offerItemType == string(tradingmodel.ResourceTransfer) {

		if fromResource == nil {
			return nil, fmt.Errorf("result should have a 'From' resource")
		}

		if toTarget == nil {
			return nil, fmt.Errorf("result should have a 'To' user/group")
		}

		return &tradingmodel.ResourceTransferItem{
			OfferItemBase: offerItemBase,
			ResourceKey:   fromResource.Key,
			ItemReceived:  offerItemNode.Props()[ReceivedKey].(bool),
			ItemGiven:     offerItemNode.Props()[GivenKey].(bool),
		}, nil

	} else {
		return nil, fmt.Errorf("unexpected offer item type: %s", offerItemType)
	}

}

func (t TradingStore) UpdateOfferItem(ctx context.Context, offerItem tradingmodel.OfferItem) error {
	session, err := t.graphDriver.GetSession()
	if err != nil {
		return err
	}
	defer session.Close()

	var result neo4j.Result

	if offerItem.IsServiceProviding() {

		service := offerItem.(*tradingmodel.ProvideServiceItem)
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
				"id":            service.Key.String(),
				FromApprovedKey: service.GiverAccepted,
				ToApprovedKey:   service.ReceiverAccepted,
				GivenKey:        service.ServiceGivenConfirmation,
				ReceivedKey:     service.ServiceReceivedConfirmation,
				CompletedKey:    service.IsCompleted(),
			})

	} else if offerItem.IsCreditTransfer() {

		creditTransfer := offerItem.(*tradingmodel.CreditTransferItem)
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
				"id":                  creditTransfer.Key.String(),
				FromApprovedKey:       creditTransfer.GiverAccepted,
				ToApprovedKey:         creditTransfer.ReceiverAccepted,
				CreditsTransferredKey: creditTransfer.CreditsTransferred,
				CompletedKey:          creditTransfer.IsCompleted(),
			})

	} else if offerItem.IsBorrowingResource() {
		resourceBorrow := offerItem.(*tradingmodel.BorrowResourceItem)
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
				"id":            resourceBorrow.Key.String(),
				FromApprovedKey: resourceBorrow.GiverAccepted,
				ToApprovedKey:   resourceBorrow.ReceiverAccepted,
				GivenKey:        resourceBorrow.ItemGiven,
				TakenKey:        resourceBorrow.ItemTaken,
				ReturnedBackKey: resourceBorrow.ItemReturnedBack,
				ReceivedBacKey:  resourceBorrow.ItemReceivedBack,
				CompletedKey:    resourceBorrow.IsCompleted(),
			})

	} else if offerItem.IsResourceTransfer() {

		service := offerItem.(*tradingmodel.ResourceTransferItem)
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
				"id":            service.Key.String(),
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
		return exceptions.ErrOfferItemNotFound
	}
	return nil
}

func (t TradingStore) ConfirmItemGiven(ctx context.Context, key tradingmodel.OfferItemKey) error {
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
		return exceptions.ErrOfferItemNotFound
	}
	return nil
}

func (t TradingStore) GetOffer(key tradingmodel.OfferKey) (*tradingmodel.Offer, error) {

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
		return nil, exceptions.ErrOfferNotFound
	}

	createdByField, _ := result.Record().Get("createdById")
	createdById := createdByField.(string)
	createdByKey := usermodel.NewUserKey(createdById)

	offerField, _ := result.Record().Get("offer")
	offerNode := offerField.(neo4j.Node)

	return MapOfferNode(offerNode, createdByKey)

}

func MapOfferNode(node neo4j.Node, createdByKey usermodel.UserKey) (*tradingmodel.Offer, error) {

	offerId := node.Props()["id"].(string)
	offerKey, err := tradingmodel.ParseOfferKey(offerId)
	if err != nil {
		return nil, err
	}

	status, err := stringToStatus(node.Props()["status"].(string))
	if err != nil {
		return nil, err
	}

	return &tradingmodel.Offer{
		Key:          offerKey,
		CreatedByKey: createdByKey,
		Status:       status,
		CreatedAt:    node.Props()["created_at"].(time.Time),
	}, nil

}

func (t TradingStore) GetOfferItemsForOffer(key tradingmodel.OfferKey) (*tradingmodel.OfferItems, error) {

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

	var offerItems []tradingmodel.OfferItem
	for result.Next() {

		offerIdField, _ := result.Record().Get("offerId")
		offerId := offerIdField.(string)
		offerKey, err := tradingmodel.ParseOfferKey(offerId)
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

	return tradingmodel.NewOfferItems(offerItems), nil

}

func (t TradingStore) GetOffersForUser(userKey usermodel.UserKey) (*trading.GetOffersResult, error) {

	session, err := t.graphDriver.GetSession()
	if err != nil {
		return nil, err
	}
	defer session.Close()

	result, err := session.Run(`
		match 
		(user:User {id:$userId})-[membership:IsMemberOf]->(group:Group),
		(group)<-[inRel:In]-(offer:Offer)<-[partOfRel:IsPartOf]-(offerItem:OfferItem),
		(from)-[fromRel:From]->(offerItem)
		optional match (to)<-[toRel:To]-(offerItem)
		
		call {
			with from, fromRel, offer, offerItem, user
			match (from)-[fromRel]->(offerItem)-[partOfRel]->(offer)
			where from = user
			return offer as o
			
			union
			
			with from, fromRel, offer, offerItem, user, group, membership
			match (from)-[fromRel]->(offerItem)-[partOfRel]->(offer),(user)-[membership]->(group)
			where from = group and membership.isAdmin
			return offer as o
			
			union
			
			with to, toRel, offerItem, offer, user
			match (to)<-[toRel]-(offerItem)
			where to = user
			return offer as o
			
			union
			
			with to, toRel, offerItem, partOfRel, offer, inRel, group, membership, user
			match (to)<-[toRel]-(offerItem)-[partOfRel]->(offer)-[inRel]->(group)<-[membership]-(user)
			where to = group and  membership.isAdmin
			return offer as o
		}
		
		with collect(distinct o) as offers
		UNWIND offers as offer
		match (from)-[:From]->(offerItem)-[:IsPartOf]->(offer)-[:CreatedBy]->(creator),(to)<-[:To]-(offerItem)
		
		return offer, creator.id as createdById, collect({offerItem: offerItem, from: from, to: to}) as offerItems`,
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
		createdByKey := usermodel.NewUserKey(createdById)

		offer, err := MapOfferNode(offerNode, createdByKey)
		if err != nil {
			return nil, err
		}

		offerItemsContainerField, _ := record.Get("offerItems")
		offerItemsContainerSlice := offerItemsContainerField.([]interface{})

		var offerItems []tradingmodel.OfferItem
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
			OfferItems: tradingmodel.NewOfferItems(offerItems),
		})

	}

	return &trading.GetOffersResult{
		Items: resultItems,
	}, nil

}

func (t TradingStore) GetTradingHistory(ctx context.Context, ids *usermodel.UserKeys) ([]tradingmodel.HistoryEntry, error) {
	panic("implement me")
}
