# Sequelize migrations (PostgreSQL)

Each **data shard** database must run the same migrations.

```bash
cd services/internal
# Windows PowerShell example — point at one shard:
$env:DATABASE_URL = "postgres://fm:admin@127.0.0.1:5432/fm"
npm run migrate
```

Repeat with a different `DATABASE_URL` for every shard. The first migration **drops** `characters` if present (see migration file) to align schema; do not run that blindly on production with data.

Config for CLI: `sequelize/config-for-cli.cjs` (not the YAML app config).

## Auto migration on server startup

Set `sequelize.auto_migrate_on_startup: true` in `services/internal/config.yaml`.
When enabled, `src/server.js` applies pending migrations to all configured PostgreSQL endpoints:
- `postgresql.unified`
- each world `postgresql.worlds.<id>.global`
- each world `postgresql.worlds.<id>.data[]`
