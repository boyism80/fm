
function on_activated_4201004(me, skill, params)
end

local function total_damage_steal(hits)
	local total = 0
	for _, amount in ipairs(hits) do
		if amount > 0 then
			total = total + amount
		end
	end
	return total
end

local function roll_percent_steal(prob)
	if prob <= 0 then
		return false
	end
	if prob >= 100 then
		return true
	end
	return math.random(1, 100) <= prob
end

function on_attack_4201004(me, skill, damages)
	local effect = skill:effect()
	if effect == nil then
		return
	end
	local prop = effect.prop
	if prop == nil or prop <= 0 then
		prop = 100
	end
	local map = me:map()
	for mob, hits in pairs(damages) do
		if total_damage_steal(hits) <= 0 then
			goto continue_steal
		end
		if mob:wz().boss then
			goto continue_steal
		end
		if mob:has_stolen() then
			goto continue_steal
		end
		if not roll_percent_steal(prop) then
			goto continue_steal
		end
		local rows = mob:drops()
		if #rows == 0 then
			goto continue_steal
		end
		local candidates = {}
		for i = 1, #rows do
			local e = rows[i]
			if e.item > 0 and e.quest == 0 and math.floor(e.item / 10000) ~= 238 then
				candidates[#candidates + 1] = e
			end
		end
		if #candidates == 0 then
			goto continue_steal
		end
		for i = #candidates, 2, -1 do
			local j = math.random(i)
			candidates[i], candidates[j] = candidates[j], candidates[i]
		end
		local channel = me:channel_drop_rate()
		local bonus = me:bonus_drop_rate()
		local mob_dr = mob:drop_rate()
		for _, e in ipairs(candidates) do
			local prob = e.prob
			if prob > 0 then
				local p = 5 * prob * channel * (bonus / 100) * (mob_dr / 100)
				if p > 1 then
					p = 1
				end
				if math.random() < p then
					local x, y = mob:position()
					local pos = { x + math.random(-20, 20), y }
					local stolen_id = e.item
					local count = 1
					if e.max > 0 and e.min >= 1 and e.max >= e.min then
						count = math.random(e.min, e.max)
					end
					map:spawn_item(stolen_id, count, pos, me, true)
					mob:record_stolen_item(stolen_id)
					break
				end
			end
		end
		::continue_steal::
	end
end
