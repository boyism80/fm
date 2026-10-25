-- Skill name (String.wz/Skill.img.xml): MP 리커버리

local function hp_taken(me)
    return math.floor(me:max_hp() / 100) * 10
end

return {
    on_activating = function(me, skill, params)
        return skill:effect() ~= nil and me:hp() > hp_taken(me)
    end,

    on_activated = function(me, skill, params)
        local effect = skill:effect()
        local decrease = hp_taken(me)
        if effect == nil or me:hp() <= decrease then
            return false
        end
        me:hp(me:hp() - decrease, false)

        local mp_gain = math.floor(decrease / 100) * effect.y
        local new_mp = math.min(me:max_mp(), me:mp() + mp_gain)
        me:mp(new_mp, false)
        me:update_stats({ STAT.Hp, STAT.Mp })
    end
}
