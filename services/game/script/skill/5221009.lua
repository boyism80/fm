-- Skill name (String.wz/Skill.img.xml): 마인드 컨트롤

return {
    on_activated = function(me, skill, params)
        local map = me:map()
        if map == nil then
            return
        end
        local effect = skill:effect()
        if effect == nil then
            return
        end
        local duration_ms = effect.time or 0
        if duration_ms <= 0 then
            return
        end
        local prop = tonumber(effect.prop) or 0
        if prop <= 0 then
            return
        end
        local mob_max = tonumber(effect.mob_count) or 0
        if mob_max <= 0 then
            mob_max = 1
        end
        local mobs = map:objects(ObjectType.Mob, { area = me:skill_area(skill) })
        local checked = 0
        for _, mob in ipairs(mobs) do
            if checked >= mob_max then
                break
            end
            checked = checked + 1
            if math.random(1, 100) <= prop then
                mob:buff(MobBuff.Hypnotize, 1, duration_ms, skill, me)
            end
        end
    end
}
