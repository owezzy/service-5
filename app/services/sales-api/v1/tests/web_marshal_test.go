package tests

import (
	"time"

	"github.com/google/uuid"
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/homegrp"
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/productgrp"
	"github.com/owezzy/service-5/app/services/sales-api/v1/handlers/usergrp"
	"github.com/owezzy/service-5/business/core/home"
	"github.com/owezzy/service-5/business/core/product"
	"github.com/owezzy/service-5/business/core/user"
)

func toAppUser(usr user.User) usergrp.AppUser {
	roles := make([]string, len(usr.Roles))
	for i, role := range usr.Roles {
		roles[i] = role.Name()
	}

	return usergrp.AppUser{
		ID:           usr.ID.String(),
		Name:         usr.Name,
		Email:        usr.Email.Address,
		Roles:        roles,
		PasswordHash: nil, // This field is not marshalled.
		Department:   usr.Department,
		Enabled:      usr.Enabled,
		DateCreated:  usr.DateCreated.Format(time.RFC3339),
		DateUpdated:  usr.DateUpdated.Format(time.RFC3339),
	}
}

func toAppUsers(users []user.User) []usergrp.AppUser {
	items := make([]usergrp.AppUser, len(users))
	for i, usr := range users {
		items[i] = toAppUser(usr)
	}

	return items
}

// =============================================================================

func toAppProductDetails(prd product.Product, usr user.User) productgrp.AppProductDetails {
	return productgrp.AppProductDetails{
		ID:          prd.ID.String(),
		Name:        prd.Name,
		Cost:        prd.Cost,
		Quantity:    prd.Quantity,
		UserID:      prd.UserID.String(),
		UserName:    usr.Name,
		DateCreated: prd.DateCreated.Format(time.RFC3339),
		DateUpdated: prd.DateUpdated.Format(time.RFC3339),
	}
}

func toAppProductsDetails(prds []product.Product, usrs map[uuid.UUID]user.User) []productgrp.AppProductDetails {
	items := make([]productgrp.AppProductDetails, len(prds))
	for i, prd := range prds {
		items[i] = toAppProductDetails(prd, usrs[prd.UserID])
	}

	return items
}

// =============================================================================

func toAppHome(hme home.Home) homegrp.AppHome {
	return homegrp.AppHome{
		ID:     hme.ID.String(),
		UserID: hme.UserID.String(),
		Type:   hme.Type,
		Address: homegrp.AppNewAddress{
			Address1: hme.Address.Address1,
			Address2: hme.Address.Address2,
			ZipCode:  hme.Address.ZipCode,
			City:     hme.Address.City,
			State:    hme.Address.State,
			Country:  hme.Address.Country,
		},
		DateCreated: hme.DateCreated.Format(time.RFC3339),
		DateUpdated: hme.DateUpdated.Format(time.RFC3339),
	}
}

func toAppHomes(homes []home.Home) []homegrp.AppHome {
	items := make([]homegrp.AppHome, len(homes))
	for i, hme := range homes {
		items[i] = toAppHome(hme)
	}

	return items
}
