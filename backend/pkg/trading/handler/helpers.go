package handler

import (
	fmt "fmt"
	"github.com/commonpool/backend/pkg/exceptions"
	"github.com/commonpool/backend/pkg/keys"
	"github.com/commonpool/backend/pkg/trading/domain"
	"github.com/labstack/echo/v4"
	"time"
)

func parseTargetFromQueryParams(c echo.Context, typeQueryParam string, valueQueryParam string) (*keys.Target, error) {
	typeParam := c.QueryParams().Get(typeQueryParam)
	if typeParam != "" {
		typeValue, err := keys.ParseOfferItemTargetType(typeParam)
		if err != nil {
			return nil, err
		}
		targetType := &typeValue
		targetIdStr := c.QueryParams().Get(valueQueryParam)
		if targetIdStr == "" {
			return nil, exceptions.ErrQueryParamRequired(valueQueryParam)
		}
		if targetType.IsGroup() {
			groupKey, err := keys.ParseGroupKey(targetIdStr)
			if err != nil {
				return nil, err
			}
			return keys.NewGroupTarget(groupKey), nil
		} else if targetType.IsUser() {
			userKey := keys.NewUserKey(targetIdStr)
			return keys.NewUserTarget(userKey), nil
		}
	}
	return nil, nil
}

func mapNewOfferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (domain.OfferItem, error) {

	itemType := tradingOfferItem.Type

	if itemType == domain.CreditTransfer {

		res, err := mapCreateCreditTransferItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.CreditTransferItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.ResourceTransfer {

		res, err := mapCreateResourceTransferItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.ResourceTransferItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.ProvideService {

		res, err := mapCreateProvideServiceItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.ProvideServiceItem: %v", err)
		}
		return res, nil

	} else if itemType == domain.BorrowResource {

		res, err := mapCreateBorrowItem(tradingOfferItem, itemKey)
		if err != nil {
			return nil, fmt.Errorf("error mapping tradingOfferItem as trading.BorrowItem: %v", err)
		}
		return res, nil

	} else {

		return nil, fmt.Errorf("unexpected item type: %s", itemType)

	}
}

func mapCreateBorrowItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.BorrowResourceItem, error) {
	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Duration' as time.Duration: %v", err)
	}

	return &domain.BorrowResourceItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.BorrowResource,
			Key:  itemKey,
			To:   &tradingOfferItem.To,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateProvideServiceItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.ProvideServiceItem, error) {
	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	duration, err := time.ParseDuration(*tradingOfferItem.Duration)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Duration' as time.Duration: %v", err)
	}

	return &domain.ProvideServiceItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.ProvideService,
			Key:  itemKey,
			To:   &tradingOfferItem.To,
		},
		ResourceKey: resourceKey,
		Duration:    duration,
	}, nil
}

func mapCreateResourceTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.ResourceTransferItem, error) {
	resourceKey, err := keys.ParseResourceKey(*tradingOfferItem.ResourceId)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.ResourceId' as keys.ResourceKey: %v", err)
	}

	return &domain.ResourceTransferItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.ResourceTransfer,
			Key:  itemKey,
			To:   &tradingOfferItem.To,
		},
		ResourceKey: resourceKey,
	}, nil
}

func mapCreateCreditTransferItem(tradingOfferItem SendOfferPayloadItem, itemKey keys.OfferItemKey) (*domain.CreditTransferItem, error) {
	amount, err := time.ParseDuration(*tradingOfferItem.Amount)
	if err != nil {
		return nil, fmt.Errorf("error parsing 'tradingOfferItem.Amount' as time.Duration: %v", err)
	}

	return &domain.CreditTransferItem{
		OfferItemBase: domain.OfferItemBase{
			Type: domain.CreditTransfer,
			Key:  itemKey,
			To:   &tradingOfferItem.To,
		},
		From:   tradingOfferItem.From,
		Amount: amount,
	}, nil
}
