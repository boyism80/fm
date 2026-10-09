-- Skill name (String.wz/Skill.img.xml): 배틀쉽

local VEHICLE = 1932000
local GAUGE = 5221999

return {
	manual_cooldown = true,

	on_activated = function(me, skill, params)
		return me:ride(skill, VEHICLE)
	end,

	on_ride = function(me, skill)
		local max_hp = math.max(200 * (me:level() + 2 * skill:level() - 120), 1)
		local _, _, _, hp = me:mount()
		if hp == 0 or hp > max_hp then
			hp = max_hp
			me:set_mount_hp(hp)
		end
		if hp < max_hp then
			me:show_skill_cooldown(GAUGE, hp)
		end
	end,

	on_rider_damaged = function(me, skill, damage)
		local _, _, _, hp = me:mount()
		hp = math.max(hp - damage, 0)
		me:set_mount_hp(hp)
		me:show_skill_cooldown(GAUGE, hp)
		if hp > 0 then
			return
		end
		me:dismount()
		skill:cooldown(skill:effect().cooldown)
	end
}
