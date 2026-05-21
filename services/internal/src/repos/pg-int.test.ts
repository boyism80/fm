import assert from "node:assert/strict";
import test from "node:test";
import { toPgInt, toPgIntOrNull } from "./pg-int";

test("toPgInt coerces pg bigint strings", () => {
    assert.equal(toPgInt("42"), 42);
    assert.equal(toPgInt(42.9), 42);
    assert.equal(toPgInt(null), 0);
});

test("toPgIntOrNull preserves null and coerces strings", () => {
    assert.equal(toPgIntOrNull(null), null);
    assert.equal(toPgIntOrNull(""), null);
    assert.equal(toPgIntOrNull("99"), 99);
});
