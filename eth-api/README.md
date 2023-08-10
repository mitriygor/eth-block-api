


my-fiber-project/
|- app/
|    |- handlers/
|    |      |- user_handler.go
|    |      |- auth_handler.go
|    |
|    |- middleware/
|    |       |- auth_middleware.go
|    |       |- logging_middleware.go
|    |
|    |- models/
|    |     |- user.go
|    |     |- post.go
|    |
|    |- routes/
|    |      |- routes.go
|    |
|    |- config/
|    |      |- config.go
|    |
|    |- server.go
|
|- migrations/
|      |- 20220101000001_create_users_table.up.sql
|      |- 20220101000001_create_users_table.down.sql
|      |- 20220102000001_create_posts_table.up.sql
|      |- 20220102000001_create_posts_table.down.sql
|
|- .env
|- go.mod
|- go.sum
|- main.go




app/: Contains the main application logic.

handlers/: HTTP request handlers.
middleware/: Custom middleware functions.
models/: Database models or structures.
routes/: Defines the routes and their handlers.
config/: Configuration-related files.
migrations/: Contains SQL migration files for database schema changes using a tool like "migrate". The migration files should be named with a timestamp prefix.

.env: Environment file for storing sensitive configuration variables.

go.mod and go.sum: Files that manage project dependencies using Go modules.

main.go: The entry point of the application.




