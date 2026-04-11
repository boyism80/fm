-- Skill name (String.wz/Skill.img.xml): 콜드 빔

function on_attack_2201004(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Freeze)
end

function on_activated_2201004(me, skill, params)
end
