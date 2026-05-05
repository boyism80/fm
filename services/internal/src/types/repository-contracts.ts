import type { PoolClient } from "pg";

export type RepositoryQuery = { text: string; values: unknown[] };
export type RepositoryTxOptions = { txClient?: PoolClient };
