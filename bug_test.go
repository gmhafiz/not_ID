package bug

import (
	"context"
	"entgo.io/bug/ent"
	//"entgo.io/bug/ent/enttest"
	"entgo.io/bug/ent/session"
	"entgo.io/ent/dialect"
	"entgo.io/ent/dialect/sql"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	_ "github.com/lib/pq"
	_ "github.com/mattn/go-sqlite3"
	"github.com/stretchr/testify/assert"
	"net"
	"strconv"
	"testing"
	"time"
)

//func TestBugSQLite(t *testing.T) {
//	client := enttest.Open(t, dialect.SQLite, "file:ent?mode=memory&cache=shared&_fk=1")
//	defer client.Close()
//	test(t, client)
//}
//
//func TestBugMySQL(t *testing.T) {
//	for version, port := range map[string]int{"56": 3306, "57": 3307, "8": 3308} {
//		addr := net.JoinHostPort("localhost", strconv.Itoa(port))
//		t.Run(version, func(t *testing.T) {
//			client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
//			defer client.Close()
//			test(t, client)
//		})
//	}
//}
//
//func TestBugPostgres(t *testing.T) {
//	for version, port := range map[string]int{"10": 5430, "11": 5431, "12": 5432, "13": 5433, "14": 5434} {
//		t.Run(version, func(t *testing.T) {
//			client := enttest.Open(t, dialect.Postgres, fmt.Sprintf("host=localhost port=%d user=postgres dbname=test password=pass sslmode=disable", port))
//			defer client.Close()
//			test(t, client)
//		})
//	}
//}

func TestBugMaria(t *testing.T) {
	for version, port := range map[string]int{"10.5": 4306} {
		t.Run(version, func(t *testing.T) {
			addr := net.JoinHostPort("localhost", strconv.Itoa(port))

			db, err := sql.Open(dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			assert.Nil(t, err)

			ctx := context.Background()

			// Schema is already created beforehand, not by ent. So we simulate this
			// before calling ent.Open()
			_, err = db.ExecContext(ctx, `
				CREATE TABLE users
				(
					id      bigint unsigned auto_increment primary key ,
					name    varchar(255),
    				created_at  datetime NOT NULL
				);
			`)
			assert.Nil(t, err)

			_, err = db.ExecContext(ctx, `
				CREATE TABLE IF NOT EXISTS sessions
				(
					token   VARCHAR(255) PRIMARY KEY,
					user_id BIGINT UNSIGNED NULL,
					created_at  datetime NOT NULL,
					
					FOREIGN KEY (user_id) REFERENCES users(id)
				);
			`)
			assert.Nil(t, err)

			// Do not use built-in migrate.
			//client := enttest.Open(t, dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))

			// Instead, open but without using ent's migration.
			client, err := ent.Open(dialect.MySQL, fmt.Sprintf("root:pass@tcp(%s)/test?parseTime=True", addr))
			assert.Nil(t, err)

			defer client.Close()
			test(t, client)
		})
	}
}

func test(t *testing.T, client *ent.Client) {
	ctx := context.Background()

	client.User.Delete().ExecX(ctx)
	client.User.Create().
		SetCreatedAt(time.Now()).
		SetName("Ariel").
		ExecX(ctx)
	if n := client.User.Query().CountX(ctx); n != 1 {
		t.Errorf("unexpected number of users: %d", n)
	}

	id, err := client.User.Query().FirstID(ctx)
	if err != nil {
		t.Fatal(err)
	}

	client.Session.Delete().ExecX(ctx)
	client.Session.Create().
		SetToken("random").
		SetCreatedAt(time.Now()).
		ExecX(ctx)

	first, err := client.Debug().Session.Query().Where(session.TokenEQ("random")).First(ctx)

	/*
		generates wrong query

		2023/05/04 12:51:09 driver.Query: query=SELECT DISTINCT `sessions`.`id`, `sessions`.`token`, `sessions`.`user_id`, `sessions`.`created_at` FROM `sessions` WHERE `sessions`.`token` = ? LIMIT 1 args=[random]

		But table `sessions` has no `id` column.

		Ent assumes all tables uses `id` as primary column. But in this case,
		table `sessions` is using `token` with type varchar as its primary key.

	*/
	if err != nil {
		t.Fatal(err)
	}

	/*
		Expected query should be:

		SELECT DISTINCT `sessions`.`token`, `sessions`.`user_id`, `sessions`.`created_at` FROM `sessions` WHERE `sessions`.`token` = ? LIMIT 1 args=[random]
	*/

	first.Update().SetUserID(id).ExecX(ctx)
}
