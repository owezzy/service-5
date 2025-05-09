package commands

import (
	"context"
	"fmt"
	"github.com/owezzy/service-5/business/core/event"
	"github.com/owezzy/service-5/business/core/user"
	"github.com/owezzy/service-5/business/core/user/stores/userdb"
	db "github.com/owezzy/service-5/business/data/dbsql/pgx"
	"github.com/owezzy/service-5/foundation/logger"
	"net/mail"
	"time"
)

// UserAdd adds new users into the database.
func UserAdd(log *logger.Logger, cfg db.Config, name, email, password string) error {
	if name == "" || email == "" || password == "" {
		fmt.Println("help: useradd <name> <email> <password>")
		return ErrHelp
	}

	db, err := db.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	evnCore := event.NewCore(log)
	core := user.NewCore(log, evnCore, userdb.NewStore(log, db))

	addr, err := mail.ParseAddress(email)
	if err != nil {
		return fmt.Errorf("parsing email: %w", err)
	}

	nu := user.NewUser{
		Name:            name,
		Email:           *addr,
		Password:        password,
		PasswordConfirm: password,
		Roles:           []user.Role{user.RoleAdmin, user.RoleUser},
	}

	usr, err := core.Create(ctx, nu)
	if err != nil {
		return fmt.Errorf("create user: %w", err)
	}

	fmt.Println("user id:", usr.ID)
	return nil
}
