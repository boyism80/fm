local pq = require("script/lib/party_quest")

local SCRIPT = "script/portal/in_fairyBoss.lua"
local ROOM = 300030310

local function enter(me)
    local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
    if me:quest(31229):started() == false or ok == false or count ~= 0 then
        me:message("이미 누군가 안에 입장해있거나 퀘스트를 진행중이시지 않습니다.", Msg.PinkText)
        return
    end
    run_on_map(ROOM, SCRIPT, "reset_map")
    me:map(ROOM)
end

return {

    reset_map = function(map)
        map:reset()
    end,

    on_enter = function(me)
        local map = me:map()
        if not pq.hold_entry(map, SCRIPT) then
            me:message("이미 누군가 안에 입장해있거나 퀘스트를 진행중이시지 않습니다.", Msg.PinkText)
            return
        end
        enter(me)
        pq.release_entry(map, SCRIPT)
    end
}
