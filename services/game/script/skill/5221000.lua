-- Skill name (String.wz/Skill.img.xml): 메이플 용사

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_activated_5221000(me, skill, params)
	combat.for_each_near_party_member(me, skill, function(ch)
		util.apply_buff_from_effect(ch, skill, BuffFlag.MapleWarrior, "x")
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end
