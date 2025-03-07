package product

import "github.com/owezzy/service-5/business/data/order"

// DefaultOrderBy represents the default way we sort.
var DefaultOrderBy = order.NewBy(OrderByProdID, order.ASC)

// Set of fields that the results can be ordered by. These are the names
// that should be used by the application layer.
const (
	OrderByProdID   = "product_id"
	OrderByName     = "name"
	OrderByCost     = "cost"
	OrderByQuantity = "quantity"
	OrderBySold     = "sold"
	OrderByRevenue  = "revenue"
	OrderByUserID   = "user_id"
)
