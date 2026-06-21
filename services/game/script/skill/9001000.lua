-- Skill name (String.wz/Skill.img.xml): 힐 + 디스펠

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_activated_9001000(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	if effect.prop <= 0 then
		return
	end
	if effect.prop < 100 and math.random(1, 100) > effect.prop then
		return
	end
	combat.for_each_character_in_skill_area(me, skill, function(ch)
		if ch == nil or not ch:is_alive() then
			return
		end
		util.apply_dispel(ch, skill)
		ch:hp(ch:max_hp())
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
