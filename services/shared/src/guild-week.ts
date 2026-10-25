const WEEK_MS = 7 * 24 * 60 * 60 * 1000;
// 2026-01-01 00:00 in the server's local time zone, a Thursday: guild ranking weeks start every Thursday at local midnight.
const FIRST_WEEK_START_MS = new Date(2026, 0, 1).getTime();

export function guildWeek(unixMs: number): number {
    return Math.floor((unixMs - FIRST_WEEK_START_MS) / WEEK_MS);
}

export function guildWeekStart(week: number): number {
    return FIRST_WEEK_START_MS + week * WEEK_MS;
}
