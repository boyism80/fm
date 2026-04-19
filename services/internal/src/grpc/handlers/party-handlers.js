"use strict";

function createPartyHandlers(partyService, messages, grpcError) {
    function setOptionalPartyId(reply, partyId) {
        if (partyId == null) {
            reply.clearPartyId();
            return;
        }
        const n = Number(partyId);
        if (Number.isInteger(n) && n >= 0) {
            reply.setPartyId(n);
        } else {
            reply.clearPartyId();
        }
    }

    function setPartyMemberDoor(mm, doorRow) {
        const d = doorRow;
        if (
            d &&
            typeof d === "object" &&
            Number.isFinite(Number(d.town)) &&
            Number.isFinite(Number(d.target)) &&
            Number.isFinite(Number(d.x)) &&
            Number.isFinite(Number(d.y))
        ) {
            const pd = new messages.PartyDoor();
            pd.setTown(Number(d.town));
            pd.setTarget(Number(d.target));
            pd.setX(Number(d.x));
            pd.setY(Number(d.y));
            mm.setDoor(pd);
        } else {
            mm.clearDoor();
        }
    }

    return {
        async createParty(call, callback) {
            try {
                const result = await partyService.createParty(
                    call.request.getWorldId(),
                    call.request.getLeaderCharacterId()
                );
                const reply = new messages.CreatePartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async joinParty(call, callback) {
            try {
                const result = await partyService.joinParty(
                    call.request.getWorldId(),
                    call.request.getPartyId(),
                    call.request.getCharacterId(),
                    call.request.getCharacterName(),
                    call.request.getLevel(),
                    call.request.getClassId()
                );
                const reply = new messages.JoinPartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async inviteParty(call, callback) {
            try {
                const result = await partyService.inviteParty(
                    call.request.getWorldId(),
                    call.request.getInviterCharacterId(),
                    call.request.getTargetCharacterName()
                );
                const reply = new messages.InvitePartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                reply.setTargetCharacterId(result.targetCharacterId ?? 0);
                reply.setTargetChannelId(result.targetChannelId ?? 0);
                setOptionalPartyId(reply, result.partyId);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async leaveParty(call, callback) {
            try {
                const result = await partyService.leaveParty(
                    call.request.getWorldId(),
                    call.request.getCharacterId()
                );
                const reply = new messages.LeavePartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                reply.setDisbanded(Boolean(result.disbanded));
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async expelParty(call, callback) {
            try {
                const result = await partyService.expelParty(
                    call.request.getWorldId(),
                    call.request.getRequesterCharacterId(),
                    call.request.getTargetCharacterId()
                );
                const reply = new messages.ExpelPartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                reply.setDisbanded(Boolean(result.disbanded));
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async changePartyLeader(call, callback) {
            try {
                const result = await partyService.changePartyLeader(
                    call.request.getWorldId(),
                    call.request.getPartyId(),
                    call.request.getRequesterCharacterId(),
                    call.request.getNewLeaderCharacterId()
                );
                const reply = new messages.ChangePartyLeaderReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async getParty(call, callback) {
            try {
                const result = await partyService.getParty(
                    call.request.getWorldId(),
                    call.request.getPartyId()
                );
                const reply = new messages.GetPartyReply();
                reply.setFound(Boolean(result.found));
                if (result.found) {
                    const p = new messages.PartySnapshot();
                    p.setWorldId(result.party.worldId);
                    p.setPartyId(result.party.partyId);
                    p.setLeaderCharacterId(result.party.leaderCharacterId);
                    p.setRevision(result.party.revision);
                    p.setState(result.party.state);
                    p.setMembersList(
                        result.members.map((m) => {
                            const mm = new messages.PartyMemberSnapshot();
                            mm.setCharacterId(m.characterId);
                            mm.setCharacterName(m.characterName);
                            mm.setLevel(m.level);
                            mm.setClassId(m.classId);
                            mm.setRole(m.role);
                            mm.setMapId(Number(m.mapId ?? 0));
                            mm.setChannelIndex(Number(m.channelIndex ?? -2));
                            setPartyMemberDoor(mm, m.door);
                            return mm;
                        })
                    );
                    reply.setParty(p);
                }
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async reportPartyMemberSnapshot(call, callback) {
            try {
                const req = call.request;
                let doorPayload = null;
                if (typeof req.hasDoor === "function" && req.hasDoor()) {
                    const d = req.getDoor();
                    if (d) {
                        doorPayload = {
                            town: d.getTown(),
                            target: d.getTarget(),
                            x: d.getX(),
                            y: d.getY(),
                        };
                    }
                }
                const result = await partyService.reportPartyMemberSnapshot(
                    req.getWorldId(),
                    req.getCharacterId(),
                    req.getLevel(),
                    req.getClassId(),
                    req.getMapId(),
                    doorPayload
                );
                const reply = new messages.ReportPartyMemberSnapshotReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                setOptionalPartyId(reply, result.partyId);
                reply.setRevision(result.revision ?? 0);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },

        async denyParty(call, callback) {
            try {
                const result = await partyService.denyParty(
                    call.request.getWorldId(),
                    call.request.getDeniedCharacterId(),
                    call.request.getInviterName(),
                    call.request.getAction()
                );
                const reply = new messages.DenyPartyReply();
                reply.setOk(Boolean(result.ok));
                reply.setErrorCode(result.code ?? messages.PartyErrorCode.UNKNOWN);
                callback(null, reply);
            } catch (err) {
                grpcError(err, callback);
            }
        },
    };
}

module.exports = { createPartyHandlers };
