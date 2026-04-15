"use strict";

class UnifiedRepository {
    constructor(internalContext) {
        this.ctx = internalContext;
    }

    _pool() {
        return this.ctx.getPgUnifiedPool();
    }

    async findAccountByLoginId(loginId) {
        const { rows } = await this._pool().query(
            "SELECT id, login_id, created_at FROM account_identity WHERE login_id = $1",
            [loginId]
        );
        return rows[0] ?? null;
    }

    async insertAccountIdentity(loginId) {
        const { rows } = await this._pool().query(
            "INSERT INTO account_identity (login_id) VALUES ($1) RETURNING id, login_id, created_at",
            [loginId]
        );
        return rows[0];
    }

    async findCharacterNameEntry(name) {
        const { rows } = await this._pool().query(
            "SELECT character_id, name, account_id, world_id, status FROM character_name_registry WHERE LOWER(name) = LOWER($1) AND deleted_at IS NULL",
            [name]
        );
        return rows[0] ?? null;
    }

    async reserveCharacterName(name, accountId, worldId) {
        const { rows } = await this._pool().query(
            `INSERT INTO character_name_registry (name, account_id, world_id)
             VALUES ($1, $2, $3)
             RETURNING character_id, name, account_id, world_id`,
            [name, accountId, worldId]
        );
        return rows[0];
    }

    async deleteCharacterName(characterId) {
        await this._pool().query(
            "UPDATE character_name_registry SET deleted_at = NOW() WHERE character_id = $1",
            [characterId]
        );
    }
}

module.exports = { UnifiedRepository };
