export function toPgInt(value: unknown): number {
    if (value == null || value === "") {
        return 0;
    }
    if (typeof value === "number" && Number.isFinite(value)) {
        return Math.trunc(value);
    }
    const n = Number(value);
    if (!Number.isFinite(n)) {
        return 0;
    }
    return Math.trunc(n);
}

export function toPgIntOrNull(value: unknown): number | null {
    if (value == null || value === "") {
        return null;
    }
    if (typeof value === "number" && Number.isFinite(value)) {
        return Math.trunc(value);
    }
    const n = Number(value);
    if (!Number.isFinite(n)) {
        return null;
    }
    return Math.trunc(n);
}
