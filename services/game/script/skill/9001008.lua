-- Skill name (String.wz/Skill.img.xml): 하이퍼 바디

local util = require("script/lib/skill")

return {
	on_activated = function(me, skill, params)
		local effect = skill:effect()
		if effect == nil then
			return
		end
		me:buff(skill, {
			[BuffFlag.MaxHp] = effect.x,
			[BuffFlag.MaxMp] = effect.x,
		})
	end,

	on_buff = function(me, skill)
		util.add_hyper_body_bonus(me, skill)
	end,

	on_unbuff = function(me, skill)
		util.remove_hyper_body_bonus(me, skill)
	end
}
