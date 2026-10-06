-- Skill name (String.wz/Skill.img.xml): 드래곤 블러드

local TIMER_KEY = "dragon_blood"
local INTERVAL_MS = 4000

local function dragon_blood_values(skill)
    if skill == nil then
        return 20
    end
    local effect = skill:effect()
    if effect == nil then
        return 0, 0
    end

    return effect.x, effect.pad
end

local function on_tick(me, hp_loss)
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

return {
	on_activated = function(me, skill, params)
		local hp_loss, bonus = dragon_blood_values(skill)
		local buff = me:buff(skill, {[BuffFlag.DragonBlood] = hp_loss, [BuffFlag.WeaponAtk] = bonus})
	end,

	on_buff = function(me, skill)
		local hp_loss = me:buff_value(BuffFlag.DragonBlood)
		if hp_loss == nil then
		    return
		end
		me:mktimer(TIMER_KEY, INTERVAL_MS, true, on_tick, hp_loss)
	end,

	on_unbuff = function(me, skill)
		me:rmtimer(TIMER_KEY)
	end
}
