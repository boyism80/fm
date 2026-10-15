local pq = require("script/integration/lib/party_quest")

local ENTRY = 230040001
local FIELD = 923000000
local KARTA = 2060100
local SEA_MOSS = 4000175
local PEARL = 4031472

test_suite {
	name = "Fourth Job: 카르타의 일그러진 차원",
	bot_count = 2,

	on_initialize = function(ctx)
		local classes = { 112, 212 }
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local command = string.format("/봇초기화 120 %d 0 %d:1 6301:1", classes[i + 1], SEA_MOSS)
			if pq.command(bot, command, "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(ENTRY) == false then
				return ctx:fail(bot:name() .. " 입장 맵 이동 실패")
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
			local oid = pq.npc(ctx, leader, KARTA)
			if oid == false then
				return false
			end
			if leader:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == FIELD
				end, 10000) == false then
				return ctx:fail("일그러진 차원 입장 실패")
			end
			return pq.wait_all(ctx, FIELD, 10000)
		end,
		function(ctx)
			return pq.collect(ctx, ctx:bot(0), nil, PEARL, 1)
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				if pq.portal(ctx, ctx:bot(i), "out00", ENTRY) == false then
					return false
				end
			end
			return true
		end,
	},
}
