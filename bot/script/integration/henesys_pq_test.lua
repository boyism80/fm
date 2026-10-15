local pq = require("script/integration/lib/party_quest")

local PARK = 100000200
local STAGE = 910010000
local SHORTCUT = 910010100
local BONUS = 910010200
local SHORTCUT_OUT = 910010400
local TORY = 1012112
local TOMMY = 1012113
local GROWLIE = 1012114
local SEED = 4001095
local FLOWER = 9102002
local MOONFLOWER = 9108000
local RICE_CAKE = 4001101
local MOON_BUNNY = 9300061
local DIALOG_DEFAULT = 0
local DIALOG_LIST = 4

local function missing_seed(leader)
	for i = 0, 5 do
		if (leader:items()[SEED + i] or 0) == 0 then
			return SEED + i
		end
	end
	return nil
end

local function shake_flower(ctx, leader, spot, seed)
	if pq.move(leader, spot.x, spot.y) == false then
		return ctx:fail("달맞이꽃 앞으로 이동 실패")
	end
	for _, r in ipairs(leader:reactors(spot.id)) do
		if r.broken == false and r.x == spot.x and r.y == spot.y then
			if pq.break_reactor(ctx, leader, r) == false then
				return false
			end
			leader:request(resp.spawn_item, nil, function(p)
				return p.item_model ~= nil and p.item_model.id == seed
			end, 1500)
		end
	end
	local drop = leader:drops(seed)[1]
	if drop ~= nil and leader:loot(drop.oid) == false then
		return ctx:fail("달맞이꽃 씨앗 줍기 실패: " .. seed)
	end
	return true
end

local function gather_seeds(ctx, leader)
	for _ = 1, 100 do
		local seed = missing_seed(leader)
		if seed == nil then
			return true
		end
		for _, spot in ipairs(leader:reactor_spots()) do
			if spot.id == FLOWER + seed - SEED and (leader:items()[seed] or 0) == 0 then
				if shake_flower(ctx, leader, spot, seed) == false then
					return false
				end
			end
		end
		if (leader:items()[seed] or 0) == 0 then
			ctx:sleep(5000)
		end
	end
	return ctx:fail("달맞이꽃 씨앗을 모두 모으지 못함: " .. tostring(missing_seed(leader)))
end

local function gather_rice_cakes(ctx, leader, count)
	for _ = 1, count * 4 do
		if (leader:items()[RICE_CAKE] or 0) >= count then
			return true
		end
		local drop = leader:drops(RICE_CAKE)[1]
		if drop == nil then
			local spawn = leader:request(resp.spawn_item, nil, function(p)
				return p.item_model ~= nil and p.item_model.id == RICE_CAKE
			end, 20000)
			if spawn == false then
				return ctx:fail("월묘가 떡을 만들지 않음")
			end
			drop = leader:drops(RICE_CAKE)[1]
		end
		if drop ~= nil then
			leader:loot(drop.oid)
		end
	end
	return ctx:fail("월묘의 떡 부족: " .. (leader:items()[RICE_CAKE] or 0))
end

test_suite {
	name = "Party Quest: 헤네시스 PQ",
	bot_count = 3,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			if pq.command(bot, "/레벨바꾸기 20", "레벨 설정") == false then
				return ctx:fail(bot:name() .. " 레벨 설정 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(PARK) == false then
				return false
			end
		end
		return true
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = {
		pq.form_party,
		function(ctx)
			local leader = ctx:bot(0)
			local oid = pq.npc(ctx, leader, TORY)
			if oid == false then
				return false
			end
			local warp = leader:request(resp.warp, req.npc_click { oid = oid }, function(p)
				return p.character.map == STAGE
			end, 15000)
			if warp == false then
				return ctx:fail("달맞이꽃 언덕 입장 실패")
			end
			return pq.wait_all(ctx, STAGE, 15000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			if gather_seeds(ctx, leader) == false then
				return false
			end
			for i = 0, 5 do
				if pq.feed(ctx, leader, MOONFLOWER + i, SEED + i, 1) == false then
					return false
				end
			end
			local bunny = leader:mobs(MOON_BUNNY)[1]
			if bunny == nil and leader:request(resp.spawn_mob, nil, function(p)
				return p.mob.mob_id == MOON_BUNNY
			end, 10000) == false then
				return ctx:fail("월묘가 나타나지 않음")
			end
			return gather_rice_cakes(ctx, leader, 10)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local oid = pq.npc(ctx, leader, GROWLIE)
			if oid == false then
				return false
			end
			local dlg = leader:npc_click(oid)
			if dlg == false or dlg.selections == nil then
				return ctx:fail("어흥이 선택지가 오지 않음")
			end
			dlg = leader:dialog(true, 0)
			if dlg == nil or dlg.text:find("월묘가 만든 떡이 아닌가", 1, true) == nil then
				return ctx:fail("어흥이가 떡을 받지 않음")
			end
			if leader:dialog(true) == nil then
				return ctx:fail("어흥이 두 번째 대화가 오지 않음")
			end
			local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
				return p.character.map == SHORTCUT
			end, 10000)
			if warp == false then
				return ctx:fail("지름길로 이동하지 않음")
			end
			return pq.wait_all(ctx, SHORTCUT, 10000)
		end,
		function(ctx)
			local leader = ctx:bot(0)
			local oid = pq.npc(ctx, leader, TOMMY)
			if oid == false then
				return false
			end
			if leader:npc_click(oid) == false then
				return ctx:fail("토미 대화가 오지 않음")
			end
			local list = leader:dialog(true)
			if list == nil or list.selections == nil then
				return ctx:fail("토미 선택지가 오지 않음")
			end
			local warp = leader:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
				return p.character.map == BONUS
			end, 15000)
			if warp == false then
				return ctx:fail("돼지의 마을 입장 실패")
			end
			return pq.wait_all(ctx, BONUS, 15000)
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				local bot = ctx:bot(i)
				local oid = pq.npc(ctx, bot, TOMMY)
				if oid == false then
					return false
				end
				local dlg = bot:npc_click(oid)
				if dlg == false or dlg.selections == nil then
					return ctx:fail(bot:name() .. " 돼지의 마을 토미 선택지가 오지 않음")
				end
				local out = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
					return p.character.map == SHORTCUT_OUT
				end, 10000)
				if out == false then
					return ctx:fail(bot:name() .. " 돼지의 마을에서 나가지 못함")
				end
				oid = pq.npc(ctx, bot, TORY)
				if oid == false then
					return false
				end
				dlg = bot:npc_click(oid)
				if dlg == false or dlg.selections == nil then
					return ctx:fail(bot:name() .. " 지름길 토리 선택지가 오지 않음")
				end
				local home = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_LIST, next = true, selected = 0 }, function(p)
					return p.character.map == PARK
				end, 10000)
				if home == false then
					return ctx:fail(bot:name() .. " 헤네시스 공원으로 돌아가지 못함")
				end
			end
			return true
		end,
	},
}
