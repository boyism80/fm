import type { ServiceError } from "@grpc/grpc-js";

export type GrpcCallback<TReply> = (error: ServiceError | null, reply?: TReply) => void;

export type GrpcCall<TRequest> = { request: TRequest; cancelled: boolean };

// Services answer { ok, code, ...fields }; a reply carries the same fields with the code as errorCode.
export function okReply<TResult extends { ok: boolean; code?: unknown }>(
    result: TResult,
    none: NonNullable<TResult["code"]>,
    unknown: NonNullable<TResult["code"]>,
) {
    const { code, ...fields } = result;
    const errorCode = (result.ok ? none : (code ?? unknown)) as NonNullable<TResult["code"]>;
    return { ...fields, errorCode };
}
