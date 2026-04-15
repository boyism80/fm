"use strict";

function createAuthHandlers(accountService, sessionService, internalConfig, messages, grpcError) {
    return {
        async loginAccount(call, callback) {
            try {
                const result = await accountService.loginAccount(
                    call.request.getLoginId(),
                    call.request.getPassword(),
                    call.request.getMacAddress(),
                    call.request.getIpAddress(),
                    call.request.getInitialRole()
                );
                const reply = new messages.LoginAccountReply();
                if (result.status === messages.LoginAccountReply.Status.SUCCESS ||
                    result.status === messages.LoginAccountReply.Status.REGISTERED) {
                    const begin = await sessionService.beginLogin(
                        internalConfig.app.world_id,
                        result.accountId,
                        `${internalConfig.grpc.host}:${internalConfig.grpc.port}`
                    );
                    if (!begin.ok) {
                        reply.setStatus(messages.LoginAccountReply.Status.ALREADY_LOGGED_IN);
                        reply.setRole(result.role ?? 0);
                        callback(null, reply);
                        return;
                    }
                }
                reply.setStatus(result.status);
                reply.setAccountId(result.accountId ?? 0);
                reply.setGender(result.gender ?? 0);
                reply.setIsChatBlocked(result.isChatBlocked ?? false);
                const blockedUntil = result.chatBlockedUntil
                    ? new Date(result.chatBlockedUntil).getTime()
                    : 0;
                reply.setChatBlockedUntilUnixMs(blockedUntil);
                reply.setCharacterSlotCount(result.characterSlotCount ?? 6);
                reply.setRole(result.role ?? 0);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createAuthHandlers };
