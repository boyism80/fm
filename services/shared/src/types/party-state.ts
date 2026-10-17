import { PartyState } from "../protobuf/generated/fminternal/internal_service";

export { PartyState };

export const DEFAULT_PARTY_STATE = PartyState.PARTY_STATE_ACTIVE;

const VALID_PARTY_STATES: ReadonlySet<PartyState> = new Set([
    PartyState.PARTY_STATE_UNSPECIFIED,
    PartyState.PARTY_STATE_ACTIVE,
]);

export function partyStateFromDb(value: number | null | undefined): PartyState {
    if (value != null && VALID_PARTY_STATES.has(value)) {
        return value;
    }
    return DEFAULT_PARTY_STATE;
}

export function partyStateToDb(state: PartyState): number {
    return state;
}
