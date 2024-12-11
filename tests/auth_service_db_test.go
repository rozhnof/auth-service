package tests

import (
	"context"

	"github.com/rozhnof/auth-service/internal/infrastructure/database/postgres"
)

type AuthServiceDatabase struct {
	db postgres.Database
}

func NewAuthServiceDatabase(db postgres.Database) AuthServiceDatabase {
	return AuthServiceDatabase{
		db: db,
	}
}

func (d *AuthServiceDatabase) Truncate(ctx context.Context) error {
	query := `
		DO $$ DECLARE
			table_name TEXT;
		BEGIN
			FOR table_name IN 
				SELECT tablename 
				FROM pg_tables 
				WHERE schemaname = 'public'
			LOOP
				EXECUTE format('TRUNCATE TABLE %I CASCADE', table_name);
			END LOOP;
		END $$;
	`
	if _, err := d.db.Exec(ctx, query); err != nil {
		return err
	}

	return nil
}
