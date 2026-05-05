import { Repository } from "./repository";

export class ValueRepository<TModel = Record<string, unknown>, TRow = Record<string, unknown>, TKey = unknown> extends Repository<TModel, TRow, TKey> {
    async getAll(): Promise<never> {
        throw new Error(`${this.constructor.name}.getAll is not supported for value repositories`);
    }

    async delAll(): Promise<never> {
        throw new Error(`${this.constructor.name}.delAll is not supported for value repositories`);
    }
}
