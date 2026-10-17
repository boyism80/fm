import type { PlayerNpcRepository, PlayerNpcRow } from "@fm/shared/repos/player-npc-repository";
import type { RabbitMQService } from "./rabbitmq-service";

const PLAYER_NPC_CREATED_EVENT = "player_npc_created";
const PLAYER_NPC_REMOVED_EVENT = "player_npc_removed";

export class PlayerNpcService {
    private readonly playerNpcRepository: PlayerNpcRepository;
    private readonly rabbitmqService: RabbitMQService;

    constructor(playerNpcRepository: PlayerNpcRepository, rabbitmqService: RabbitMQService) {
        this.playerNpcRepository = playerNpcRepository;
        this.rabbitmqService = rabbitmqService;
    }

    async list(worldId: number): Promise<PlayerNpcRow[]> {
        return this.playerNpcRepository.list(worldId);
    }

    // A lost notice only delays the other channels until they restart, so it does not undo the write.
    async create(worldId: number, npc: PlayerNpcRow): Promise<boolean> {
        if (!(await this.playerNpcRepository.insert(worldId, npc))) {
            return false;
        }
        try {
            await this.rabbitmqService.publish(`fm.${worldId}.all.global`, PLAYER_NPC_CREATED_EVENT, npc);
        } catch (err) {
            console.error("[player-npc] created notice failed", err);
        }
        return true;
    }

    async remove(worldId: number, npcId: number): Promise<boolean> {
        if (!(await this.playerNpcRepository.remove(worldId, npcId))) {
            return false;
        }
        try {
            await this.rabbitmqService.publish(`fm.${worldId}.all.global`, PLAYER_NPC_REMOVED_EVENT, { npc_id: npcId });
        } catch (err) {
            console.error("[player-npc] removed notice failed", err);
        }
        return true;
    }
}
