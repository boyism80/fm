import type { PoolClient, QueryResult } from "pg";
import type { EquipmentBonusStatsJson } from "./equipment-bonus-stats";
import type { GuildLogo, GuildRankTitles } from "./guild-json";
import type { KeyLayoutJsonRecord } from "./key-layout-json";

export type RepositoryQueryJson = Record<string, string | number | boolean | null>;

export type PartyMemberDoorQuery = { town: number; target: number; x: number; y: number };

export type CharacterLooksQuery = Record<string, number>;

export type BuffFlagValuesQuery = { mask: number; position: number; value: number }[];

export type RepositoryQueryValue =
    | string
    | number
    | boolean
    | null
    | Date
    | number[]
    | RepositoryQueryJson
    | PartyMemberDoorQuery
    | KeyLayoutJsonRecord
    | EquipmentBonusStatsJson
    | GuildLogo
    | GuildRankTitles
    | CharacterLooksQuery
    | BuffFlagValuesQuery;

export type RepositoryQuery = { text: string; values: RepositoryQueryValue[] };

export type RepositoryTxOptions = { txClient?: PoolClient };

export type RepositoryQueryResult = QueryResult;
