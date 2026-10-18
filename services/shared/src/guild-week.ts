const WEEK_MS = 7 * 24 * 60 * 60 * 1000;
// 2026-01-01 00:00 KST, a Thursday: guild ranking weeks start every Thursday at midnight KST.
const FIRST_WEEK_START_MS = Date.UTC(2025, 11, 31, 15);

export function guildWeek(unixMs: number): number {
    return Math.floor((unixMs - FIRST_WEEK_START_MS) / WEEK_MS);
}

export function guildWeekStart(week: number): number {
    return FIRST_WEEK_START_MS + week * WEEK_MS;
}
