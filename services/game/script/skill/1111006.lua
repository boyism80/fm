-- Skill name (String.wz/Skill.img.xml): 코마 : 도끼

local util = require("script/lib/skill")
local combat = require("script/lib/combat")

function on_attack_1111006(me, skill, damages)
	combat.apply_prob_status(me, skill, damages, MobBuff.Stun)
end

function on_activated_1111006(me, skill, params)
	-- Coma (Axe) shares the same combo orb consumption behavior as Coma (Sword).
	util.consume_combo_orbs(me)
end
