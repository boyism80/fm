import type { Pool } from "pg";
import { RequestRepository } from "./request-repository";

export class NoteRemoveRequestRepository extends RequestRepository<{ population: number }> {
    protected readonly kind = "note_remove";

    protected override pool(worldId: number): Pool {
        return this.ctx.getPgGlobalPool(worldId);
    }
}
