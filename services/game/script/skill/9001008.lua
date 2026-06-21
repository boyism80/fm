-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

local util = require("script/lib/skill")

function on_activated_9001008(me, skill, params)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	me:buff(skill, {
		[BuffFlag.MaxHp] = effect.x,
		[BuffFlag.MaxMp] = effect.x,
	})
end

function on_buff_9001008(me, skill)
	util.add_hyper_body_bonus(me, skill)
end

function on_unbuff_9001008(me, skill)
	util.remove_hyper_body_bonus(me, skill)
end
