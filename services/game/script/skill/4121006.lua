-- Skill name (String.wz/Skill.img.xml): 스피릿 자벨린

local SPIRIT_CLAW_CONSUME_COUNT = 200

function on_activating_4121006(me, skill, params)
    local visited = {}
    for slot, item in pairs(me:items(InventoryType.Use)) do
        local wz = item:wz()
        if wz and wz:consume_type() == ConsumeType.Shuriken and not visited[wz:id()] then
            local item_id = wz:id()
            local slots = me:item(item_id)
            local count = 0
            for _, it in pairs(slots) do
                count = count + it:count()
            end
            if count >= SPIRIT_CLAW_CONSUME_COUNT then
                me:rmitem(item_id, SPIRIT_CLAW_CONSUME_COUNT)
                return true
            end
            visited[item_id] = true
        end
    end
    return false
end

function on_activated_4121006(me, skill, params)
    local effect = skill:effect()
    local value = 1
    if effect ~= nil then
        value = effect.x
    end
    me:buff(skill, BuffFlag.SpiritClaw, value)
end
