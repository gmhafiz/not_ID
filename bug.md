Ent selects `id` when it does not exist in the database because assumes each schema has an `id ` column.


- [x] The issue is present in the latest release.
- [x] I have searched the [issues](https://github.com/ent/ent/issues) of this repository and believe that this is not a duplicate.

## Current Behavior 😯

/ent/schema/sessions.go
```go
// Fields of the User.
func (Session) Fields() []ent.Field {
	return []ent.Field{
		field.String("token"),
		field.Uint64("user_id").Optional().Nillable(),
		field.Time("created_at"),
	}
}

// Edges of the User.
func (Session) Edges() []ent.Edge {
	return []ent.Edge{
		edge.From("user", User.Type).
			Field("user_id").
			Ref("session").
			Unique(),
	}
}
```

Schema is created externally, not by ent

```sql
CREATE TABLE IF NOT EXISTS sessions
(
    token   VARCHAR(255) PRIMARY KEY,
    user_id BIGINT UNSIGNED NULL,
    created_at  datetime NOT NULL,
    
    FOREIGN KEY (user_id) REFERENCES users(id)
);
```

When you want to select a record in `sesssion` table using its primary key which is `token` (not `id`):

bug_test.go
```go
first, err := client.Debug().Session.Query().
    Where(session.TokenEQ("random")).
    First(ctx)
if err != nil {
    t.Fatal(err)
}
```

`Debug()` shows that it selects non-existent `id`

```
=== RUN   TestBugMySQL/56
2023/05/04 12:51:09 driver.Query: query=SELECT DISTINCT `sessions`.`id`, `sessions`.`token`, `sessions`.`user_id`, `sessions`.`created_at` FROM `sessions` WHERE `sessions`.`token` = ? LIMIT 1 args=[random]
```


## Expected Behavior 🤔

Ent should not select id column like this

```
SELECT DISTINCT `sessions`.`token`, `sessions`.`user_id`, `sessions`.`created_at` FROM `sessions` WHERE `sessions`.`token` = ? LIMIT 1
```

## Steps to Reproduce 🕹

```sh
git clone https://github.com/gmhafiz/not_ID
cd not_ID
docker-compose up -d
go test -v bug_test.go
```

Returns

```sh
=== RUN   TestBugMaria
=== RUN   TestBugMaria/10.5
2023/05/04 13:07:14 driver.Query: query=SELECT `sessions`.`id`, `sessions`.`token`, `sessions`.`user_id`, `sessions`.`created_at` FROM `sessions` WHERE `sessions`.`token` = ? LIMIT 1 args=[random]
    bug_test.go:132: Error 1054 (42S22): Unknown column 'sessions.id' in 'field list'
--- FAIL: TestBugMaria (0.02s)
    --- FAIL: TestBugMaria/10.5 (0.02s)
FAIL
FAIL    command-line-arguments  0.024s
FAIL
```


## Your Environment 🌎

| Tech     | Version                                                                                  |
|----------|------------------------------------------------------------------------------------------|
| Go       | 1.20.3                                                                                   |
| Ent      | 0.12.3                                                                                   |
| Database | mariadb, but applies to all                                                              |
| Driver   | https://github.com/go-sql-driver/mysql , github.com/lib/pq , github.com/mattn/go-sqlite3 |
