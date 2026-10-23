local ROOM = 240093300

return {

    on_enter = function(me)
        local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
        if me:quest(31351):started() == false or ok == false or count ~= 0 then
            me:message("길이 막혀있거나 이미 누군가 입장해 있습니다.", Msg.PinkText)
            return
        end
        me:open_npc(2210014)
    end
}
