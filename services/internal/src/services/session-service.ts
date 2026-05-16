import { SessionDisconnectSource, SessionErrorCode } from "../protobuf/generated/fminternal/internal_service";
import type { InternalContext } from "../context/internal-context";
import type { CharacterRealtimeStateRepository } from "../repos/character-realtime-state-repository";
import type { SessionRepository } from "../repos/session-repository";
import type { CharacterService } from "./character-service";
import type { PartyService } from "./party-service";

export const SessionState = {
    LOGIN: "LOGIN",
    TRANSITION: "TRANSITION",
    GAME: "GAME",
};

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
    private readonly characterService: CharacterService;
    private readonly partyService: PartyService | null;

    constructor(
        internalContext: InternalContext,
        sessionRepository: SessionRepository,
        characterRealtimeStateRepository: CharacterRealtimeStateRepository,
        characterService: CharacterService,
        partyService: PartyService | null
    ) {
        void internalContext;
        void characterRealtimeStateRepository;
        this.repo = sessionRepository;
        this.characterService = characterService;
        this.partyService = partyService;
    }

    private now() {
        return formatDateTime(new Date());
    }

    private ttlByState(state: string) {
        switch (state) {
            case SessionState.LOGIN: return 120;
            case SessionState.TRANSITION: return 30;
            case SessionState.GAME: return 180;
            default: return 120;
        }
    }

    private atomicResultTuple(raw: unknown): [number, number] {
        if (Array.isArray(raw)) {
            const ok = Number(raw[0] ?? 0);
            const code = Number(raw[1] ?? SessionErrorCode.SESSION_UNKNOWN);
            return [ok, code];
        }
        return [0, SessionErrorCode.SESSION_UNKNOWN];
    }

    async beginLogin(worldId: number, accountId: number, loginServerId: string) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.beginLoginAtomic(worldId, accountId, loginServerId, this.now(), this.ttlByState(SessionState.LOGIN))
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_ALREADY_LOGGED_IN };
        }
        return { ok: true };
    }

    async beginTransition(worldId: number, accountId: number, characterId: number, characterName: string) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.beginTransitionAtomic(worldId, accountId, characterId, characterName, this.now(), this.ttlByState(SessionState.TRANSITION))
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        return { ok: true };
    }

    async attachGameSession(
        worldId: number,
        accountId: number,
        characterId: number,
        characterName: string,
        channelId: number
    ) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.attachGameSessionAtomic(
                worldId,
                accountId,
                characterId,
                characterName,
                channelId,
                this.now(),
                this.ttlByState(SessionState.GAME)
            )
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_NOT_FOUND };
        }
        if (this.partyService) await this.partyService.applyMemberChannelIndex(worldId, characterId, channelId);
        return { ok: true };
    }

    async refresh(worldId: number, accountId: number, characterId?: number) {
        const [ok, code] = this.atomicResultTuple(
            await this.repo.refreshAtomic(
                worldId,
                accountId,
                this.ttlByState(SessionState.LOGIN),
                this.ttlByState(SessionState.TRANSITION),
                this.ttlByState(SessionState.GAME),
                characterId
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
            characterName?: string;
            channelId?: number;
        } = {}
    ) {
        let characterName = options.characterName;
        if (
            (characterName == null || characterName === "") &&
            options.characterId != null &&
            options.characterId > 0
        ) {
            const row = await this.characterService.getCharacter(worldId, options.characterId);
            characterName = row?.name;
        }
        let channelId = options.channelId;
        if (
            (channelId === undefined || !Number.isInteger(channelId) || channelId < 0) &&
            characterName != null &&
            characterName !== "" &&
            options.disconnectSource === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER
        ) {
            const sess = await this.repo.getCharacterSessionByName(worldId, characterName);
            const sch = sess?.gameServer?.channelId;
            if (sch != null && Number.isInteger(sch) && sch >= 0) {
                channelId = sch;
            }
        }
        const chForRepo =
            channelId !== undefined && Number.isInteger(channelId) && channelId >= 0 ? channelId : undefined;
        const [ok, code] = this.atomicResultTuple(
            await this.repo.logoutAtomic(
                worldId,
                accountId,
                this.ttlByState(SessionState.TRANSITION),
                characterName,
                chForRepo
            )
        );
        if (ok !== 1) {
            return { ok: false, code: Number.isInteger(code) ? code : SessionErrorCode.SESSION_LOGOUT_FAILED };
        }
        const src = options.disconnectSource ?? SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_UNSPECIFIED;
        const transferDisconnect = options.transferDisconnect ?? false;
        const gameNormalDisconnect = src === SessionDisconnectSource.SESSION_DISCONNECT_SOURCE_GAME_SERVER && !transferDisconnect;
        const cid = options.characterId ?? null;
        if (cid != null && this.partyService && gameNormalDisconnect) await this.partyService.applyMemberChannelIndex(worldId, cid, -2);
        return { ok: true, code: SessionErrorCode.SESSION_NONE };
    }
}
