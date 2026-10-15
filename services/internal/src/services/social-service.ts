import { randomUUID } from "node:crypto";
import type { InternalContext } from "../context/internal-context";
import type { CharacterRepository } from "../repos/character-repository";
import type { NoteRemoveRequestRepository } from "../repos/note-request-repository";
import type { UnifiedRepository } from "../repos/unified-repository";
import type { SessionRepository } from "../repos/session-repository";
import type { DistributedLockService } from "./distributed-lock-service";
import type { RabbitMQService } from "./rabbitmq-service";
import type { CharacterModel } from "../types/repository-models";
import type {
    CompleteWhisperRequest,
    MessengerReply,
    MessengerRequest,
    Note,
    NotesReply,
    NotesRequest,
    WhisperReply,
    WhisperRequest,
} from "../protobuf/generated/fminternal/internal_service";
import {
    MessengerAction,
    NoteResult,
    NotesAction,
    WhisperMode,
} from "../protobuf/generated/fminternal/internal_service";
import { AccountSessionState } from "../session-state";
import { redisSessionKey } from "../redis-session-key";

type Member = {
    characterId: number;
    name: string;
    channelId: number;
    position: number;
    look: string;
    stamp: string;
};

type Invite = {
    sender: number;
    receiver: number;
    expires: number;
};

type Room = {
    id: number;
    members: Member[];
    invites: Invite[];
};

const AMQ_DIRECT_EXCHANGE = "amq.direct";
const MESSENGER_LOCK = "social:messengers";
const MESSENGER_CAPACITY = 3;
const MESSENGER_INVITE_LIMIT = 20;
const MESSENGER_INVITE_TTL_MS = 60000;
const MESSENGER_LOOK_LIMIT = 2048;
const NAME_LIMIT = 13;
const TEXT_LIMIT = 200;
const NOTE_LIMIT = 100;

enum NoteKind {
    NORMAL = 0,
    POPULATION = 1,
}

function emptyMessenger(): MessengerReply {
    return { ok: false, roomId: 0, position: 0, members: [] };
}

function offline(): WhisperReply {
    return { found: false, characterId: 0, name: "", channelId: -1, cashShop: false, mapId: 0 };
}

function notesReply(status: NoteResult): NotesReply {
    return { status, notes: [], population: 0 };
}

function memberReply(member: Member) {
    return { ...member, look: Buffer.from(member.look, "base64") };
}

export class SocialService {
    constructor(
        private readonly internalContext: InternalContext,
        private readonly characterRepository: CharacterRepository,
        private readonly noteRemoveRequestRepository: NoteRemoveRequestRepository,
        private readonly unifiedRepository: UnifiedRepository,
        private readonly sessionRepository: SessionRepository,
        private readonly distributedLockService: DistributedLockService,
        private readonly rabbitmqService: RabbitMQService,
    ) {}

    private key(world: number, rest: string) {
        return redisSessionKey(`w${world}:social:${rest}`);
    }

    private redis(world: number) {
        return this.internalContext.getRedisGlobalAccess(world).client;
    }

    private validText(text: string, max: number) {
        return text.trim().length > 0 && text.length <= max && !/[\x00-\x1f\x7f]/u.test(text);
    }

    private async named(world: number, name: string) {
        if (!this.validText(name, NAME_LIMIT)) {
            return null;
        }

        const entry = await this.unifiedRepository.findCharacterNameEntry(name);
        if (entry?.world_id !== world) {
            return null;
        }
        return this.characterRepository.get(world, entry.character_id);
    }

    private async session(world: number, character: CharacterModel) {
        const account = await this.sessionRepository.getAccountSession(world, character.accountId);
        if (account?.worldId !== world || account.characterId !== character.characterId) {
            return null;
        }
        return account;
    }

    async publish(world: number, channel: number, eventType: string, payload: Record<string, unknown>) {
        await this.rabbitmqService.assertDirectExchange(AMQ_DIRECT_EXCHANGE);
        await this.rabbitmqService.publish(AMQ_DIRECT_EXCHANGE, `fm.${world}.${channel}.parcel`, eventType, {
            occurred_at: new Date().toISOString(),
            ...payload,
        });
    }

