local pq = require("script/integration/lib/party_quest")

local ENTRY = 220050300
local BOSS_MAP = 922020100
local FLO = 2041023
local THANATOS = { 9300086, 9300100 }

test_suite {
	name = "Fourth Job: 속성의 타나토스",
	bot_count = 2,

	on_initialize = function(ctx)
		local classes = { 212, 222 }
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local command = string.format("/봇초기화 120 %d 0 - 6225:1,6226:2", classes[i + 1])
			if pq.command(bot, command, "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(ENTRY) == false then
				return ctx:fail(bot:name() .. " 시간의 통로 이동 실패")
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
			local oid = pq.npc(ctx, leader, FLO)
			if oid == false then
				return false
			end
			if leader:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == BOSS_MAP
				end, 10000) == false then
				return ctx:fail("타나토스의 방 입장 실패")
			end
			return pq.wait_all(ctx, BOSS_MAP, 10000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local killed = 0
			for _, spot in ipairs(bot:mob_spots()) do
				for _, id in ipairs(THANATOS) do
					if spot.id == id then
						if pq.kill_at(ctx, bot, id, spot.x, spot.y) == false then
							return false
						end
						killed = killed + 1
					end
				end
			end
			if killed == 0 then
				return ctx:fail("속성의 타나토스 배치가 없음")
			end
			return true
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
