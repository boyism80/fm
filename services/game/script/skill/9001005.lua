-- Skill name (String.wz/Skill.img.xml): 부활

local combat = require("script/lib/combat")

return {
	on_activated = function(me, skill, params)
		combat.for_each_character_in_skill_area(me, skill, function(ch)
			if ch == nil or ch:is_alive() then
				return
			end
			ch:stance(0)
			ch:hp(ch:max_hp())
			ch:mp(ch:max_mp())
			if ch ~= me then
				ch:show_skill_effect(skill, SkillEffectType.Affected)
			end
		end)
	end
}