    // The target's channel claims the request through completeWhisper and
    // pushes the whisper result to the reply list this call waits on.
    async whisper(req: WhisperRequest): Promise<WhisperReply> {
        const modes = [WhisperMode.WHISPER_MODE_FIND, WhisperMode.WHISPER_MODE_SEND, WhisperMode.WHISPER_MODE_FIND_BUDDY];
        if (!modes.includes(req.mode) || !this.validText(req.targetName, NAME_LIMIT)) {
            return offline();
        }
        if (req.mode === WhisperMode.WHISPER_MODE_SEND && !this.validText(req.message, TEXT_LIMIT)) {
            return offline();
        }

        const sender = await this.characterRepository.get(req.worldId, req.characterId);
        const target = await this.named(req.worldId, req.targetName);
        if (!sender || !target) {
            return offline();
        }
        if (target.hidden && (sender.role ?? 0) < (target.role ?? 0)) {
            return offline();
        }
        if ((await this.sessionRepository.findChannel(req.worldId, sender.characterId)) !== req.channelId) {
            return offline();
        }

        const session = await this.session(req.worldId, target);
        if (req.mode !== WhisperMode.WHISPER_MODE_SEND && session?.state === AccountSessionState.ACCOUNT_SESSION_STATE_CASH_SHOP) {
            return { ...offline(), found: true, characterId: target.characterId, name: target.name, cashShop: true };
        }

        const channel = await this.sessionRepository.findChannel(req.worldId, target.characterId);
        if (channel == null) {
            return offline();
        }

        const redis = this.redis(req.worldId);
        const requestId = randomUUID();
        const pendingKey = this.key(req.worldId, `whisper:${requestId}`);
        const replyKey = `${pendingKey}:reply`;
        await redis.set(pendingKey, String(target.characterId), "EX", 8);
        const waiter = redis.duplicate();
        try {
            await this.publish(req.worldId, channel, "whisper", {
                character_id: target.characterId,
                sender_id: sender.characterId,
                sender_name: sender.name,
                sender_role: sender.role ?? 0,
                channel_id: req.channelId,
                mode: req.mode,
                message: req.message,
                request_id: requestId,
            });
            const result = await waiter.blpop(replyKey, 6);
            if (!result) {
                return offline();
            }

            const outcome = JSON.parse(result[1]) as { mapId: number; received: boolean };
            if (!outcome.received) {
                return offline();
            }
            return {
                found: true,
                characterId: target.characterId,
                name: target.name,
                channelId: channel,
                cashShop: false,
                mapId: outcome.mapId,
            };
        } finally {
            waiter.disconnect();
            await redis.del(pendingKey, replyKey);
        }
    }

    async completeWhisper(req: CompleteWhisperRequest) {
        const redis = this.redis(req.worldId);
        const key = this.key(req.worldId, `whisper:${req.requestId}`);
        const claimed = `claimed:${req.characterId}`;
        const expected = req.complete ? claimed : String(req.characterId);
        const ok = await redis.eval(
            `
            if redis.call('GET', KEYS[1]) ~= ARGV[1] then return 0 end
            if ARGV[2] == '0' then
                redis.call('SET', KEYS[1], ARGV[3], 'KEEPTTL')
            else
                redis.call('DEL', KEYS[1])
                redis.call('LPUSH', KEYS[2], ARGV[4])
                redis.call('EXPIRE', KEYS[2], 8)
            end
            return 1`,
            2,
            key,
            `${key}:reply`,
            expected,
            req.complete ? "1" : "0",
            claimed,
            JSON.stringify({ mapId: req.mapId, received: req.received }),
        );
        return { ok: Number(ok) === 1 };
    }

    private async rooms(world: number) {
        const raw = await this.redis(world).hgetall(this.key(world, "messengers"));
        return Object.values(raw).map(value => JSON.parse(value) as Room);
    }

    private async store(world: number, room: Room) {
        const redis = this.redis(world);
        const key = this.key(world, "messengers");
        if (room.members.length === 0) {
            await redis.hdel(key, String(room.id));
        } else {
            await redis.hset(key, String(room.id), JSON.stringify(room));
        }
    }

    private async broadcast(world: number, room: Room, payload: Record<string, unknown>, except = 0) {
        for (const member of room.members) {
            if (member.characterId === except) {
                continue;
            }
            await this.publish(world, member.channelId, "messenger", {
                room_id: room.id,
                character_id: member.characterId,
                ...payload,
            });
        }
    }

    private async remove(world: number, room: Room, id: number) {
        const member = room.members.find(m => m.characterId === id);
        if (!member) {
            return;
        }

        room.members = room.members.filter(m => m.characterId !== id);
        room.invites = room.invites.filter(i => i.sender !== id && i.receiver !== id);
        await this.store(world, room);
        await this.broadcast(world, room, { action: MessengerAction.MESSENGER_ACTION_LEAVE, position: member.position });
    }

