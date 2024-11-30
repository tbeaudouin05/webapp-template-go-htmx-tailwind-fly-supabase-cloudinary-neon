package getNeonConnection

import (
	"context"
	"webapp-template-go-htmx-tailwind-fly-supabase-cloudinary-neon/goEnv"

	"github.com/jackc/pgx/v5"
)

func GetNeonConnection(ctx context.Context) *pgx.Conn {
	connStr := goEnv.GlobalEnvVar.NeonDatabaseUrl
	neonConnection, err := pgx.Connect(ctx, connStr)
	if err != nil {
		println("Error connecting to Neon database: %s", err)
		panic(err)
	}

	return neonConnection
}
