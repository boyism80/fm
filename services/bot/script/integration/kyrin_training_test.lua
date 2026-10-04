local pq = require("script/integration/lib/party_quest")

local NAVIGATION_ROOM = 120000101
local TRAINING = 912010000
local CLEAR = 912010200
local KYRIN = 1090000
local KYRIN_MOB = 9300159

local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1

test_suite {
	name = "Job: 카이린의 훈련장",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 120 510 0 - 6330:1", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(NAVIGATION_ROOM)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KYRIN)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("카이린 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return p.character.map == TRAINING
				end, 10000) == false then
				return ctx:fail("카이린의 훈련장 입장 실패")
			end
			if bot:mobs(KYRIN_MOB)[1] == nil and bot:request(resp.spawn_mob, nil, function(p)
					return p.mob.mob_id == KYRIN_MOB
				end, 5000) == false then
				return ctx:fail("훈련장에 카이린이 없음")
			end
			return true
		end,
		function(ctx)
			return pq.wait_map(ctx, ctx:bot(0), CLEAR, 200000)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, KYRIN)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("훈련을 마친 뒤 카이린 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
					return p.character.map == NAVIGATION_ROOM
				end, 10000) == false then
				return ctx:fail("항해실로 돌아가지 않음")
			end
			return true
		end,
	},
}
