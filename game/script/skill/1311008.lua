-- Skill name (String.wz/Skill.img.xml): 드래곤 블러드

local TIMER_KEY = "dragon_blood"
local INTERVAL_MS = 4000

local function params(skill)
    if skill == nil then
        return 20
    end
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return 20
    end
    local effect = wz.effects[skill:level()]
    if effect == nil then
        return 0, 0
    end

    return effect.x, effect.pad
end

local function get_attack_bonus(skill)
    local wz = skill:wz()
    if wz == nil or wz.effects == nil then
        return 0
    end
    local effect = wz.effects[skill:level()]
    if effect == nil then
        return 0
    end
    return effect.pad or 0
end

function on_activated(me, skill, params)
    local hp_loss, bonus = params(skill)
    local buff = me:buff(skill, {[BuffFlag.DragonBlood] = hp_loss, [BuffFlag.WeaponAtk] = bonus})
end

function on_buff(me, skill)
    local hp_loss = me:buff_value(BuffFlag.DragonBlood)
    if hp_loss == nil then
        return
    end
    me:mktimer(TIMER_KEY, INTERVAL_MS, true, on_tick, hp_loss)
end

function on_tick(me, hp_loss)
    if me == nil then
        return
    end
    local v = math.floor(hp_loss)
    if v <= 0 then
        return
    end
    if me:hp() > v then
        me:add_hp(-v)
    else
        me:unbuff(BuffFlag.DragonBlood)
    end
end

function on_unbuff(me, skill)
    local success = me:rmtimer(TIMER_KEY)
    if success then
        me:chat('remove timer success')
    else
        me:chat('remove timer failed')
    end
end
