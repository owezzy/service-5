package productgrp

import (
	"errors"
	"github.com/owezzy/service-5/business/core/product"
	"github.com/owezzy/service-5/business/data/order"
	"github.com/owezzy/service-5/foundation/validate"
	"net/http"


)

func parseOrder(r *http.Request) (order.By, error) {
	const (
		orderByProdID   = "product_id"
		orderByName     = "name"
		orderByCost     = "cost"
		orderByQuantity = "quantity"
		orderBySold     = "sold"
		orderByRevenue  = "revenue"
		orderByUserID   = "user_id"
	)

	var orderByFields = map[string]string{
		orderByProdID:   product.OrderByProdID,
		orderByName:     product.OrderByName,
		orderByCost:     product.OrderByCost,
		orderByQuantity: product.OrderByQuantity,
		orderBySold:     product.OrderBySold,
		orderByRevenue:  product.OrderByRevenue,
		orderByUserID:   product.OrderByUserID,
	}

	orderBy, err := order.Parse(r, order.NewBy(orderByProdID, order.ASC))
	if err != nil {
		return order.By{}, err
	}

	if _, exists := orderByFields[orderBy.Field]; !exists {
		return order.By{}, validate.NewFieldsError(orderBy.Field, errors.New("order field does not exist"))
	}

	orderBy.Field = orderByFields[orderBy.Field]

	return orderBy, nil
}
