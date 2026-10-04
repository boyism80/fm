local M = {}

function M.command(bot, text, expected)
	return bot:request(resp.notice, req.normal_chat { message = text }, function(p)
		return p.message:find(expected, 1, true) ~= nil
	end)
end

function M.move(bot, x, y)
	bot:move(x, y)
	return M.command(bot, "/좌표", string.format("Position: %d, %d", x, y)) ~= false
end

function M.npc(ctx, bot, template)
	local x, y = bot:npc_position(template)
	if x == nil then
		return ctx:fail(bot:name() .. " 맵에 NPC 배치가 없음: " .. template)
	end
	if M.move(bot, x, y) == false then
		return ctx:fail(bot:name() .. " NPC 앞으로 이동 실패: " .. template)
	end
	local oid = bot:npc(template)
	if oid == nil then
		return ctx:fail(bot:name() .. " NPC가 보이지 않음: " .. template)
	end
	return oid
end

function M.collect(ctx, bot, mob_id, item_id, count)
	local spots = {}
	for _, spot in ipairs(bot:mob_spots()) do
		if mob_id == nil or spot.id == mob_id then
			spots[#spots + 1] = spot
		end
	end

	local visited = 0
	while (bot:items()[item_id] or 0) < count do
		local drop = bot:drops(item_id)[1]
		local mob = bot:mobs(mob_id)[1]
		if drop ~= nil then
			if bot:loot(drop.oid) == false then
				return ctx:fail(bot:name() .. " 줍기 실패: " .. item_id)
			end
			visited = 0
		elseif mob ~= nil then
			if bot:kill(mob.oid) == false then
				return ctx:fail(bot:name() .. " 처치 실패: " .. mob.id)
			end
			bot:request(resp.spawn_item, nil, function(p)
				return p.item_model ~= nil and p.item_model.id == item_id
			end, 2000)
			visited = 0
		elseif visited < #spots then
			local spot = spots[visited % #spots + 1]
			if M.move(bot, spot.x, spot.y) == false then
				return ctx:fail(bot:name() .. " 몹 위치로 이동 실패")
			end
			visited = visited + 1
		else
			local spawn = bot:request(resp.spawn_mob, nil, function(p)
				return mob_id == nil or p.mob.mob_id == mob_id
			end, 30000)
			if spawn == false then
				return ctx:fail(bot:name() .. " 몹 재생성 대기 시간 초과")
			end
		end
	end
	return true
end

function M.give(ctx, from, to, item_id, count)
	local x, y = to:position()
	if M.move(from, x, y) == false then
		return ctx:fail(from:name() .. " " .. to:name() .. "에게 이동 실패")
	end
	local oid = from:drop(item_id, count)
	if oid == nil then
		return ctx:fail(from:name() .. " 아이템 버리기 실패: " .. item_id)
	end
	if to:loot(oid) == false then
		return ctx:fail(to:name() .. " 아이템 줍기 실패: " .. item_id)
	end
	return true
end

function M.combinations(n, k)
	local result = {}
	local picked = {}
	local function pick(start)
		if #picked == k then
			result[#result + 1] = { unpack(picked) }
			return
		end
		for i = start, n do
			picked[#picked + 1] = i
			pick(i + 1)
			picked[#picked] = nil
		end
	end
	pick(1)
	return result
end

return M
