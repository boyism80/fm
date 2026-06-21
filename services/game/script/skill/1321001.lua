-- Skill name (String.wz/Skill.img.xml): 몬스터 마그넷

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_attack_1321001(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activating_1321001(me, skill, params)
	return true
end

function on_activated_1321001(me, skill, params)
	util.apply_monster_magnet(me, skill, params)
end
