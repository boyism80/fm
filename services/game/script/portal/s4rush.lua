local GROUP_NAME = "warrior_rush"
local QUEST = 6110
local PARTY_SIZE = 2

local function warrior(class_id)
    return math.floor(class_id / 100) == 1 or class_id == 900
end

return {
    on_enter = function(me)
        if not me:quest(QUEST):started() then
            me:message("알 수 없는 힘으로 봉인되어 지나갈 수 없습니다.")
            return
        end
        local party = me:party()
        if party == nil then
            me:message("파티를 구성한 후 도전할 수 있습니다.")
            return
        end
        if party:leader_id() ~= me:id() then
            me:message("파티장이 입장을 시작할 수 있습니다.")
            return
        end
        local characters = me:map():characters()
        local here = 0
        for _, member in ipairs(party:members()) do
            if not warrior(member:class_id()) then
                me:message("파티원중 전사가 아닌 직업이 있어 도전할 수 없습니다.")
                return
            end
            if characters[member:id()] ~= nil then
                here = here + 1
            end
        end
        if here < PARTY_SIZE then
            me:message("파티 인원이 2명이 아니기 때문에 도전할 수 없습니다.")
            return
        end
        local group = state_machine(GROUP_NAME)
        if group == nil then
            me:message("지금은 시험장에 들어갈 수 없습니다.")
            return
        end
        local sm, err = group:start_party(me, party)
        if sm == nil then
            me:message("다른 파티가 시험 중입니다. 잠시 후 다시 시도하세요.")
            if err ~= nil then
                log("warrior_rush start_party:", err)
            end
        end
    end,
}
