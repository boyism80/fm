import { LoginAccountReply_Status } from "../../protobuf/generated/fminternal/internal_service";
import type { AccountService } from "../../services/account-service";
import type { SessionService } from "../../services/session-service";
import type { LoginAccountReply, LoginAccountRequest } from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCallback, GrpcErrorHandler } from "./types";

type LoginAccountCall = { request: LoginAccountRequest };

type LoginResult = {
    status: number;
    accountId?: number;
    gender?: number;
    isChatBlocked?: boolean;
    chatBlockedUntil?: string | number | Date | null;
    characterSlotCount?: number;
    role?: number;
};
type InternalConfig = { app: { world_id: number }; grpc: { host: string; port: number } };

export function createAuthHandlers(
    accountService: AccountService,
    sessionService: SessionService,
    internalConfig: InternalConfig,
    grpcError: GrpcErrorHandler
) {
    return {
        async loginAccount(call: LoginAccountCall, callback: GrpcCallback<LoginAccountReply>) {
            try {
                const result = await accountService.loginAccount(
                    call.request.loginId,
                    call.request.password,
                    call.request.macAddress,
                    call.request.ipAddress,
                    call.request.initialRole
                );
                if (result.status === LoginAccountReply_Status.SUCCESS || result.status === LoginAccountReply_Status.REGISTERED) {
                    const accountId = result.accountId ?? 0;
                    const begin = await sessionService.beginLogin(
                        internalConfig.app.world_id,
                        accountId,
                        `${internalConfig.grpc.host}:${internalConfig.grpc.port}`
                    );
                    if (!begin.ok) {
                        callback(null, {
                            status: LoginAccountReply_Status.ALREADY_LOGGED_IN,
                            accountId: 0,
                            gender: 0,
                            isChatBlocked: false,
                            chatBlockedUntilUnixMs: 0,
                            characterSlotCount: 6,
                            role: result.role ?? 0,
                        });
                        return;
                    }
                }
                const blockedUntil = result.chatBlockedUntil ? new Date(result.chatBlockedUntil).getTime() : 0;
                callback(null, {
                    status: result.status,
                    accountId: result.accountId ?? 0,
                    gender: result.gender ?? 0,
                    isChatBlocked: result.isChatBlocked ?? false,
                    chatBlockedUntilUnixMs: blockedUntil,
                    characterSlotCount: result.characterSlotCount ?? 6,
                    role: result.role ?? 0,
                });
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}
