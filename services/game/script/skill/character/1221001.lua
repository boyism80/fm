-- Skill name (String.wz/Skill.img.xml): 몬스터 마그넷

function on_attack_1221001(me, skill, damages)
	apply_prob_status_on_skill_hit(me, skill, damages, MobBuff.Stun)
end

function on_activating_1221001(me, skill, params)
	return true
end

function on_activated_1221001(me, skill, params)
	apply_monster_magnet(me, skill, params)
end
