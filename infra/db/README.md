# Sequelize migrations (PostgreSQL)

Migrations live here, outside any service. The runner is in `@fm/shared` (`sequelize-migrator.ts`, `auto-migrate.ts`), so every TypeScript service that uses the databases can apply them without depending on another service.

Each **data shard** database must run the same migrations.

```bash
cd services/shared
# Windows PowerShell example — point at one shard:
$env:DATABASE_URL = "postgres://fm:admin@127.0.0.1:5432/fm"
npm run migrate
```

Repeat with a different `DATABASE_URL` for every shard. The first migration **drops** `characters` if present (see migration file) to align schema; do not run that blindly on production with data.

## Auto migration on server startup

Set `sequelize.auto_migrate_on_startup: true` in the service config (`services/internal/config.yaml`).
When enabled, the service calls `autoMigrateAllIfEnabled` on startup and applies pending migrations to all configured PostgreSQL endpoints. A PostgreSQL advisory lock per database lets several services start at once:
- `postgresql.unified`
- each world `postgresql.worlds.<id>.global`
- each world `postgresql.worlds.<id>.data[]`


## Social (#277–#281)

`20261011120000-social.ts` adds `characters.family_progress` and `characters.population_history` (JSONB: last given time and per-target times) on data shards. The `social_*` tables are used in each world's global database; startup migrations currently provision the same schema on all configured endpoints. Apply migrations before starting the updated internal/game services.

`20261011130000-social-durability.ts` adds `characters.family_reputation`, plus family recipient credits and movement reservations. Notes stay in global `social_notes`; the game server keeps a synced copy and grants thank-you population from the remove reply (a lost reply loses it until #298). Family credits keep their recipients from the time EXP was earned and replay only the unsaved difference. Teleport/summon reservations become final on destination entry; unconfirmed reservations expire after 30 seconds and return reputation and the original daily use, even across midnight. Internal sweeps expired reservations every 10 seconds, including those left by a service restart.

`20261021100000-family-leader.ts` adds `families` (leader_id, precept) and `social_family.leader_id`, the leader's character id or NULL outside a family, and moves the leader's notice to `families.precept`. A family is changed only while its `families` row is locked. `social_family_invites` and `social_family_summons` are dropped: invites and summons are held by the answering character on the game server.

Before rollout, verify with two channels and three characters: cross-channel whisper/find and offline targets; invite/join/decline/leave in a three-person messenger; per-character 24-hour population and 30-day target cooldown across reconnect; offline notes and repeated/foreign note IDs; cash-gift thank-you and repeated 0x6E claims; family invite, two-junior limit, subtree merge/split, precept permissions, reputation after EXP/save, once-daily perks, blocked-map teleport/summon, and perk expiry/reconnect.

Packet layouts were checked against the active client with IDA: whisper send `0x4F64CB`, population result `0x84F121`, note receive `0x85177C`/record `0x4BD87B`, note deletion `0x590C87`, messenger join `0x6D2AA3`/look update `0x6D537E`, family pedigree `0x6ADA77` and info `0x6ACA79`. Ordinary note success is action 4; action 2 loads the inbox; action 7 requests a refresh. Messenger channel and notification are separate bytes. Family pedigree ends with two counted DWORD-pair arrays and a WORD capacity.
