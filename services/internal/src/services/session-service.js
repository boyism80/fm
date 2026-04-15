"use strict";

const SessionState = {
    LOGIN: "LOGIN",
    TRANSITION: "TRANSITION",
    GAME: "GAME",
};

function formatDateTime(date = new Date()) {
    const pad = (v) => String(v).padStart(2, "0");
    const yyyy = date.getFullYear();
    const MM = pad(date.getMonth() + 1);
    const dd = pad(date.getDate());
    const HH = pad(date.getHours());
    const mm = pad(date.getMinutes());
    const ss = pad(date.getSeconds());
    return `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
}

class SessionService {
    constructor(internalContext, sessionRepository) {
        this.ctx = internalContext;
        this.repo = sessionRepository;
    }

    _now() {
        return formatDateTime(new Date());
    }

    _ttlByState(state) {
        switch (state) {
        case SessionState.LOGIN:
            return 120;
        case SessionState.TRANSITION:
            return 30;
        case SessionState.GAME:
            return 180;
        default:
            return 120;
        }
    }

    async beginLogin(worldId, accountId, loginServerId) {
        const now = this._now();
        const [ok, code] = await this.repo.beginLoginAtomic(
            worldId,
            accountId,
            loginServerId,
            now,
            this._ttlByState(SessionState.LOGIN)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: code || "ALREADY_LOGGED_IN" };
        }
        return { ok: true };
    }

    async beginTransition(worldId, accountId, characterId, characterName) {
        const now = this._now();
        const [ok, code] = await this.repo.beginTransitionAtomic(
            worldId,
            accountId,
            characterId,
            characterName,
            now,
            this._ttlByState(SessionState.TRANSITION)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: code || "SESSION_NOT_FOUND" };
        }
        return { ok: true };
    }

    async attachGameSession(worldId, accountId, channelId) {
        const now = this._now();
        const [ok, code] = await this.repo.attachGameSessionAtomic(
            worldId,
            accountId,
            channelId,
            now,
            this._ttlByState(SessionState.GAME)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: code || "SESSION_NOT_FOUND" };
        }
        return { ok: true };
    }

    async refresh(worldId, accountId) {
        const [ok, code] = await this.repo.refreshAtomic(
            worldId,
            accountId,
            this._ttlByState(SessionState.LOGIN),
            this._ttlByState(SessionState.TRANSITION),
            this._ttlByState(SessionState.GAME)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: code || "SESSION_NOT_FOUND" };
        }
        return { ok: true };
    }

    async logout(worldId, accountId) {
        const [ok, code] = await this.repo.logoutAtomic(
            worldId,
            accountId,
            this._ttlByState(SessionState.TRANSITION)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: code || "LOGOUT_FAILED" };
        }
        return { ok: true, code: code || "" };
    }
}

module.exports = { SessionService, SessionState, formatDateTime };

