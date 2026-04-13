-- Skill name (String.wz/Skill.img.xml): 트리플 스로우

function on_attack_4121007(me, skill, damages)
	apply_venom(me, damages, Skill.Venom)
end

function on_activated_4121007(me, skill, params)
end
