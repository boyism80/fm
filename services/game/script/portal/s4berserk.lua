local pq = require("script/lib/party_quest")

local GROUP_NAME = "warrior_berserk"
local QUEST = 6153
local KEY = 4031475
local SHIELD = 4031471

return {
    on_enter = function(me)
        if not me:quest(QUEST):started() then
            me:message("아직 이 시험에 들어갈 수 없습니다.")
            return
        end
        local class_id = me:class()
        if class_id ~= 112 and class_id ~= 122 and class_id ~= 132 and class_id ~= 900 then
            me:message("전사만 이 시험을 받을 수 있습니다.")
            return
        end
        if not pq.has_item(me, KEY) then
            me:message("시험장에 들어갈 열쇠가 없습니다.")
            return
        end
        if pq.has_item(me, SHIELD) then
            me:message("이미 방패를 얻었습니다.")
            return
        end
        local group = state_machine(GROUP_NAME)
        if group == nil then
            me:message("지금은 시험장에 들어갈 수 없습니다.")
            return
        end
        local sm, err = group:start_solo(me)
        if sm == nil then
            me:message("다른 사람이 시험 중입니다. 잠시 후 다시 시도하세요.")
            if err ~= nil then
                log("warrior_berserk start_solo:", err)
            end
        end
    end,
}
