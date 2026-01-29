# SQL Migration Guide

## Overview

When modifying database schemas in jzero, you MUST create corresponding migration files in `desc/sql_migration/` to track changes and enable rollbacks.

## ⚠️ Critical Rules

1. **Always create migration files** for any schema change (tables, columns, indexes, etc.)
2. **Development**: Use `jzero migrate` commands with `.jzero.yaml`
3. **Production**: Use code-based migration in `cmd/server.go` - automatic on startup
4. **Both modes**: Create files in `desc/sql_migration/` with up/down versions

## Migration Execution Modes

### Development Mode: jzero migrate Command

Use during development for quick testing:

```bash
jzero migrate up              # Apply pending migrations
jzero migrate down            # Rollback last migration
jzero migrate up 4            # Apply specific number of migrations
jzero migrate goto 3          # Go to specific version
jzero migrate version         # Check current version
```

**Requires**: `.jzero.yaml` configuration (see below)

### Production Mode: Code-Based Migration

Automatic migration on server startup - **recommended for production**:

```go
// cmd/server.go
var serverCmd = &cobra.Command{
    Use: "server",
    Run: func(cmd *cobra.Command, args []string) {
        cc := configcenter.MustNewConfigCenter[config.Config](...)

        // Run migrations before starting server
        m, err := migrate.NewMigrate(cc.MustGetConfig().Sqlx.SqlConf)
        logx.Must(err)
        defer m.Close()
        logx.Must(m.Up())

        // Continue with server startup...
        svcCtx := svc.NewServiceContext(cc)
        restServer := rest.MustNewServer(...)
        restServer.Start()
    },
}
```

**For multi-database support**:
```go
m, err := migrate.NewMigrate(
    cc.MustGetConfig().Sqlx.SqlConf,
    migrate.WithSourceAppendDriver(true),  // Use desc/sql_migration/{driver}/
)
```

**Benefits**: Automatic migration, uses `etc/etc.yaml` config, no manual steps needed

## Migration File Structure

### Default Mode (Single Database)

All migrations in `desc/sql_migration/`:

```
desc/sql_migration/
├── 1_create_users_table.up.sql
├── 1_create_users_table.down.sql
├── 2_add_email_index.up.sql
└── 2_add_email_index.down.sql
```

### Multi-Database Mode (sourceAppendDriver)

Separate directories per database type:

```
desc/sql_migration/
├── mysql/
│   ├── 1_create_users_table.up.sql
│   └── 1_create_users_table.down.sql
├── pgx/
│   ├── 1_create_users_table.up.sql
│   └── 1_create_users_table.down.sql
└── sqlite/
    ├── 1_create_users_table.up.sql
    └── 1_create_users_table.down.sql
```

**Enable with**: Add `sourceAppendDriver: true` to `.jzero.yaml` (development) OR use `migrate.WithSourceAppendDriver(true)` in code (production)

**Why separate?** Different databases use different SQL syntax:
- MySQL: `AUTO_INCREMENT`, `ENGINE=InnoDB`, backticks
- PostgreSQL: `SERIAL`, no engine, double quotes
- SQLite: `AUTOINCREMENT`, `TEXT` type, no engine

## Naming Convention

- **Format**: `{sequence_number}_{description}.{direction}.sql`
- **Sequence**: 1, 2, 3, ... (consecutive numbers)
- **Description**: Snake_case (e.g., `add_phone_column`)
- **Direction**: `up.sql` (apply) or `down.sql` (rollback)

## Common Migration Patterns

### Create Table

```sql
-- up.sql
CREATE TABLE `users` (
  `id` bigint unsigned NOT NULL AUTO_INCREMENT,
  `username` varchar(50) NOT NULL,
  `email` varchar(100) NOT NULL,
  `created_at` datetime NOT NULL DEFAULT CURRENT_TIMESTAMP,
  PRIMARY KEY (`id`),
  UNIQUE KEY `idx_username` (`username`)
) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4;

-- down.sql
DROP TABLE IF EXISTS `users`;
```

### Add Column

```sql
-- up.sql
ALTER TABLE `users` ADD COLUMN `phone` varchar(20) DEFAULT NULL AFTER `email`;

-- down.sql
ALTER TABLE `users` DROP COLUMN `phone`;
```

### Add Index

```sql
-- up.sql
CREATE INDEX `idx_email` ON `users` (`email`);

-- down.sql
DROP INDEX `idx_email` ON `users`;
```

### Modify Column

```sql
-- up.sql
ALTER TABLE `users` MODIFY COLUMN `username` varchar(100) NOT NULL;

-- down.sql
ALTER TABLE `users` MODIFY COLUMN `username` varchar(50) NOT NULL;
```

## Configuration

### Development Mode (.jzero.yaml)

Required for `jzero migrate` commands:

```yaml
# .jzero.yaml (development only)
migrate:
  driver: "mysql"                    # mysql, pgx, or sqlite
  datasource-url: "root:password@tcp(127.0.0.1:3306)/mydb"
  sourceAppendDriver: true           # Optional: for multi-database support
```

