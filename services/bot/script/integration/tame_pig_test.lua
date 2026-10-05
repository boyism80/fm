local pq = require("script/integration/lib/party_quest")

local ZOO = 230000003
local PIG_ROOM = 923010000
local KENTA_OUTSIDE = 2060005
local KENTA_INSIDE = 9060000
local PIG = 9300102
local REPORT_BOX = 2302005
local PHEROMONE = 4031507
local RESEARCH_REPORT = 4031508

local DIALOG_DEFAULT = 0
local COUNT = 5

test_suite {
	name = "Quest: 켄타의 사육실",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 70 100 0 - 6002:1", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(ZOO)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KENTA_OUTSIDE)
			if oid == false then
				return false
			end
			if bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == PIG_ROOM
				end, 10000) == false then
				return ctx:fail("사육실 입장 실패")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if pq.move(bot, -26, 335) == false then
				return ctx:fail("멧돼지 앞으로 이동 실패")
			end
			if bot:mobs(PIG)[1] == nil then
				return ctx:fail("호위용 멧돼지가 없음")
			end
			for _ = 1, 40 do
				if (bot:items()[PHEROMONE] or 0) >= COUNT then
					return true
				end
				local drop = bot:drops(PHEROMONE)[1]
				if drop == nil then
					bot:request(resp.spawn_item, nil, function(p)
						return p.item_model ~= nil and p.item_model.id == PHEROMONE
					end, 20000)
					drop = bot:drops(PHEROMONE)[1]
				end
				if drop ~= nil then
					bot:loot(drop.oid)
				end
			end
			return ctx:fail("페로몬 부족: " .. (bot:items()[PHEROMONE] or 0))
		end,
		function(ctx)
			local bot = ctx:bot(0)
			for _ = 1, 8 do
				if (bot:items()[RESEARCH_REPORT] or 0) >= COUNT then
					return true
				end
				local result = pq.visit_reactors(ctx, bot, pq.reactor_by_id(REPORT_BOX), function(r)
					if (bot:items()[RESEARCH_REPORT] or 0) >= COUNT then
						return true
					end
					if pq.break_reactor(ctx, bot, r) == false then
						return false
					end
					bot:request(resp.spawn_item, nil, function(p)
						return p.item_model ~= nil and p.item_model.id == RESEARCH_REPORT
					end, 2000)
					local drop = bot:drops(RESEARCH_REPORT)[1]
					if drop ~= nil then
						return bot:loot(drop.oid)
					end
					return true
				end)
				if result == false then
					return false
				end
				if (bot:items()[RESEARCH_REPORT] or 0) < COUNT then
					ctx:sleep(10000)
				end
			end
			return ctx:fail("연구 보고서 부족: " .. (bot:items()[RESEARCH_REPORT] or 0))
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KENTA_INSIDE)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("켄타 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == ZOO
				end, 10000) == false then
				return ctx:fail("5+5 완료 후 아쿠아리움 동물원으로 나가지 않음")
			end
			return true
		end,
	},
}
