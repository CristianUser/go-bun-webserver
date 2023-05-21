package migrations

import (
	"context"
	"fmt"

	"github.com/cristianuser/go-bun-webserver/bunapp"
	"github.com/cristianuser/go-bun-webserver/org"
	"github.com/uptrace/bun"
	"github.com/uptrace/bun/dbfixture"
)

func init() {
	Migrations.MustRegister(func(ctx context.Context, db *bun.DB) error {
		db.RegisterModel((*org.User)(nil)) // required by dbfixture
		db.NewCreateTable().Model((*org.Session)(nil)).IfNotExists().Exec(ctx)
		fixture := dbfixture.New(db, dbfixture.WithRecreateTables())
		return fixture.Load(ctx, bunapp.FS(), "fixture/fixture.yml")
	}, func(ctx context.Context, db *bun.DB) error {
		fmt.Print(" [down migration] ")
		return nil
	})
}
