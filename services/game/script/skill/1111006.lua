-- Skill name (String.wz/Skill.img.xml): 코마 : 도끼

function on_attack_1111006(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activated_1111006(me, skill, params)
	-- Coma (Axe) shares the same combo orb consumption behavior as Coma (Sword).
	consume_combo_orbs(me)
end
