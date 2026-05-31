import { LoginAccountReply_Status } from "../../protobuf/generated/fminternal/internal_service";
import type { AccountService } from "../../services/account-service";
import type { SessionService } from "../../services/session-service";
import type { LoginAccountReply, LoginAccountRequest } from "../../protobuf/generated/fminternal/internal_service";
import type { GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

type LoginAccountCall = { request: LoginAccountRequest };
type InternalConfig = { app: { world_id: number }; grpc: { host: string; port: number } };

@Controller("authController")
export class AuthGrpcController {
    private readonly accountService: AccountService;
    private readonly sessionService: SessionService;
    private readonly internalConfig: InternalConfig;
    private readonly grpcError: GrpcErrorHandler;

    constructor(
        accountService: AccountService,
        sessionService: SessionService,
        internalConfig: InternalConfig,
        grpcError: GrpcErrorHandler
    ) {
        this.accountService = accountService;
        this.sessionService = sessionService;
        this.internalConfig = internalConfig;
        this.grpcError = grpcError;
    }

    @Method("loginAccount")
    async loginAccount(call: LoginAccountCall, callback: GrpcCallback<LoginAccountReply>) {
        try {
            const result = await this.accountService.loginAccount(
                call.request.loginId,
                call.request.password,
                call.request.macAddress,
                call.request.ipAddress,
                call.request.initialRole
            );
            if (result.status === LoginAccountReply_Status.SUCCESS || result.status === LoginAccountReply_Status.REGISTERED) {
                const accountId = result.accountId ?? 0;
                const begin = await this.sessionService.beginLogin(
                    this.internalConfig.app.world_id,
                    accountId,
                    `${this.internalConfig.grpc.host}:${this.internalConfig.grpc.port}`
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
            this.grpcError(err, callback);
        }
    }
}