**Database examples**:
- MySQL: `root:password@tcp(127.0.0.1:3306)/mydb`
- PostgreSQL: `postgres://user:password@localhost:5432/mydb?sslmode=disable`
- SQLite: `file:./mydb.db`

### Production Mode (etc/etc.yaml)

Uses existing `Sqlx` configuration - **no migrate section needed**:

```yaml
# etc/etc.yaml (production)
Name: myproject
Host: 0.0.0.0
Port: 8888

Sqlx:
  DataSource: "root:prodpassword@tcp(prod-db:3306)/mydb"
  DriverName: mysql
```

**Code reads from**: `cc.MustGetConfig().Sqlx.SqlConf`

## Schema Generation Modes

### Local SQL Mode
- Schema in `desc/sql/*.sql` files
- Generate models: `jzero gen --desc desc/sql/table.sql`
- When changing: Update `.sql` file + create migrations

### Remote Datasource Mode
- Schema from database connection
- Generate models: `jzero gen --datasource="..."`
- When changing: Create migrations directly (no `.sql` files needed)

## Workflow Example

### Development Workflow

```bash
# 1. Create migration files
echo "ALTER TABLE users ADD COLUMN phone varchar(20);" > desc/sql_migration/4_add_phone.up.sql
echo "ALTER TABLE users DROP COLUMN phone;" > desc/sql_migration/4_add_phone.down.sql

# 2. Apply migrations (development only)
jzero migrate up

# 3. Regenerate models (Local SQL Mode)
jzero gen --desc desc/sql/users.sql

# OR (Remote Datasource Mode)
jzero gen --datasource="root:password@tcp(127.0.0.1:3306)/mydb"

# 4. Test and rollback if needed
jzero migrate down
```

### Production Workflow

```bash
# 1. Create migration files (same as development)
# 2. Ensure cmd/server.go has migration code (see Production Mode above)
# 3. Deploy - migrations run automatically on server startup!
./myproject server --config etc/etc.yaml
```

## Best Practices

### ✅ Always Follow

1. **Create both up and down migrations** - Every change must be reversible
2. **Test down migrations** - Verify rollback works
3. **Use descriptive names** - Make files self-documenting
4. **Keep migrations focused** - One logical change per migration
5. **Run migrations before creating ServiceContext** - Ensures models match schema
6. **Use sequence numbers correctly** - Maintain consecutive ordering (1, 2, 3...)

### ❌ Never Do

1. **Never modify existing migration files** - Create new ones instead
2. **Never skip down migrations** - Always provide rollback logic
3. **Never use destructive operations in down migrations** - Down should safely undo
4. **Never mix multiple unrelated changes** - Keep each migration focused
5. **Never skip sequence numbers** - Maintain consecutive numbering

### When to Create Migrations

- ✅ Adding/modifying/dropping tables
- ✅ Adding/dropping columns
- ✅ Adding/dropping indexes or constraints
- ✅ Adding foreign keys
- ✅ Renaming tables or columns

## Multi-Database Decision Guide

| Factor | Single Database | Multi-Database |
|--------|----------------|----------------|
| Directory | `desc/sql_migration/` | `desc/sql_migration/{driver}/` |
| Configuration | Default | Add `sourceAppendDriver: true` |
| Maintenance | One set of migrations | Separate migrations per DB |
| SQL syntax | Single database type | Database-specific for each |
| Best for | One database type | Dev=SQLite, Prod=MySQL/PostgreSQL |

**Use multi-database when**:
- Development uses SQLite, production uses MySQL/PostgreSQL
- Need to test with different databases
- Want database-specific SQL syntax optimization

## Troubleshooting

### Migration files not found

**Check directory structure**:
- Single DB mode: Files in `desc/sql_migration/`
- Multi-DB mode: Files in `desc/sql_migration/{driver}/`

**Verify configuration**:
```yaml
# .jzero.yaml
migrate:
  driver: "mysql"
  datasource-url: "..."
  sourceAppendDriver: true  # Matches your directory structure?
```

### Wrong database syntax errors

Ensure SQL syntax matches your database:
- MySQL: Backticks, `AUTO_INCREMENT`, `ENGINE=InnoDB`
- PostgreSQL: No quotes (or double quotes), `SERIAL`, no `ENGINE`
- SQLite: No quotes, `AUTOINCREMENT`, `TEXT` type

### Do I need .jzero.yaml in production?

**NO!** Production uses code-based migration with `etc/etc.yaml`. `.jzero.yaml` is development-only.

### Should I use jzero migrate or code migration?

- **Development**: `jzero migrate` - quick testing, easy rollback
- **Production**: Code in `cmd/server.go` - automatic, reliable

## Related Documentation

- [Database Best Practices](./best-practices.md) - Database operation guidelines
- [Model Generation](./model-generation.md) - Generating models from SQL
- [Database Connection](./database-connection.md) - Setting up database connections