    // A member whose game session changed (relog, channel change) left the
    // room without a leave request.
    private async prune(world: number, rooms: Room[]) {
        for (const room of rooms) {
            for (const member of [...room.members]) {
                const character = await this.characterRepository.get(world, member.characterId);
                const active = character ? await this.session(world, character) : null;
                if (
                    active?.state !== AccountSessionState.ACCOUNT_SESSION_STATE_GAME ||
                    active.channelId !== member.channelId ||
                    active.timestamps.stateChangedAt !== member.stamp
                ) {
                    await this.remove(world, room, member.characterId);
                }
            }
            room.invites = room.invites.filter(i => i.expires > Date.now());
        }
    }

    async leaveMessenger(world: number, characterId: number) {
        await using lock = await this.distributedLockService.acquireWorldDataLock(world, MESSENGER_LOCK);
        for (const room of await this.rooms(world)) {
            await this.remove(world, room, characterId);
        }
    }

    async messenger(req: MessengerRequest): Promise<MessengerReply> {
        const world = req.worldId;
        const id = req.characterId;
        const character = await this.characterRepository.get(world, id);
        if (!character || (await this.sessionRepository.findChannel(world, id)) !== req.channelId) {
            return emptyMessenger();
        }

        const session = await this.session(world, character);
        if (!session) {
            return emptyMessenger();
        }

        await using lock = await this.distributedLockService.acquireWorldDataLock(world, MESSENGER_LOCK);
        const rooms = await this.rooms(world);
        await this.prune(world, rooms);

        const room = rooms.find(r => r.members.some(m => m.characterId === id));
        if (req.action === MessengerAction.MESSENGER_ACTION_ENTER) {
            return this.join(world, character, req, rooms, room, session.timestamps.stateChangedAt ?? "");
        }

        if (req.action === MessengerAction.MESSENGER_ACTION_DECLINE) {
            return this.decline(world, character, req.text, rooms);
        }

        if (!room) {
            return emptyMessenger();
        }

        switch (req.action) {
            case MessengerAction.MESSENGER_ACTION_LEAVE:
                await this.remove(world, room, id);
                break;
            case MessengerAction.MESSENGER_ACTION_INVITE:
                if (!(await this.invite(world, character, req.text, rooms, room))) {
                    return emptyMessenger();
                }
                break;
            case MessengerAction.MESSENGER_ACTION_CHAT: {
                if (!this.validText(req.text, TEXT_LIMIT)) {
                    return emptyMessenger();
                }

                // The client sends its own "name0"/"name1" join and leave lines;
                // every other line gets the sender's name so it cannot be forged.
                const prefix = `${character.name} : `;
                const own = req.text === `${character.name}0` || req.text === `${character.name}1`;
                const text = own || req.text.startsWith(prefix) ? req.text : prefix + req.text;
                await this.broadcast(world, room, { action: MessengerAction.MESSENGER_ACTION_CHAT, text }, id);
                break;
            }
            case MessengerAction.MESSENGER_ACTION_AVATAR: {
                if (!req.look.length || req.look.length > MESSENGER_LOOK_LIMIT) {
                    return emptyMessenger();
                }

                const member = room.members.find(m => m.characterId === id)!;
                member.look = Buffer.from(req.look).toString("base64");
                await this.store(world, room);
                await this.broadcast(world, room, { action: MessengerAction.MESSENGER_ACTION_AVATAR, member }, id);
                break;
            }
            default:
                return emptyMessenger();
        }
        return { ...emptyMessenger(), ok: true, roomId: room.id };
    }

