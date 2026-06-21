-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_activated_1301007(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	combat.for_each_near_party_member(me, skill, function(ch)
		ch:buff(skill, {
			[BuffFlag.MaxHp] = effect.x,
			[BuffFlag.MaxMp] = effect.x,
		})
		if ch ~= me then
			ch:show_skill_effect(skill, SkillEffectType.Affected)
		end
	end)
end

function on_buff_1301007(me, skill)
	util.add_hyper_body_bonus(me, skill)
end

function on_unbuff_1301007(me, skill)
	util.remove_hyper_body_bonus(me, skill)
end
