-- Skill name (String.wz/Skill.img.xml): 회복

local TIMER_KEY = "recovery"
local INTERVAL_MS = 5000

-- A full HP ends the recovery at the next tick.
local function on_tick(me, amount)
    if me:hp() >= me:max_hp() then
        me:unbuff(BuffFlag.Recovery)
        return
    end
    me:add_hp(amount)
    me:show_hp_healed_effect(amount)
end

return {
    on_activated = function(me, skill, params)
        local effect = skill:effect()
        if effect == nil then
            return
        end
        me:buff(skill, { [BuffFlag.Recovery] = effect.x })
    end,

    on_buff = function(me, skill)
        local amount = me:buff_value(BuffFlag.Recovery)
        if amount == nil or amount <= 0 then
            return
        end
        me:mktimer(TIMER_KEY, INTERVAL_MS, true, on_tick, amount)
    end,

    on_unbuff = function(me, skill)
        me:rmtimer(TIMER_KEY)
    end
}
