-- Skill name (String.wz/Skill.img.xml): 매직 부스터

local util = require("script/lib/skill")

function on_activated_2111005(me, skill, params)
	util.apply_buff_from_effect(me, skill, BuffFlag.Booster, "x")
end

function on_unbuff_2111005(me, skill)
end

