import type { ServiceError } from "@grpc/grpc-js";

export type GrpcCallback<TReply> = (error: ServiceError | null, reply?: TReply) => void;

export type GrpcErrorHandler = (err: unknown, callback: (error: ServiceError) => void) => void;

export type GrpcCall<TRequest> = { request: TRequest };
