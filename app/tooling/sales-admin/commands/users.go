package commands

import (
	"context"
	"encoding/json"
	"fmt"
	"github.com/owezzy/service-5/business/core/event"
	"github.com/owezzy/service-5/business/core/user"
	"github.com/owezzy/service-5/business/core/user/stores/userdb"
	db "github.com/owezzy/service-5/business/data/dbsql/pgx"
	"github.com/owezzy/service-5/foundation/logger"
	"os"
	"strconv"
	"time"
)

// Users retrieves all users from the database.
func Users(log *logger.Logger, cfg db.Config, pageNumber string, rowsPerPage string) error {
	db, err := db.Open(cfg)
	if err != nil {
		return fmt.Errorf("connect database: %w", err)
	}
	defer db.Close()

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	page, err := strconv.Atoi(pageNumber)
	if err != nil {
		return fmt.Errorf("converting page number: %w", err)
	}

	rows, err := strconv.Atoi(rowsPerPage)
	if err != nil {
		return fmt.Errorf("converting rows per page: %w", err)
	}

	evnCore := event.NewCore(log)
	core := user.NewCore(log, evnCore, userdb.NewStore(log, db))

	users, err := core.Query(ctx, user.QueryFilter{}, user.DefaultOrderBy, page, rows)
	if err != nil {
		return fmt.Errorf("retrieve users: %w", err)
	}

	return json.NewEncoder(os.Stdout).Encode(users)
}
