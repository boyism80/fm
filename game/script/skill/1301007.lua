-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

function on_activated(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local percent = effect.x or 0
    me:buff(skill, {
        [BuffFlag.MaxHp] = percent,
        [BuffFlag.MaxMp] = percent,
    })

    -- TODO: Apply this buff to nearby party members when party system is implemented.
end

function on_buff(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local percent = effect.x or 0
    local hp_fixed, hp_percent = me:bonus_max_hp()
    local mp_fixed, mp_percent = me:bonus_max_mp()
    me:bonus_max_hp(hp_fixed, hp_percent + percent)
    me:bonus_max_mp(mp_fixed, mp_percent + percent)
end

function on_unbuff(me, skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return
    end

    local effect = wz.effects[skill:level()]
    if effect == nil then
        return
    end

    local percent = effect.x or 0
    local hp_fixed, hp_percent = me:bonus_max_hp()
    local mp_fixed, mp_percent = me:bonus_max_mp()
    me:bonus_max_hp(hp_fixed, hp_percent - percent)
    me:bonus_max_mp(mp_fixed, mp_percent - percent)
end
