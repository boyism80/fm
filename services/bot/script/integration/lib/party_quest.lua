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

function M.talk(ctx, bot, template, expected)
	local oid = M.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail(bot:name() .. " NPC 대화 응답 없음: " .. template)
	end
	bot:dialog(false)
	if dlg.text:find(expected, 1, true) == nil then
		return ctx:fail(bot:name() .. " 예상과 다른 대화: " .. dlg.text)
	end
	return true
end

function M.find_reactor(bot, oid)
	for _, r in ipairs(bot:reactors()) do
		if r.oid == oid then
			return r
		end
	end
	return nil
end

function M.break_reactor(ctx, bot, reactor)
	for _ = 1, 10 do
		local p, name = bot:hit_reactor(reactor.oid)
		if p == false then
			return ctx:fail(string.format("%s 리액터 타격 응답 없음: %d (맵 %d, 상태 %d)", bot:name(), reactor.id, bot:map(), reactor.state))
		end
		if name == resp.destroy_reactor then
			return true
		end
		local current = M.find_reactor(bot, reactor.oid)
		if current == nil or current.broken then
			return true
		end
	end
	return ctx:fail(string.format("%s 10번 때려도 리액터가 부서지지 않음: %d (맵 %d)", bot:name(), reactor.id, bot:map()))
end

function M.clear_visible(ctx, bot, item_id, skip)
	while true do
		local drop = bot:drops(item_id)[1]
		local mob = bot:mobs()[1]
		local reactor = nil
		for _, r in ipairs(bot:reactors()) do
			if r.broken == false and (skip == nil or skip[r.id] == nil) then
				reactor = r
				break
			end
		end
		if drop ~= nil then
			if bot:loot(drop.oid) == false then
				return ctx:fail(bot:name() .. " 줍기 실패: " .. item_id)
			end
		elseif mob ~= nil then
			if bot:kill(mob.oid) == false then
				return ctx:fail(bot:name() .. " 처치 실패: " .. mob.id)
			end
			bot:request(resp.spawn_item, nil, nil, 2000)
		elseif reactor ~= nil then
			if M.break_reactor(ctx, bot, reactor) == false then
				return false
			end
			bot:request(resp.spawn_item, nil, nil, 2000)
		else
			return true
		end
	end
end

