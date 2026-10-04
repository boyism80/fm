local pq = require("script/integration/lib/party_quest")

local SLEEPY_DUNGEON = 105040305
local PASSAGE = 108010100
local ARENA = 108010101
local CRACK = 1061009
local CRYSTAL = 1061010
local CLONE = 9001002
local BLACK_CHARM = 4031059

local DIALOG_YES_NO = 1

test_suite {
	name = "Job: 3차 전직 차원의 균열",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 70 310 0 - - 195000=job3_trial1_2", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return bot:map_move(SLEEPY_DUNGEON)
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, CRACK)
			if oid == false then
				return false
			end
			if bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
					return p.character.map == PASSAGE
				end, 10000) == false then
				return ctx:fail("빛나는수정의통로 입장 실패")
			end
			return pq.portal(ctx, bot, "in00", ARENA)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local spot = bot:mob_spots()[1]
			if spot == nil then
				return ctx:fail("헬레나의 분신 배치가 없음")
			end
			if pq.kill_at(ctx, bot, CLONE, spot.x, spot.y) == false then
				return false
			end
			return pq.loot_spawn(ctx, bot, BLACK_CHARM)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			local oid = pq.npc(ctx, bot, CRYSTAL)
			if oid == false then
				return false
			end
			if bot:npc_click(oid) == false then
				return ctx:fail("빛나는 수정 대화가 오지 않음")
			end
			if bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
					return p.character.map == SLEEPY_DUNGEON
				end, 10000) == false then
				return ctx:fail("슬리피던전으로 돌아가지 않음")
			end
			return true
		end,
	},
}
