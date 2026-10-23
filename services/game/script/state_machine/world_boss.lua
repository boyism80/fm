-- State machine (old/scripts/event/WorldBossEvent.js): World boss

local boss_instance = require("script/lib/boss_instance")

return boss_instance.create({
    map = 123456785,
    exit_map = 123456789,
    duration_ms = 3600000,
    bosses = {
        { id = 5000007, x = 156, y = 258 },
        { id = 5000007, x = 559, y = 545 },
    },
})
