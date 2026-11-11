-- The bosses the boss waiting room's NPCs send a party to. A party fights the phases in order, each on its own map
-- with its own clock; a multi-phase boss names how its win is told to the world.
return {
    hilla = {
        name = "힐라",
        difficulty = "#k★#k☆☆☆☆#b#n",
        limit = 5,
        record = "boss_entry.hilla",
        phases = {
            { map = 123356780, time = 40, mobs = { 8870000 }, x = 164, y = 196 },
        },
    },
    magnus = {
        name = "매그너스",
        difficulty = "#k★★#k☆☆☆#b#n",
        limit = 5,
        record = "boss_entry.magnus",
        phases = {
            { map = 123356785, time = 30, mobs = { 8880000 }, x = 2365, y = -1347 },
        },
    },
    suu = {
        name = "스우",
        difficulty = "#k★#k☆☆☆☆#b#n",
        limit = 5,
        record = "boss_entry.suu",
        victory = "스우와의 전투에서 승리하였습니다.",
        phases = {
            { map = 123356782, time = 40, mobs = { 9801028 }, x = -110, y = -16 },
            { map = 123356783, time = 30, mobs = { 9801029 }, x = 4, y = -16 },
            { map = 123356784, time = 30, mobs = { 9801030 }, x = 35, y = 19 },
        },
    },
    chaos_papulatus = {
        name = "카오스 파풀라투스",
        difficulty = "#k★★★#k☆☆#b#n",
        limit = 10,
        record = "boss_entry.chaos_papulatus",
        phases = {
            { map = 123356781, time = 50, mobs = { 8820119 }, x = -393, y = -386 },
        },
    },
    damien = {
        name = "데미안",
        difficulty = "#k★★★★#k☆#b#n",
        limit = 3,
        record = "boss_entry.damien",
        victory = "데미안과의 전투에서 승리하였습니다.",
        phases = {
            { map = 123356786, time = 60, mobs = { 9300890 }, x = 1222, y = 16 },
            { map = 123356787, time = 60, mobs = { 9300891 }, x = 1222, y = 16 },
        },
    },
    black_mage = {
        name = "검은 마법사",
        difficulty = "#r★★★★★(헬)#b#n",
        limit = 3,
        record = "boss_entry.black_mage",
        victory = "검은 마법사와의 전투에서 승리하였습니다.",
        phases = {
            { map = 123356788, time = 120, mobs = { 8880500 }, x = -3, y = 85 },
            { map = 123356789, time = 60, mobs = { 8880501 }, x = -235, y = 215 },
        },
    },
    lucid = {
        name = "루시드",
        difficulty = "#k★★★★#k☆#b#n",
        limit = 4,
        record = "boss_entry.lucid",
        victory = "루시드와의 전투에서 승리하였습니다.",
        phases = {
            { map = 450004150, time = 60, mobs = { 8880166, 8880140 }, x = 1008, y = 48 },
            { map = 450004550, time = 60, mobs = { 8880150 }, x = 407, y = -125 },
        },
    },
}
