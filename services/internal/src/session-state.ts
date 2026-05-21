import { AccountSessionState, CharacterSessionState } from "./protobuf/generated/fminternal/internal_service";

export { AccountSessionState, CharacterSessionState };

const VALID_ACCOUNT_SESSION_STATES: ReadonlySet<AccountSessionState> = new Set([
    AccountSessionState.ACCOUNT_SESSION_STATE_UNSPECIFIED,
    AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN,
    AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION,
    AccountSessionState.ACCOUNT_SESSION_STATE_GAME,
]);

const VALID_CHARACTER_SESSION_STATES: ReadonlySet<CharacterSessionState> = new Set([
    CharacterSessionState.CHARACTER_SESSION_STATE_TRANSITION,
    CharacterSessionState.CHARACTER_SESSION_STATE_ONLINE,
]);

export function accountSessionStateFromRedisHash(value: string | undefined): AccountSessionState {
    const n = Number(value ?? 0);
    if (VALID_ACCOUNT_SESSION_STATES.has(n)) {
        return n;
    }
    return AccountSessionState.ACCOUNT_SESSION_STATE_UNSPECIFIED;
}

export function accountSessionStateToRedisHash(state: AccountSessionState): string {
    return String(state);
}

export function isValidCharacterSessionState(value: number): value is CharacterSessionState {
    return Number.isInteger(value) && VALID_CHARACTER_SESSION_STATES.has(value);
}
