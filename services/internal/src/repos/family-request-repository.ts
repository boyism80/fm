import type { Pool } from "pg";
import { RequestRepository } from "./request-repository";

export class FamilyMoveRequestRepository extends RequestRepository<{ payer_id: number; feature: number }> {
    protected readonly kind = "family_move";

    protected override pool(worldId: number): Pool {
        return this.ctx.getPgGlobalPool(worldId);
    }
}
