import { LoginAccountReply_Status, SecondPasswordStatus } from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { AccountService } from "../../services/account-service";
import type { SessionService } from "../../services/session-service";
import type {
    LoginAccountReply,
    LoginAccountRequest,
    SetAccountRoleReply,
    SetAccountRoleRequest,
    SecondPasswordReply,
    SetSecondPasswordRequest,
    RemoveSecondPasswordRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type { GrpcCall, GrpcCallback, GrpcErrorHandler } from "./types";
import { Controller, Method } from "../grpc-method-decorator";

type LoginAccountCall = { request: LoginAccountRequest };
type InternalConfig = { grpc: { host: string; port: number } };

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
                role: result.role ?? 0,
            });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("setAccountRole")
    async setAccountRole(call: GrpcCall<SetAccountRoleRequest>, callback: GrpcCallback<SetAccountRoleReply>) {
        try {
            const ok = await this.accountService.setAccountRole(call.request.accountId >>> 0, call.request.role >>> 0);
            callback(null, { ok });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("setSecondPassword")
    async setSecondPassword(call: GrpcCall<SetSecondPasswordRequest>, callback: GrpcCallback<SecondPasswordReply>) {
        try {
            const result = await this.accountService.setSecondPassword(call.request.accountId >>> 0, call.request.password);
            let status = SecondPasswordStatus.SECOND_PASSWORD_INVALID;
            switch (result) {
                case "success":
                    status = SecondPasswordStatus.SECOND_PASSWORD_SUCCESS;
                    break;
                case "already_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_ALREADY_SET;
                    break;
                case "not_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_NOT_SET;
                    break;
                case "locked":
                    status = SecondPasswordStatus.SECOND_PASSWORD_LOCKED;
                    break;
            }
            callback(null, { status });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }

    @Method("removeSecondPassword")
    async removeSecondPassword(call: GrpcCall<RemoveSecondPasswordRequest>, callback: GrpcCallback<SecondPasswordReply>) {
        try {
            const result = await this.accountService.removeSecondPassword(call.request.accountId >>> 0, call.request.password);
            let status = SecondPasswordStatus.SECOND_PASSWORD_INVALID;
            switch (result) {
                case "success":
                    status = SecondPasswordStatus.SECOND_PASSWORD_SUCCESS;
                    break;
                case "already_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_ALREADY_SET;
                    break;
                case "not_set":
                    status = SecondPasswordStatus.SECOND_PASSWORD_NOT_SET;
                    break;
                case "locked":
                    status = SecondPasswordStatus.SECOND_PASSWORD_LOCKED;
                    break;
            }
            callback(null, { status });
        } catch (err) {
            this.grpcError(err, callback);
        }
    }
}
