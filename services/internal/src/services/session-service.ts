import { AccountSessionState, SessionDisconnectSource, SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { DebuffPersisted } from "../protobuf/generated/fminternal/internal_service";
import { isRedisEvalArray } from "../types/redis-eval";
import type { SessionRepository } from "../repos/session-repository";
import type { UnifiedRepository } from "../repos/unified-repository";
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
    private readonly unifiedRepo: UnifiedRepository;
    private readonly partyService: PartyService | null;
    private readonly buddyService: BuddyService | null;
    private readonly guildService: GuildService | null;

    constructor(
        sessionRepository: SessionRepository,
        unifiedRepository: UnifiedRepository,
        partyService: PartyService | null,
        buddyService: BuddyService | null,
        guildService: GuildService | null
    ) {
        this.repo = sessionRepository;
        this.unifiedRepo = unifiedRepository;
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

    async beginTransition(
        worldId: number,
        accountId: number,
        characterId: number,
        characterName: string,
        clientIp: string,
        debuffs: DebuffPersisted[],
        sourceChannelId: number | null
    ) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.beginTransition(
                worldId,
                accountId,
                characterId,
                characterName,
                clientIp,
                this.now(),
                this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION),
                sourceChannelId,
                debuffs
            )
        );
        if (ok !== 1) {
            console.log(`[session] begin_transition failed account=${accountId} character=${characterId} code=${code}`);
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        console.log(`[session] begin_transition ok account=${accountId} character=${characterId} source_channel=${sourceChannelId ?? ""}`);
        await this.releaseCharacterNameReservation(accountId);
        return { ok: true };
    }

    async enterGame(worldId: number, accountId: number, characterId: number, channelId: number, clientIp: string) {
        const account = await this.repo.getAccountSession(worldId, accountId);
        const gameToGameTransfer = account?.gameToGameTransfer === true;
        const [ok, code] = this.atomicResultTuple(
            await this.repo.enterGame(
                worldId,
                accountId,
                characterId,
                channelId,
                clientIp,
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
                await this.partyService.publishMemberLogOnOff(worldId, characterId);
            }
            if (this.buddyService) {
                await this.buddyService.publishBuddyChannel(worldId, characterId, channelId);
            }
            if (this.guildService && !gameToGameTransfer) {
                await this.guildService.publishMemberOnline(worldId, characterId, true);
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
        return { ok: true, debuffs: gameToGameTransfer ? account?.debuffs ?? [] : [] };
    }

    async refresh(worldId: number, accountId: number, owner: { characterId: number; channelId: number } | null) {
        const raw = await this.repo.refresh(
            worldId,
            accountId,
            this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN),
            this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION),
            this.ttlByState(AccountSessionState.ACCOUNT_SESSION_STATE_GAME),
            owner
        );
        const [ok, code] = this.atomicResultTuple(raw);
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        if (Number(raw[3]) === AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN) {
            await this.unifiedRepo.refreshCharacterNameReservation(accountId).catch((err) => {
                console.log(`[session] refresh name reservation failed account=${accountId}: ${err}`);
            });
        }
        return { ok: true };
    }

    // The reservation also expires on its own, so a failure here only delays the release.
    private async releaseCharacterNameReservation(accountId: number) {
        await this.unifiedRepo.releaseCharacterNameReservation(accountId).catch((err) => {
            console.log(`[session] release name reservation failed account=${accountId}: ${err}`);
        });
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
        const ownerCharacterId = src === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER ? options.characterId : undefined;
        const raw = await this.repo.logout(worldId, accountId, keep, channelId, ownerCharacterId);
        const [ok, code] = this.atomicResultTuple(raw);
        console.log(
            `[session] logout account=${accountId} source=${src} transfer=${transferDisconnect}` +
                ` character=${options.characterId ?? ""} ok=${ok} prev_state=${Array.isArray(raw) ? raw[2] : ""}`
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_LOGOUT_FAILED };
        }
        await this.releaseCharacterNameReservation(accountId);
        const gameNormalDisconnect = src === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER && !transferDisconnect;
        const cid = options.characterId ?? null;
        if (cid != null && gameNormalDisconnect) {
            if (this.partyService) {
                await this.partyService.publishMemberLogOnOff(worldId, cid);
            }
            if (this.buddyService) {
                await this.buddyService.publishBuddyChannel(worldId, cid, -1);
            }
            if (this.guildService) {
                await this.guildService.publishMemberOnline(worldId, cid, false);
            }
        }
        return { ok: true, code: SessionErrorCode.SESSION_NONE };
    }
}
