package testNeonDatabase

import (
	"context"
	"fmt"

	"github.com/jackc/pgx/v5"
)

func TestNeonDatabase(ctx context.Context, neonConnection *pgx.Conn) {

	_, err := neonConnection.Exec(ctx, "CREATE TABLE IF NOT EXISTS playing_with_neon(id SERIAL PRIMARY KEY, name TEXT NOT NULL, value REAL);")
	if err != nil {
		panic(err)
	}
	_, err = neonConnection.Exec(ctx, "INSERT INTO playing_with_neon(name, value) SELECT LEFT(md5(i::TEXT), 10), random() FROM generate_series(1, 10) s(i);")
	if err != nil {
		panic(err)
	}
	rows, err := neonConnection.Query(ctx, "SELECT * FROM playing_with_neon")
	if err != nil {
		panic(err)
	}
	defer rows.Close()
	for rows.Next() {
		var id int32
		var name string
		var value float32
		if err := rows.Scan(&id, &name, &value); err != nil {
			panic(err)
		}
		fmt.Printf("%d | %s | %f\n", id, name, value)
	}
}