function M.sweep(ctx, bot, item_id, count, skip)
	local spots = {}
	for _, spot in ipairs(bot:mob_spots()) do
		spots[#spots + 1] = spot
	end
	for _, spot in ipairs(bot:reactor_spots()) do
		if skip == nil or skip[spot.id] == nil then
			spots[#spots + 1] = spot
		end
	end
	if M.clear_visible(ctx, bot, item_id, skip) == false then
		return false
	end
	for _, spot in ipairs(spots) do
		if (bot:items()[item_id] or 0) >= count then
			return true
		end
		if M.move(bot, spot.x, spot.y) == false then
			return ctx:fail(bot:name() .. " 위치 이동 실패")
		end
		if M.clear_visible(ctx, bot, item_id, skip) == false then
			return false
		end
	end
	return true
end

function M.form_party(ctx)
	local leader = ctx:bot(0)
	if leader:request(resp.party_created, req.party_operation { operation = PARTY.Create }) == false then
		return ctx:fail("파티 생성 실패")
	end
	for i = 1, ctx:bot_count() - 1 do
		local member = ctx:bot(i)
		local invite = leader:request_on(member, resp.party_invite,
			req.party_operation { operation = PARTY.Invite, target_name = member:name() })
		if invite == false then
			return ctx:fail(member:name() .. " 초대 패킷 없음")
		end
		local joined = member:request(resp.party_update_join,
			req.party_operation { operation = PARTY.AcceptInvite, party_id = invite.party_id })
		if joined == false then
			return ctx:fail(member:name() .. " 파티 가입 실패")
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

function M.wait_map(ctx, bot, map_id, timeout)
	if bot:map() == map_id then
		return true
	end
	local warp = bot:request(resp.warp, nil, function(p)
		return p.character.map == map_id
	end, timeout)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 없음: " .. map_id .. " (현재 " .. bot:map() .. ")")
	end
	return true
end

function M.wait_all(ctx, map_id, timeout)
	for i = 0, ctx:bot_count() - 1 do
		if M.wait_map(ctx, ctx:bot(i), map_id, timeout) == false then
			return false
		end
	end
	return true
end

function M.portal(ctx, bot, portal, map_id)
	local warp = bot:warp(portal)
	if warp == false or bot:map() ~= map_id then
		return ctx:fail(bot:name() .. " 포탈 이동 실패: " .. portal .. " (현재 " .. bot:map() .. ")")
	end
	return true
end

function M.seek_reactor(ctx, bot, match)
	local function visible()
		for _, r in ipairs(bot:reactors()) do
			if match(r) then
				return r
			end
		end
		return nil
	end
	local found = visible()
	for _, spot in ipairs(bot:reactor_spots()) do
		if found ~= nil then
			break
		end
		if M.move(bot, spot.x, spot.y) == false then
			return ctx:fail(bot:name() .. " 리액터 위치로 이동 실패")
		end
		found = visible()
	end
	if found == nil then
		return ctx:fail(bot:name() .. " 리액터를 찾지 못함 (맵 " .. bot:map() .. ")")
	end
	if M.move(bot, found.x, found.y) == false then
		return ctx:fail(bot:name() .. " 리액터 앞으로 이동 실패: " .. found.id)
	end
	return found
end

function M.reactor_by_id(id)
	return function(r)
		return r.id == id
	end
end

function M.feed(ctx, bot, reactor_id, item_id, count)
	local reactor = M.seek_reactor(ctx, bot, M.reactor_by_id(reactor_id))
	if reactor == false then
		return false
	end
	if bot:drop(item_id, count) == nil then
		return ctx:fail(bot:name() .. " 아이템 버리기 실패: " .. item_id)
	end
	local p = bot:request({ resp.trigger_reactor, resp.destroy_reactor }, nil, function(p)
		return p.reactor.oid == reactor.oid
	end, 12000)
	if p == false then
		return ctx:fail(string.format("%s 리액터가 아이템에 반응하지 않음: %d <- %d", bot:name(), reactor_id, item_id))
	end
	return reactor
end

function M.pick_up(ctx, bot, item_id)
	if (bot:items()[item_id] or 0) > 0 then
		return true
	end
	return M.loot_spawn(ctx, bot, item_id)
end

function M.loot_spawn(ctx, bot, item_id)
	local drop = bot:drops(item_id)[1]
	if drop == nil then
		local spawn = bot:request(resp.spawn_item, nil, function(p)
			return p.item_model ~= nil and p.item_model.id == item_id
		end, 5000)
		if spawn == false then
			return ctx:fail(bot:name() .. " 아이템이 떨어지지 않음: " .. item_id)
		end
		drop = bot:drops(item_id)[1]
	end
	if drop == nil or bot:loot(drop.oid) == false then
		return ctx:fail(bot:name() .. " 줍기 실패: " .. item_id)
	end
	return true
end

function M.open_box(ctx, bot, box_id, item_id)
	local box = M.seek_reactor(ctx, bot, M.reactor_by_id(box_id))
	if box == false then
		return false
	end
	if M.break_reactor(ctx, bot, box) == false then
		return false
	end
	return M.pick_up(ctx, bot, item_id)
end

function M.kill_at(ctx, bot, mob_id, x, y)
	if M.move(bot, x, y) == false then
		return ctx:fail(bot:name() .. " 몹 위치로 이동 실패")
	end
	local mob = bot:mobs(mob_id)[1]
	if mob == nil then
		local spawn = bot:request(resp.spawn_mob, nil, function(p)
			return p.mob.mob_id == mob_id
		end, 10000)
		if spawn == false then
			return ctx:fail(bot:name() .. " 몹이 보이지 않음: " .. mob_id)
		end
		mob = { oid = spawn.mob.oid }
	end
	if bot:kill(mob.oid) == false then
		return ctx:fail(bot:name() .. " 처치 실패: " .. mob_id)
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
