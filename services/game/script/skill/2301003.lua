-- Skill name (String.wz/Skill.img.xml): 인빈서블

local util = require("script/lib/skill")

function on_activated_2301003(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Invincible, "x")
end

function on_unbuff_2301003(me, skill)
end

