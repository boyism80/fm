import { CharacterRecordRepository } from "./character-record-repository";

export class AccountRecordRepository extends CharacterRecordRepository {
    protected override get table() {
        return "account_records";
    }

    protected override get ownerColumn() {
        return "account_id";
    }
}
