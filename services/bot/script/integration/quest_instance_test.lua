local pq = require("script/integration/lib/party_quest")

local DIALOGS = { resp.dialog, resp.dialog_yes_no, resp.dialog_accept, resp.warp }
local LAB_HALL = 926130100

local function setup(ctx, bot, args, map_id)
	if pq.command(bot, "/봇초기화 " .. args, "봇초기화 완료") == false then
		return ctx:fail("봇 초기화 실패: " .. args)
	end
	if bot:map_move(map_id) == false then
		return ctx:fail("맵 이동 실패: " .. map_id)
	end
	return true
end

local function answer(p, name)
	if name == resp.dialog_yes_no then
		return req.dialog { dialog_type = 1, next = true }
	elseif name == resp.dialog_accept then
		return req.dialog { dialog_type = p.enable_escape and 11 or 12, next = true }
	end
	return req.dialog { dialog_type = 0, next = true }
end

local function talk_until_warp(ctx, bot, template, maps)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	local p, name = bot:request(DIALOGS, req.npc_click { oid = oid })
	for _ = 1, 10 do
		if p == false then
			return ctx:fail(string.format("NPC %d 대화 중 응답 없음 (맵 %d)", template, bot:map()))
		end
		if name == resp.warp then
			if maps[p.character.map] == nil then
				return ctx:fail(string.format("NPC %d 대화 뒤 다른 맵으로 이동: %d", template, p.character.map))
			end
			return true
		end
		p, name = bot:request(DIALOGS, answer(p, name))
	end
	return ctx:fail(string.format("NPC %d 대화가 끝나지 않음", template))
end

local function pass(ctx, bot, x, y, portal, map_id)
	if pq.move(bot, x, y) == false then
		return ctx:fail("포탈 앞으로 이동 실패: " .. portal)
	end
	return pq.portal(ctx, bot, portal, map_id)
end

local function hunt(ctx, bot, mob_id, count)
	local spots = {}
	for _, spot in ipairs(bot:mob_spots()) do
		if spot.id == mob_id then
			spots[#spots + 1] = spot
		end
	end
	local killed = 0
	local visited = 0
	while killed < count do
		local mob = bot:mobs(mob_id)[1]
		if mob ~= nil then
			if bot:kill(mob.oid) == false then
				return ctx:fail("처치 실패: " .. mob_id)
			end
			killed = killed + 1
			visited = 0
		elseif visited < #spots then
			visited = visited + 1
			if pq.move(bot, spots[visited].x, spots[visited].y) == false then
				return ctx:fail("몹 위치로 이동 실패: " .. mob_id)
			end
		elseif pq.regen(bot) then
			visited = 0
		elseif bot:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == mob_id
			end, 30000) == false then
			return ctx:fail(string.format("몹 재생성 대기 시간 초과: %d (%d/%d)", mob_id, killed, count))
		end
	end
	return true
end

local function milk(ctx, bot, template)
	local oid = pq.npc(ctx, bot, template)
	if oid == false then
		return false
	end
	if bot:npc_click(oid) == false then
		return ctx:fail("어미 젖소 대화가 오지 않음: " .. template)
	end
	bot:dialog(false)
	return true
end

local function read_documents(ctx, bot, count)
	for _, spot in ipairs(bot:npc_spots()) do
		if (bot:items()[4031797] or 0) >= count then
			return true
		end
		if spot.id == 2112016 then
			if pq.move(bot, spot.x, spot.y) == false then
				return ctx:fail("연구 문서 앞으로 이동 실패")
			end
			local nearest = nil
			for _, npc in ipairs(bot:npcs(2112016)) do
				if nearest == nil or math.abs(npc.x - spot.x) + math.abs(npc.y - spot.y) < math.abs(nearest.x - spot.x) + math.abs(nearest.y - spot.y) then
					nearest = npc
				end
			end
			if nearest == nil then
				return ctx:fail("연구 문서가 보이지 않음")
			end
			if bot:request(resp.inventory_operation, req.npc_click { oid = nearest.oid }) == false then
				return ctx:fail(string.format("연구 문서를 받지 못함 (%d, %d)", spot.x, spot.y))
			end
		end
	end
	if (bot:items()[4031797] or 0) < count then
		return ctx:fail("연구 문서가 모자람: " .. (bot:items()[4031797] or 0))
	end
	return true
end

