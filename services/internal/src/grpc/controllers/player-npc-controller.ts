import type {
    PlayerNpcRow,
} from "@fm/shared/repos/player-npc-repository";
import type {
    PlayerNpcService,
} from "../../services/player-npc-service";
import type {
    CreatePlayerNpcReply,
    CreatePlayerNpcRequest,
    ListPlayerNpcsReply,
    ListPlayerNpcsRequest,
    PlayerNpc,
    RemovePlayerNpcReply,
    RemovePlayerNpcRequest,
} from "@fm/shared/protobuf/generated/fminternal/internal_service";
import type {
    GrpcCall,
    GrpcCallback,
} from "./types";
import { Controller, Method } from "../grpc-method-decorator";

@Controller("playerNpcController")
export class PlayerNpcGrpcController {
    private readonly playerNpcService: PlayerNpcService;

    constructor(playerNpcService: PlayerNpcService) {
        this.playerNpcService = playerNpcService;
    }

    @Method("listPlayerNpcs")
    async listPlayerNpcs(call: GrpcCall<ListPlayerNpcsRequest>, callback: GrpcCallback<ListPlayerNpcsReply>) {
        const rows = await this.playerNpcService.list(call.request.worldId >>> 0);
        callback(null, {
            npcs: rows.map((row) => ({
                npcId: row.npc_id,
                characterId: row.character_id,
                name: row.name,
                gender: row.gender,
                skin: row.skin,
                face: row.face,
                hair: row.hair,
                equips: row.equips.map((equip) => ({ slot: equip.slot, itemId: equip.item_id })),
                mapId: row.map_id,
                x: row.x,
                y: row.y,
            })),
        });
    }

    @Method("createPlayerNpc")
    async createPlayerNpc(call: GrpcCall<CreatePlayerNpcRequest>, callback: GrpcCallback<CreatePlayerNpcReply>) {
        const npc: PlayerNpc | undefined = call.request.npc;
        if (npc === undefined || npc.npcId === 0) {
            callback(null, { created: false });
            return;
        }
        const row: PlayerNpcRow = {
            npc_id: npc.npcId >>> 0,
            character_id: npc.characterId >>> 0,
            name: npc.name,
            gender: npc.gender >>> 0,
            skin: npc.skin >>> 0,
            face: npc.face >>> 0,
            hair: npc.hair >>> 0,
            equips: npc.equips.map((equip) => ({ slot: equip.slot, item_id: equip.itemId >>> 0 })),
            map_id: npc.mapId >>> 0,
            x: npc.x,
            y: npc.y,
        };
        const created = await this.playerNpcService.create(call.request.worldId >>> 0, row);
        callback(null, { created });
    }

    @Method("removePlayerNpc")
    async removePlayerNpc(call: GrpcCall<RemovePlayerNpcRequest>, callback: GrpcCallback<RemovePlayerNpcReply>) {
        const removed = await this.playerNpcService.remove(call.request.worldId >>> 0, call.request.npcId >>> 0);
        callback(null, { removed });
    }
}
