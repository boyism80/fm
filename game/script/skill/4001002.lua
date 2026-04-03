-- Skill name (String.wz/Skill.img.xml): 디스오더

local function hits_have_damage(hits)
	if hits == nil then
		return false
	end
	for i = 1, #hits do
		local v = hits[i]
		if v ~= nil and v > 0 then
			return true
		end
	end
	return false
end

function on_attack_4001002(me, skill, damages)
	if me == nil or skill == nil or damages == nil then
		return
	end
	local effect = get_skill_effect(skill)
	if effect == nil then
		return
	end
	local watk = math.floor(tonumber(effect.x) or 0)
	local wdef = math.floor(tonumber(effect.y) or 0)
	local duration_sec = math.floor(tonumber(effect.time) or 0)
	if duration_sec <= 0 then
		return
	end
	local duration_ms = duration_sec * 1000
	for mob, hits in pairs(damages) do
		if mob ~= nil and hits_have_damage(hits) then
			mob:set_status({ [MobStatus.Watk] = watk, [MobStatus.Wdef] = wdef }, duration_ms, skill, me)
		end
	end
end

function on_activated_4001002(me, skill, params)
end
