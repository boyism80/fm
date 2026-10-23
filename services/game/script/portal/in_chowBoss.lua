local SCRIPT = "script/portal/in_chowBoss.lua"
local ROOM = 300010420

return {

    reset_map = function(map)
        map:reset()
    end,

    on_enter = function(me)
        local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
        if ok == false or count ~= 0 then
            me:message("이미 누군가 안에 입장해있습니다.", Msg.PinkText)
            return
        end
        run_on_map(ROOM, SCRIPT, "reset_map")
        me:map(ROOM)
    end
}
