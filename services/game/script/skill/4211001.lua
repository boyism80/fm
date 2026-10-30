-- Skill name (String.wz/Skill.img.xml): 챠크라

local function make_heal_hp(rate, stat, lower_f, upper_f)
    local low = math.floor(stat * lower_f * rate)
    local hi = math.floor(stat * upper_f * rate)
    local span = hi - low + 1
    if span < 1 then
        span = 1
    end
    return math.floor(math.random() * span) + low
end

return {
    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        local y = effect.y
        if y == nil then
            y = 0
        end
        local v42 = y + 100
        local v38 = math.random(1, 100) + 100
        local hpchange = math.floor((v38 * me:base_luk() * 0.033 + me:base_dex()) * v42 * 0.002)
        hpchange = hpchange + make_heal_hp(y / 100.0, me:total_luk(), 2.3, 3.5)
        if hpchange <= 0 then
            return
        end
        me:add_hp(hpchange)
        me:update_stats({ STAT.Hp })
    end
}
