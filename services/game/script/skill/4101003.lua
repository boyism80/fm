-- Skill name (String.wz/Skill.img.xml): 자벨린 부스터

local util = require("script/lib/skill")

function on_activated_4101003(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_4101003(me, skill)
end

