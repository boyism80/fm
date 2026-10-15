local pq = require("script/integration/lib/party_quest")

local ENTRY = 105090200
local FIELD = 910500000
local NPC = 1061012
local BALROG = 9300085

test_suite {
	name = "Fourth Job: 이계의 궁수 수련장",
	bot_count = 2,

	on_initialize = function(ctx)
		local classes = { 312, 322 }
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local command = string.format("/봇초기화 120 %d 0 - 6107:2,6108:1", classes[i + 1])
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
			local oid = pq.npc(ctx, leader, NPC)
			if oid == false then
				return false
			end
			if leader:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == FIELD
				end, 10000) == false then
				return ctx:fail("이계의 궁수 수련장 입장 실패")
			end
			return pq.wait_all(ctx, FIELD, 10000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			for _, spot in ipairs(bot:mob_spots()) do
				if bot:mobs(BALROG)[1] ~= nil then
					break
				end
				if pq.move(bot, spot.x, spot.y) == false then
					return ctx:fail("몹 위치로 이동 실패")
				end
			end
			local mob = bot:mobs(BALROG)[1]
			if mob == nil then
				return ctx:fail("이계의 주니어 발록이 보이지 않음")
			end
			if bot:kill(mob.oid) == false then
				return ctx:fail("이계의 주니어 발록 처치 실패")
			end
			return true
		end,
		function(ctx)
			for i = 0, ctx:bot_count() - 1 do
				if pq.portal(ctx, ctx:bot(i), "west00", ENTRY) == false then
					return false
				end
			end
			return true
		end,
	},
}