    // Joining a room the character is already in replays the join reply,
    // which covers a reply lost on the way back to the game server.
    private async join(
        world: number,
        character: CharacterModel,
        req: MessengerRequest,
        rooms: Room[],
        current: Room | undefined,
        stamp: string,
    ): Promise<MessengerReply> {
        const id = character.characterId;
        if (current) {
            if (req.roomId !== 0 && req.roomId !== current.id) {
                return emptyMessenger();
            }

            const member = current.members.find(m => m.characterId === id)!;
            return {
                ok: true,
                roomId: current.id,
                position: member.position,
                members: current.members.filter(m => m !== member).map(memberReply),
            };
        }

        if (!req.look.length || req.look.length > MESSENGER_LOOK_LIMIT) {
            return emptyMessenger();
        }

        let room: Room | undefined;
        if (req.roomId === 0) {
            const roomId = await this.redis(world).incr(this.key(world, "messenger-sequence"));
            if (roomId > 0xffffffff) {
                throw new Error("Messenger ID exhausted");
            }
            room = { id: roomId, members: [], invites: [] };
        } else {
            room = rooms.find(r => r.id === req.roomId && r.members.length > 0);
            if (!room || room.members.length >= MESSENGER_CAPACITY || !room.invites.some(i => i.receiver === id)) {
                return emptyMessenger();
            }
        }

        const joined = room;
        const position = [0, 1, 2].find(p => !joined.members.some(m => m.position === p))!;
        const member: Member = {
            characterId: id,
            name: character.name,
            channelId: req.channelId,
            position,
            look: Buffer.from(req.look).toString("base64"),
            stamp,
        };
        const others = [...room.members];
        room.members.push(member);
        room.invites = room.invites.filter(i => i.receiver !== id);
        await this.store(world, room);
        await this.broadcast(world, room, { action: MessengerAction.MESSENGER_ACTION_ENTER, member }, id);
        return { ok: true, roomId: room.id, position, members: others.map(memberReply) };
    }

    private async decline(world: number, character: CharacterModel, inviterName: string, rooms: Room[]): Promise<MessengerReply> {
        const id = character.characterId;
        const inviter = await this.named(world, inviterName);
        if (!inviter) {
            return emptyMessenger();
        }

        const room = rooms.find(r => r.invites.some(i => i.receiver === id && i.sender === inviter.characterId));
        if (!room) {
            return emptyMessenger();
        }

        room.invites = room.invites.filter(i => i.receiver !== id);
        await this.store(world, room);
        const sender = room.members.find(m => m.characterId === inviter.characterId);
        if (sender) {
            await this.publish(world, sender.channelId, "messenger", {
                character_id: sender.characterId,
                room_id: room.id,
                action: MessengerAction.MESSENGER_ACTION_DECLINE,
                name: character.name,
            });
        }
        return { ...emptyMessenger(), ok: true };
    }

    private async invite(world: number, character: CharacterModel, targetName: string, rooms: Room[], room: Room) {
        const id = character.characterId;
        const target = await this.named(world, targetName);
        if (!target || target.characterId === id) {
            return false;
        }
        if (target.hidden && (character.role ?? 0) < (target.role ?? 0)) {
            return false;
        }

        const channel = await this.sessionRepository.findChannel(world, target.characterId);
        if (channel == null || room.members.length >= MESSENGER_CAPACITY) {
            return false;
        }
        if (rooms.some(r => r.members.some(m => m.characterId === target.characterId))) {
            return false;
        }

        room.invites = room.invites.filter(i => i.receiver !== target.characterId);
        if (room.invites.length >= MESSENGER_INVITE_LIMIT) {
            return false;
        }

        room.invites.push({ sender: id, receiver: target.characterId, expires: Date.now() + MESSENGER_INVITE_TTL_MS });
        await this.store(world, room);
        await this.publish(world, channel, "messenger", {
            character_id: target.characterId,
            action: MessengerAction.MESSENGER_ACTION_INVITE,
            room_id: room.id,
            name: character.name,
        });
        return true;
    }

    // A cash gift serial permits exactly one thank-you note to its actual sender.
    async gift(world: number, characterId: number, senderName: string, serial: string) {
        const sender = await this.named(world, senderName);
        if (!sender) {
            return;
        }

        await this.internalContext.withPgGlobalTransaction(world, tx =>
            tx.query(
                "INSERT INTO social_thanks (serial, character_id, sender_id) VALUES ($1, $2, $3) ON CONFLICT DO NOTHING",
                [serial, characterId, sender.characterId],
            ),
        );
    }

    async notes(req: NotesRequest): Promise<NotesReply> {
        const sender = await this.characterRepository.get(req.worldId, req.characterId);
        if (!sender) {
            return notesReply(NoteResult.NOTE_RESULT_WRONG_NAME);
        }

        switch (req.action) {
            case NotesAction.NOTES_ACTION_SEND:
                return this.sendNote(sender, req);
            case NotesAction.NOTES_ACTION_REMOVE:
            case NotesAction.NOTES_ACTION_REMOVE_GIFTS:
                return this.removeNotes(req);
            case NotesAction.NOTES_ACTION_LIST:
                return { ...notesReply(NoteResult.NOTE_RESULT_OK), notes: await this.listNotes(req.worldId, req.characterId) };
        }
        return notesReply(NoteResult.NOTE_RESULT_WRONG_NAME);
    }

