-- Skill name (String.wz/Skill.img.xml): 스피릿 자벨린

local SPIRIT_CLAW_CONSUME_COUNT = 200

local function star_to_consume(me)
    local visited = {}
    for _, item in pairs(me:items(InventoryType.Use)) do
        local wz = item:wz()
        if wz and wz:consume_type() == ConsumeType.Shuriken and not visited[wz:id()] then
            local item_id = wz:id()
            local count = 0
            for _, it in pairs(me:item(item_id)) do
                count = count + it:count()
            end
            if count >= SPIRIT_CLAW_CONSUME_COUNT then
                return item_id
            end
            visited[item_id] = true
        end
    end
    return nil
end

return {
    on_activating = function(me, skill, params)
        return star_to_consume(me) ~= nil
    end,

    on_activated = function(me, skill, params)
        local item_id = star_to_consume(me)
        if item_id == nil then
            return false
        end
        me:rmitem(item_id, SPIRIT_CLAW_CONSUME_COUNT)
        local effect = skill:effect()
        local value = 1
        if effect ~= nil then
            value = effect.x
        end
        me:buff(skill, BuffFlag.SpiritClaw, value)
    end
}
