local ROOM = 900000000

return {

    on_enter = function(me)
        if me:quest(2073):started() == false then
            me:message("알 수 없는 힘으로 가로막혀 있어 입장할 수 없다.", Msg.PinkText)
            return
        end
        local ok, count = run_on_map(ROOM, "script/lib/party_quest.lua", "character_count")
        if ok == false or count ~= 0 then
            me:message("이미 누군가가 이 안에서 퀘스트를 진행하는 중입니다.", Msg.PinkText)
            return
        end
        local function on_arrive(me)
            me:clock(600, function(me)
                me:map(100030000)
            end)
        end
        me:map(ROOM, 0, { callback = on_arrive })
    end
}
