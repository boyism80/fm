"use strict";

const messages = require("../../protobuf/fminternal/internal_service_pb");

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
    constructor(
        internalContext,
        sessionRepository,
        characterRealtimeStateRepository,
        partyService
    ) {
        this.ctx = internalContext;
        this.repo = sessionRepository;
        this.characterRealtimeStateRepo = characterRealtimeStateRepository;
        this.partyService = partyService;
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
            return { ok: false, code: Number.isInteger(code) ? code : messages.SessionErrorCode.SESSION_ALREADY_LOGGED_IN };
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
            return { ok: false, code: Number.isInteger(code) ? code : messages.SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { ok: true };
    }

    async attachGameSession(worldId, accountId, channelId) {
        const now = this._now();
        const sessionBeforeAttach = await this.repo.getAccountSession(worldId, accountId);
        const [ok, code] = await this.repo.attachGameSessionAtomic(
            worldId,
            accountId,
            channelId,
            now,
            this._ttlByState(SessionState.GAME)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : messages.SessionErrorCode.SESSION_NOT_FOUND };
        }
        const cid = sessionBeforeAttach?.character?.id;
        if (cid != null && this.partyService) {
            await this.partyService.applyMemberChannelIndex(worldId, cid, channelId);
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
            return { ok: false, code: Number.isInteger(code) ? code : messages.SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { ok: true };
    }

    async logout(worldId, accountId) {
        const sessionBeforeLogout = await this.repo.getAccountSession(worldId, accountId);
        const [ok, code] = await this.repo.logoutAtomic(
            worldId,
            accountId,
            this._ttlByState(SessionState.TRANSITION)
        );
        if (Number(ok) !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : messages.SessionErrorCode.SESSION_LOGOUT_FAILED };
        }
        const cid = sessionBeforeLogout?.character?.id;
        if (cid != null && this.partyService) {
            await this.partyService.applyMemberChannelIndex(worldId, cid, -2);
        }
        return { ok: true, code: messages.SessionErrorCode.SESSION_NONE };
    }
}

module.exports = { SessionService, SessionState, formatDateTime };