    private async sendNote(sender: CharacterModel, req: NotesRequest): Promise<NotesReply> {
        const world = req.worldId;
        const id = sender.characterId;
        if ((req.gift !== NoteKind.NORMAL && req.gift !== NoteKind.POPULATION) || !this.validText(req.message, TEXT_LIMIT)) {
            return notesReply(NoteResult.NOTE_RESULT_WRONG_NAME);
        }

        const target = await this.named(world, req.targetName);
        if (!target || target.characterId === id) {
            return notesReply(NoteResult.NOTE_RESULT_WRONG_NAME);
        }

        const ok = await this.internalContext.withPgGlobalTransaction(world, async tx => {
            await tx.query("SELECT pg_advisory_xact_lock(280, $1)", [target.characterId]);
            const count = await tx.query<{ count: string }>(
                "SELECT COUNT(*) FROM social_notes WHERE receiver_id=$1",
                [target.characterId],
            );
            if (Number(count.rows[0]!.count) >= NOTE_LIMIT) {
                return false;
            }

            if (req.gift === NoteKind.POPULATION) {
                const token = await tx.query(
                    `UPDATE social_thanks SET consumed=TRUE WHERE serial=(
                        SELECT serial FROM social_thanks
                        WHERE character_id=$1 AND sender_id=$2 AND NOT consumed
                        ORDER BY serial LIMIT 1 FOR UPDATE
                    ) RETURNING serial`,
                    [id, target.characterId],
                );
                if (!token.rowCount) {
                    return false;
                }
            }

            await tx.query(
                "INSERT INTO social_notes (receiver_id, sender_id, sender_name, message, gift) VALUES ($1, $2, $3, $4, $5)",
                [target.characterId, id, sender.name, req.message, req.gift],
            );
            return true;
        });
        if (!ok) {
            return notesReply(NoteResult.NOTE_RESULT_INBOX_FULL);
        }

        const channel = await this.sessionRepository.findChannel(world, target.characterId);
        if (channel != null) {
            await this.publish(world, channel, "note_arrived", { character_id: target.characterId }).catch(() => {});
        }
        return notesReply(NoteResult.NOTE_RESULT_OK);
    }

    // The reward counts only thank-you notes this call actually deleted; the
    // client's gift flag never creates a reward on its own.
    private async removeNotes(req: NotesRequest): Promise<NotesReply> {
        const id = req.characterId;
        if (req.noteIds.length > NOTE_LIMIT) {
            return notesReply(NoteResult.NOTE_RESULT_WRONG_NAME);
        }

        const population = await this.internalContext.withPgGlobalTransaction(req.worldId, async tx => {
            await tx.query("SELECT pg_advisory_xact_lock(280, $1)", [id]);
            if (req.requestId !== "0") {
                const recorded = await this.noteRemoveRequestRepository.find(req.worldId, id, req.requestId, { txClient: tx });
                if (recorded?.status === "applied") {
                    return recorded.outcome.population;
                }
            }

            const removed = req.action === NotesAction.NOTES_ACTION_REMOVE_GIFTS
                ? await tx.query<{ gift: number }>(
                    "DELETE FROM social_notes WHERE receiver_id=$1 AND gift=$2 RETURNING gift",
                    [id, NoteKind.POPULATION],
                )
                : await tx.query<{ gift: number }>(
                    "DELETE FROM social_notes WHERE receiver_id=$1 AND id=ANY($2::int[]) RETURNING gift",
                    [id, req.noteIds],
                );
            const count = removed.rows.filter(n => n.gift === NoteKind.POPULATION).length;
            if (req.requestId !== "0") {
                await this.noteRemoveRequestRepository.insert(req.worldId, id, req.requestId, { population: count }, tx);
            }
            return count;
        });
        return {
            ...notesReply(NoteResult.NOTE_RESULT_OK),
            notes: await this.listNotes(req.worldId, id),
            population,
        };
    }

    async listNotes(world: number, characterId: number): Promise<Note[]> {
        const result = await this.internalContext.withPgGlobalTransaction(world, tx =>
            tx.query<{ id: number; sender_name: string; message: string; created_at: Date; gift: number }>(
                "SELECT id, sender_name, message, created_at, gift FROM social_notes WHERE receiver_id=$1 ORDER BY id LIMIT 100",
                [characterId],
            ),
        );
        return result.rows.map(n => ({
            id: n.id,
            senderName: n.sender_name,
            message: n.message,
            sentAtUnixMs: String(n.created_at.getTime()),
            gift: n.gift,
        }));
    }
}