test_suite {
	name = "Quest: 1인 퀘스트 인스턴스",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		if pq.command(bot, "/무적", "무적 상태: enabled") == false then
			return ctx:fail("무적 설정 실패")
		end
		if pq.command(bot, "/즉사", "즉사 상태: enabled") == false then
			return ctx:fail("즉사 설정 실패")
		end
		return true
	end,
	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "30 100 0 - 2174:2,2175:1", 120000100) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 1092007, { [912000000] = true }) == false then
				return false
			end
			if hunt(ctx, bot, 9300155, 1) == false then
				return false
			end
			local x, y = bot:position()
			if pq.kill_at(ctx, bot, 9300156, x, y) == false then
				return false
			end
			return pass(ctx, bot, -290, -211, "ntq1", 120000100)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "30 100 0 - 2179:2,2180:1", 120000103) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 1092000, { [912000100] = true }) == false then
				return false
			end
			for _, cow in ipairs({ 1092090, 1092091, 1092090 }) do
				if milk(ctx, bot, cow) == false then
					return false
				end
			end
			if (bot:items()[4031850] or 0) == 0 then
				return ctx:fail("신선한 우유를 받지 못함")
			end
			return pass(ctx, bot, 793, 148, "ntq2", 120000103)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "40 100 0 - 3239:1", 220020000) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 2040003, { [922000000] = true }) == false then
				return false
			end
			local opened = pq.visit_reactors(ctx, bot, pq.reactor_by_id(2202001), function(r)
				if pq.break_reactor(ctx, bot, r) == false then
					return false
				end
				return pq.loot_spawn(ctx, bot, 4031092)
			end)
			if opened == false then
				return false
			end
			if (bot:items()[4031092] or 0) < 10 then
				return ctx:fail("기계 부품이 모자람: " .. (bot:items()[4031092] or 0))
			end
			if talk_until_warp(ctx, bot, 2040003, { [220020000] = true }) == false then
				return false
			end
			if bot:quests()[3239] ~= 2 then
				return ctx:fail("기계 부품 퀘스트가 완료되지 않음")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "70 100 0 - 3309:2,3310:1", 261000010) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 2111000, { [926120100] = true }) == false then
				return false
			end
			if pq.collect(ctx, bot, 9300141, 4031698, 1) == false then
				return false
			end
			if pq.feed(ctx, bot, 2619000, 4031698, 1) == false then
				return false
			end
			if pq.pick_up(ctx, bot, 4031709) == false then
				return false
			end
			return pass(ctx, bot, 854, 311, "out00", 261000010)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "70 100 0 4000361:10 3334:2,3335:1", 261000000) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 2111003, { [926120300] = true }) == false then
				return false
			end
			local function plant(r)
				if (bot:items()[4031695] or 0) > 0 then
					return true
				end
				if bot:drop(4000361, 1) == nil then
					return ctx:fail("이슬 버리기 실패")
				end
				if bot:request({ resp.trigger_reactor, resp.destroy_reactor }, nil, function(p)
						return p.reactor.oid == r.oid
					end) == false then
					return ctx:fail("눈꽃 자리가 이슬에 반응하지 않음")
				end
				if pq.break_reactor(ctx, bot, r) == false then
					return false
				end
				local drop = bot:drops(4031695)[1]
				if drop == nil then
					bot:request(resp.spawn_item, nil, function(p)
						return p.item_model ~= nil and p.item_model.id == 4031695
					end, 2000)
					drop = bot:drops(4031695)[1]
				end
				if drop ~= nil and bot:loot(drop.oid) == false then
					return ctx:fail("눈꽃 장미 줍기 실패")
				end
				return true
			end
			for round = 1, 2 do
				if round > 1 and pq.command(bot, "/리액터초기화", "리액터") == false then
					return ctx:fail("리액터 초기화 실패")
				end
				if pq.visit_reactors(ctx, bot, pq.reactor_by_id(2612000), plant) == false then
					return false
				end
			end
			if (bot:items()[4031695] or 0) == 0 then
				return ctx:fail("눈꽃 장미가 피지 않음")
			end
			return pass(ctx, bot, -26, 174, "out00", 261000000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "70 100 0 - 3365:1,3366:1", LAB_HALL) == false then
				return false
			end
			if pass(ctx, bot, 313, 42, "in00", 926130101) == false then
				return false
			end
			for item = 4031780, 4031784 do
				if pq.collect(ctx, bot, 9300154, item, 1) == false then
					return false
				end
			end
			for item = 4031780, 4031784 do
				if pq.feed(ctx, bot, 2612005, item, 1) == false then
					return false
				end
			end
			if pq.pick_up(ctx, bot, 4031798) == false then
				return false
			end
			if pass(ctx, bot, -806, 163, "out00", 926130200) == false then
				return false
			end
			return pass(ctx, bot, -277, 65, "out00", LAB_HALL)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "70 100 0 - 3365:1,3367:1", LAB_HALL) == false then
				return false
			end
			if pass(ctx, bot, 491, 45, "in01", 926130102) == false then
				return false
			end
			if read_documents(ctx, bot, 20) == false then
				return false
			end
			if pass(ctx, bot, -122, 309, "out00", 926130201) == false then
				return false
			end
			return pass(ctx, bot, -277, 65, "out00", LAB_HALL)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "70 100 0 - 3365:1,3368:1", LAB_HALL) == false then
				return false
			end
			if pass(ctx, bot, 666, 43, "in02", 926130103) == false then
				return false
			end
			if hunt(ctx, bot, 9300153, 100) == false then
				return false
			end
			if pass(ctx, bot, 186, 175, "out00", 926130203) == false then
				return false
			end
			return pass(ctx, bot, -277, 65, "out00", LAB_HALL)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "40 100 0 - 3933:1", 260000200) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 2101003, { [926000000] = true }) == false then
				return false
			end
			if hunt(ctx, bot, 9100013, 1) == false then
				return false
			end
			return pass(ctx, bot, -502, 272, "st00", 260000200)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "30 100 0 - -", 100000103) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 1052004, { [103000004] = true }) == false then
				return false
			end
			return pass(ctx, bot, 297, 179, "out00", 200000204)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if setup(ctx, bot, "30 500 0 - 2192:1", 120000101) == false then
				return false
			end
			if talk_until_warp(ctx, bot, 1090000, { [108000500] = true, [108000501] = true }) == false then
				return false
			end
			if pq.collect(ctx, bot, 9001005, 4031857, 15) == false then
				return false
			end
			return talk_until_warp(ctx, bot, 1072008, { [120000101] = true })
		end,
	},
}
