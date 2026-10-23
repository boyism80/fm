-- State machine (old/scripts/quest/2238.js): Sleepywood hidden dungeon

local solo_instance = require("script/lib/solo_instance")

return solo_instance.create({
    maps = { 910520000 },
    exit_map = 105100100,
    duration_ms = 600000,
    setup = function(sm, maps)
        maps[1]:respawn({ include_one_time = true })
    end,
})
