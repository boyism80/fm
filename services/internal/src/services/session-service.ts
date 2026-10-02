import { AccountSessionState, SessionDisconnectSource, SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import { isRedisEvalArray } from "../types/redis-eval";
import type { SessionRepository } from "../repos/session-repository";
import type { PartyService } from "./party-service";
import type { BuddyService } from "./buddy-service";
import type { GuildService } from "./guild-service";

export function formatDateTime(date = new Date()) {
    const pad = (v: number) => String(v).padStart(2, "0");
    const yyyy = date.getFullYear();
    const MM = pad(date.getMonth() + 1);
    const dd = pad(date.getDate());
    const HH = pad(date.getHours());
    const mm = pad(date.getMinutes());
    const ss = pad(date.getSeconds());
    return `${yyyy}-${MM}-${dd} ${HH}:${mm}:${ss}`;
}

export class SessionService {
    private readonly repo: SessionRepository;
    private readonly partyService: PartyService | null;
    private readonly buddyService: BuddyService | null;
    private readonly guildService: GuildService | null;

    constructor(
        sessionRepository: SessionRepository,
        partyService: PartyService | null,
        buddyService: BuddyService | null,
        guildService: GuildService | null
    ) {
        this.repo = sessionRepository;
        this.partyService = partyService;
        this.buddyService = buddyService;
        this.guildService = guildService;
    }

    private now() {
        return formatDateTime(new Date());
    }

    private ttlByState(state: AccountSessionState) {
        switch (state) {
            case AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN:
                return 120;
            case AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION:
                return 30;
            case AccountSessionState.ACCOUNT_SESSION_STATE_GAME:
                return 180;
            default:
                return 120;
        }
    }

    private atomicResultTuple(raw: Array<number | string> | null | undefined): [number, number] {
        if (raw != null && isRedisEvalArray(raw)) {
            const ok = Number(raw[0] ?? 0);
            const code = Number(raw[1] ?? SessionErrorCode.SESSION_UNKNOWN);
            return [ok, code];
        }
        return [0, SessionErrorCode.SESSION_UNKNOWN];
    }

    async beginLogin(worldId: number, accountId: number, loginServerId: string) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.beginLogin(
                worldId,
                accountId,
                loginServerId,
                this.now(),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN)
            )
        );
        if (ok !== 1) {
            const existing = await this.repo.getAccountSession(worldId, accountId);
            const ttl = await this.repo.getAccountSessionTtl(worldId, accountId);
            console.log(
                `[session] begin_login rejected account=${accountId} login_server=${loginServerId} code=${code}` +
                    ` existing_state=${existing?.state} existing_login_server=${existing?.loginServer?.id}` +
                    ` login_connected=${existing?.loginServer?.connected} created_at=${existing?.timestamps?.createdAt}` +
                    ` updated_at=${existing?.timestamps?.updatedAt} ttl=${ttl}`
            );
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_ALREADY_LOGGED_IN };
        }
        console.log(`[session] begin_login ok account=${accountId} login_server=${loginServerId}`);
        return { ok: true };
    }

    async beginTransition(worldId: number, accountId: number, characterId: number, characterName: string) {
        const account = await this.repo.getAccountSession(worldId, accountId);
        const gameToGameTransfer = account?.state === AccountSessionState.ACCOUNT_SESSION_STATE_GAME;
        const [ok, code] = this.atomicResultTuple(
            await this.repo.beginTransition(
                worldId,
                accountId,
                characterId,
                characterName,
                this.now(),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION),
                gameToGameTransfer,
                account?.channelId ?? null
            )
        );
        if (ok !== 1) {
            console.log(`[session] begin_transition failed account=${accountId} character=${characterId} code=${code}`);
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        console.log(`[session] begin_transition ok account=${accountId} character=${characterId} game_to_game=${gameToGameTransfer}`);
        return { ok: true };
    }

    async enterGame(worldId: number, accountId: number, characterId: number, channelId: number) {
        const account = await this.repo.getAccountSession(worldId, accountId);
        const gameToGameTransfer = account?.gameToGameTransfer === true;
        const [ok, code] = this.atomicResultTuple(
            await this.repo.enterGame(
                worldId,
                accountId,
                characterId,
                channelId,
                this.now(),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_GAME)
            )
        );
        if (ok !== 1) {
            console.log(`[session] enter_game failed account=${accountId} character=${characterId} channel=${channelId} code=${code}`);
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        console.log(`[session] enter_game ok account=${accountId} character=${characterId} channel=${channelId}`);
        try {
            if (this.partyService) {
                await this.partyService.applyMemberChannelIndex(worldId, characterId, channelId);
            }
            if (this.buddyService) {
                await this.buddyService.applyBuddyChannelIndex(worldId, characterId, channelId);
            }
            if (this.guildService && !gameToGameTransfer) {
                await this.guildService.applyMemberOnlineState(worldId, characterId, true);
            }
        } catch (err) {
            // The caller answers with an error, which the game server reads as "did not enter the game".
            await this.logout(worldId, accountId, {
                disconnectSource: SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER,
                characterId,
                channelId,
            }).catch((logoutErr) => {
                console.log(`[session] enter_game rollback failed account=${accountId} character=${characterId}: ${logoutErr}`);
            });
            throw err;
        }
        return { ok: true };
    }

    async refresh(worldId: number, accountId: number) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.refresh(
                worldId,
                accountId,
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_GAME)
            )
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { ok: true };
    }

    async logout(
        worldId: number,
        accountId: number,
        options: {
            disconnectSource?: SessionDisconnectSource;
            transferDisconnect?: boolean;
            characterId?: number;
            channelId?: number;
        } = {}
    ) {
        const src = options.disconnectSource ?? SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_UNSPECIFIED;
        const transferDisconnect = options.transferDisconnect ?? false;
        const keep = src === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_LOGIN_SERVER && transferDisconnect;
        let channelId = options.channelId;
        if (keep === false && (channelId === undefined || !Number.isInteger(channelId) || channelId < 0)) {
            const account = await this.repo.getAccountSession(worldId, accountId);
            const sch = account?.channelId;
            if (sch != null && Number.isInteger(sch) && sch >= 0) {
                channelId = sch;
            }
        }
        const raw = await this.repo.logout(worldId, accountId, keep, channelId);
        const [ok, code] = this.atomicResultTuple(raw);
        console.log(
            `[session] logout account=${accountId} source=${src} transfer=${transferDisconnect}` +
                ` character=${options.characterId ?? ""} ok=${ok} prev_state=${Array.isArray(raw) ? raw[2] : ""}`
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_LOGOUT_FAILED };
        }
        const gameNormalDisconnect = src === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER && !transferDisconnect;
        const cid = options.characterId ?? null;
        if (cid != null && gameNormalDisconnect) {
            if (this.partyService) {
                await this.partyService.applyMemberChannelIndex(worldId, cid, -2);
            }
            if (this.buddyService) {
                await this.buddyService.applyBuddyChannelIndex(worldId, cid, -1);
            }
            if (this.guildService) {
                await this.guildService.applyMemberOnlineState(worldId, cid, false);
            }
        }
        return { ok: true, code: SessionErrorCode.SESSION_NONE };
    }
}
