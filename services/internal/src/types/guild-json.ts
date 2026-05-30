export type GuildLogo = {
    logo: number;
    logoColor: number;
    logoBG: number;
    logoBGColor: number;
};

export type GuildRankTitles = [string, string, string, string, string];

export const DEFAULT_GUILD_LOGO: GuildLogo = {
    logo: 0,
    logoColor: 0,
    logoBG: 0,
    logoBGColor: 0,
};

export const DEFAULT_GUILD_RANK_TITLES: GuildRankTitles = [
    "Master",
    "Jr. Master",
    "Member",
    "Member",
    "Member",
];
