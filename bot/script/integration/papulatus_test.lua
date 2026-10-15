local pq = require("script/integration/lib/party_quest")

local LOBBY = 220080000
local BOSS = 220080001
local MEDAL = 4031172
local CRACK = 4031179
local MACHINE = 2041025
local DROP_X = -124
local DROP_Y = -538
local PHASES = { 8500000, 8500001, 8500002 }
local DIALOG_YES_NO = 1

local function kill_phase(ctx, bot, mob_id)
	local mob = bot:mobs(mob_id)[1]
	if mob == nil then
		local spawn = bot:request(resp.spawn_mob, nil, function(p)
			return p.mob.mob_id == mob_id
		end, 15000)
		if spawn == false then
			return ctx:fail("파풀라투스 단계가 나타나지 않음: " .. mob_id)
		end
		mob = { oid = spawn.mob.oid }
	end
	if bot:kill(mob.oid) == false then
		return ctx:fail("파풀라투스 단계 처치 실패: " .. mob_id)
	end
	return true
end

local function leave(ctx, bot)
	local oid = pq.npc(ctx, bot, MACHINE)
	if oid == false then
		return false
	end
	if bot:npc_click(oid) == false then
		return ctx:fail(bot:name() .. " 기계장치 대화가 오지 않음")
	end
	local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
		return p.character.map == LOBBY
	end, 10000)
	if warp == false then
		return ctx:fail(bot:name() .. " 기계장치로 나가지 못함")
	end
	return true
end

test_suite {
	name = "Boss: 파풀라투스",
	bot_count = 2,

	on_initialize = function(ctx)
		for i = 0, ctx:bot_count() - 1 do
			local bot = ctx:bot(i)
			local items = string.format("%d:1,%d:1", MEDAL, CRACK)
			if pq.command(bot, "/봇초기화 120 100 0 " .. items, "봇초기화 완료") == false then
				return ctx:fail(bot:name() .. " 봇 초기화 실패")
			end
			if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
				return ctx:fail(bot:name() .. " 플레이어 모드 설정 실패")
			end
			if bot:map_move(LOBBY) == false then
				return false
			end
		end
		return true
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if pq.portal(ctx, bot, "in00", BOSS) == false then
				return false
			end
			if pq.move(bot, DROP_X, DROP_Y) == false then
				return ctx:fail("균열 위치로 이동 실패")
			end
			if bot:drop(CRACK, 1) == nil then
				return ctx:fail("차원 균열의 조각 버리기 실패")
			end
			return kill_phase(ctx, bot, PHASES[1])
		end,
		function(ctx)
			local late = ctx:bot(1)
			late:warp("in00")
			if late:map() ~= LOBBY then
				return ctx:fail("전투 중인데 입장됨")
			end
			return true
		end,
		function(ctx)
			local bot = ctx:bot(0)
			for i = 2, #PHASES do
				if kill_phase(ctx, bot, PHASES[i]) == false then
					return false
				end
			end
			return leave(ctx, bot)
		end,
		function(ctx)
			local bot = ctx:bot(1)
			for i = 1, 5 do
				bot:warp("in00")
				if bot:map() == BOSS then
					break
				end
				ctx:sleep(2000)
			end
			if bot:map() ~= BOSS then
				return ctx:fail("전투가 끝난 뒤 다시 입장하지 못함")
			end
			return leave(ctx, bot)
		end,
	},
}
