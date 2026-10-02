import { AccountSessionState } from "./protobuf/generated/fminternal/internal_service";

export { AccountSessionState };

const VALID_ACCOUNT_SESSION_STATES: ReadonlySet<AccountSessionState> = new Set([
    AccountSessionState.ACCOUNT_SESSION_STATE_UNSPECIFIED,
    AccountSessionState.ACCOUNT_SESSION_STATE_LOGIN,
    AccountSessionState.ACCOUNT_SESSION_STATE_TRANSITION,
    AccountSessionState.ACCOUNT_SESSION_STATE_GAME,
]);

export function accountSessionStateFromRedisHash(value: string | undefined): AccountSessionState {
    const n = Number(value ?? 0);
    if (VALID_ACCOUNT_SESSION_STATES.has(n)) {
        return n;
    }
    return AccountSessionState.ACCOUNT_SESSION_STATE_UNSPECIFIED;
}
