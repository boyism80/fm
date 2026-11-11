-- Reactor name (Reactor.wz/2709000.img.xml): 핑크빈 파워업

local KIRSTON = 2141000
-- Beyond the Ruins (3522 for adventurers, 3538 for Cygnus) ends once Pink Bean wakes, which its record quest 7402 shows.
local BEYOND_THE_RUINS = { 3522, 3538 }
local AWAKENED = 7402

return {
    on_reactor = function(reactor)
        local map = reactor:map()
        if map == nil then
            return
        end
        map:spawn_mob(8820008, 8, -53)

        for _, ch in pairs(map:characters()) do
            for _, quest_id in ipairs(BEYOND_THE_RUINS) do
                local quest = ch:quest(quest_id)
                if quest ~= nil and quest:started() then
                    local awakened = ch:quest(AWAKENED)
                    if awakened:started() then
                        awakened:record("1")
                    else
                        awakened:start(KIRSTON, "1")
                    end
                    ch:show_quest_completion(quest_id)
                end
            end
        end
    end
}
