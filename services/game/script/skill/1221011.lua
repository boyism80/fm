-- Skill name (String.wz/Skill.img.xml): 생츄어리

function on_attack_1221011(me, skill, damages)
	for mob, hits in pairs(damages) do
		if mob and hits then
			local wz = mob:wz()
			local is_boss = wz and wz.boss
			local d
			if is_boss then
				d = 5000000
			else
                d = mob:hp() - 1
			end
			hits[1] = d
		end
	end
end
