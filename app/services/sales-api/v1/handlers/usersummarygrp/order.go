package usersummarygrp

import (
	"errors"
	"net/http"

	"github.com/owezzy/service-5/business/core/usersummary"
	"github.com/owezzy/service-5/business/data/order"
	"github.com/owezzy/service-5/foundation/validate"
)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByUserID   = "user_id"
		orderByUserName = "user_name"
	)

	var orderByFields = map[string]string{
		orderByUserID:   usersummary.OrderByUserID,
		orderByUserName: usersummary.OrderByUserName,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByUserID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
