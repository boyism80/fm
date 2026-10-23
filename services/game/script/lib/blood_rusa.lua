local solo_instance = require("script/lib/solo_instance")

local function create(map, exit_map)
    return solo_instance.create({
        maps = { map },
        exit_map = exit_map,
        duration_ms = 190000,
        timeout_message = "제한시간이 다 되어 블러드 루사 본부로 강제 이동 됩니다.",
        setup = function(sm, maps)
            maps[1]:respawn({ include_one_time = true })
        end,
    })
end

return create
