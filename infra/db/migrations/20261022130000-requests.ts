import type { MigrationModule } from "@fm/shared/sequelize-migration";

const migration: MigrationModule = {
    async up(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE IF NOT EXISTS requests (
                request_id TEXT NOT NULL,
                kind TEXT NOT NULL,
                owner_id INTEGER NOT NULL,
                status TEXT NOT NULL DEFAULT 'applied',
                outcome JSONB NOT NULL DEFAULT '{}'::jsonb,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW(),
                PRIMARY KEY (request_id, kind)
            );

            INSERT INTO requests (request_id, kind, owner_id, status, created_at)
            SELECT request_id::text, 'parcel_send', receiver_id, CASE WHEN aborted THEN 'aborted' ELSE 'applied' END, created_at
            FROM parcel_send_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, outcome, created_at)
            SELECT request_id::text, 'parcel_claim', receiver_id, jsonb_build_object('parcel', parcel), created_at
            FROM parcel_claim_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, outcome, created_at)
            SELECT request_id::text, 'note_remove', character_id, jsonb_build_object('population', population), created_at
            FROM note_remove_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, status, created_at)
            SELECT request_id::text, 'wedding_gift_send', receiver_id, CASE WHEN aborted THEN 'aborted' ELSE 'applied' END, created_at
            FROM wedding_gift_send_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, outcome, created_at)
            SELECT request_id::text, 'wedding_gift_claim', receiver_id, jsonb_build_object('gift', gift), created_at
            FROM wedding_gift_claim_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, status, created_at)
            SELECT request_id::text, 'marriage', 0, CASE WHEN aborted THEN 'aborted' ELSE 'applied' END, created_at
            FROM marriage_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, status, outcome)
            SELECT operation_id, 'family_move', character_id,
                CASE status WHEN 'pending' THEN 'applied' WHEN 'completed' THEN 'completed' ELSE 'aborted' END,
                jsonb_build_object('payer_id', payer_id, 'feature', feature)
            FROM social_family_moves ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, status, outcome, created_at)
            SELECT request_id::text, 'shop_open', account_id, CASE WHEN aborted THEN 'aborted' ELSE 'applied' END,
                CASE WHEN shop_id IS NULL THEN '{}'::jsonb ELSE jsonb_build_object('shop_id', shop_id) END, created_at
            FROM shop_open_requests ON CONFLICT DO NOTHING;

            INSERT INTO requests (request_id, kind, owner_id, created_at)
            SELECT r.request_id, 'cash', COALESCE(t.target_id, c.used_by, 0), r.created_at
            FROM cash_transfer_receipts r
            LEFT JOIN cash_transfer_requests t ON t.request_id::text = r.request_id
            LEFT JOIN cash_coupons c ON c.request_id = r.request_id
            ON CONFLICT DO NOTHING;

            DROP TABLE parcel_send_requests, parcel_claim_requests, note_remove_requests, wedding_gift_send_requests,
                wedding_gift_claim_requests, guild_requests, marriage_requests, social_family_moves, shop_open_requests,
                cash_transfer_receipts;
        `);
    },

    async down(queryInterface) {
        await queryInterface.sequelize.query(`
            CREATE TABLE parcel_send_requests (
                request_id BIGINT PRIMARY KEY, receiver_id INTEGER NOT NULL, aborted BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO parcel_send_requests SELECT request_id::bigint, owner_id, status = 'aborted', created_at FROM requests WHERE kind = 'parcel_send';

            CREATE TABLE parcel_claim_requests (
                request_id BIGINT PRIMARY KEY, receiver_id INTEGER NOT NULL, parcel JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO parcel_claim_requests SELECT request_id::bigint, owner_id, outcome->'parcel', created_at FROM requests WHERE kind = 'parcel_claim';

            CREATE TABLE note_remove_requests (
                request_id BIGINT PRIMARY KEY, character_id INTEGER NOT NULL, population INTEGER NOT NULL,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO note_remove_requests SELECT request_id::bigint, owner_id, (outcome->>'population')::int, created_at FROM requests WHERE kind = 'note_remove';

            CREATE TABLE wedding_gift_send_requests (
                request_id BIGINT PRIMARY KEY, receiver_id INTEGER NOT NULL, aborted BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO wedding_gift_send_requests SELECT request_id::bigint, owner_id, status = 'aborted', created_at FROM requests WHERE kind = 'wedding_gift_send';

            CREATE TABLE wedding_gift_claim_requests (
                request_id BIGINT PRIMARY KEY, receiver_id INTEGER NOT NULL, gift JSONB NOT NULL, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO wedding_gift_claim_requests SELECT request_id::bigint, owner_id, outcome->'gift', created_at FROM requests WHERE kind = 'wedding_gift_claim';

            CREATE TABLE guild_requests (
                request_id BIGINT PRIMARY KEY, owner_id INTEGER NOT NULL, result_id INTEGER, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO guild_requests
            SELECT request_id::bigint, owner_id, COALESCE((outcome->>'guild_id')::int, (outcome->>'alliance_id')::int), created_at
            FROM requests WHERE kind IN ('guild_create', 'guild_change', 'alliance_create', 'alliance_change');

            CREATE TABLE marriage_requests (
                request_id BIGINT PRIMARY KEY, aborted BOOLEAN NOT NULL DEFAULT FALSE, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            INSERT INTO marriage_requests SELECT request_id::bigint, status = 'aborted', created_at FROM requests WHERE kind = 'marriage';

            CREATE TABLE social_family_moves (
                operation_id TEXT PRIMARY KEY, character_id INTEGER NOT NULL, payer_id INTEGER NOT NULL, feature INTEGER NOT NULL,
                status TEXT NOT NULL DEFAULT 'pending', expires_at TIMESTAMPTZ NOT NULL DEFAULT NOW()+INTERVAL '30 seconds'
            );
            INSERT INTO social_family_moves (operation_id, character_id, payer_id, feature, status)
            SELECT request_id, owner_id, COALESCE((outcome->>'payer_id')::int, 0), COALESCE((outcome->>'feature')::int, 0),
                CASE status WHEN 'applied' THEN 'pending' WHEN 'completed' THEN 'completed' ELSE 'refunded' END
            FROM requests WHERE kind = 'family_move';

            CREATE TABLE shop_open_requests (
                request_id BIGINT PRIMARY KEY, account_id INTEGER NOT NULL, shop_id BIGINT, aborted BOOLEAN NOT NULL DEFAULT FALSE,
                created_at TIMESTAMPTZ NOT NULL DEFAULT NOW()
            );
            CREATE INDEX idx_shop_open_requests_account ON shop_open_requests (account_id);
            INSERT INTO shop_open_requests SELECT request_id::bigint, owner_id, (outcome->>'shop_id')::bigint, status = 'aborted', created_at FROM requests WHERE kind = 'shop_open';

            CREATE TABLE cash_transfer_receipts (request_id TEXT PRIMARY KEY, created_at TIMESTAMPTZ NOT NULL DEFAULT NOW());
            INSERT INTO cash_transfer_receipts SELECT request_id, created_at FROM requests WHERE kind = 'cash';

            DROP TABLE requests;
        `);
    },
};

export default migration;
