local pq = require("script/lib/party_quest")

local GROUP_NAME = "warrior_berserk"
local QUEST = 6153
local KEY = 4031475
local SHIELD = 4031471

return {
    on_enter = function(me)
        if not me:quest(QUEST):started() then
            me:message("알 수 없는 힘으로 봉인되어 있습니다.")
            return
        end
        local class_id = me:class()
        if class_id ~= 112 and class_id ~= 122 and class_id ~= 132 and class_id ~= 900 then
            me:message("전사만 이 시험을 받을 수 있습니다.")
            return
        end
        if not pq.has_item(me, KEY) then
            me:message("잊혀진 신전으로 가는 열쇠가 필요합니다.")
            return
        end
        if pq.has_item(me, SHIELD) then
            me:message("이미 세이람의 방패를 갖고 있습니다.")
            return
        end
        local group = state_machine(GROUP_NAME)
        if group == nil then
            me:message("지금은 시험장에 들어갈 수 없습니다.")
            return
        end
        local sm, err = group:start_solo(me)
        if sm == nil then
            me:message("이미 이 안에 다른 누군가가 퀘스트에 도전중입니다.")
            if err ~= nil then
                log("warrior_berserk start_solo:", err)
            end
        end
    end,
}
