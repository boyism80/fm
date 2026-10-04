local pq = require("script/integration/lib/party_quest")

local DIALOG_DEFAULT = 0
local DIALOG_YES_NO = 1

local function click_next(ctx, bot, npc, map_id)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	if bot:npc_click(oid) == false then
		return ctx:fail(bot:name() .. " 대화가 오지 않음: " .. npc)
	end
	local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_DEFAULT, next = true }, function(p)
		return p.character.map == map_id
	end, 10000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

local function click_warp(ctx, bot, npc, map_id)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	local warp = bot:request(resp.warp, req.npc_click { oid = oid }, function(p)
		return p.character.map == map_id
	end, 10000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

local function confirm_warp(ctx, bot, npc, map_id)
	local oid = pq.npc(ctx, bot, npc)
	if oid == false then
		return false
	end
	if bot:npc_click(oid) == false then
		return ctx:fail(bot:name() .. " 대화가 오지 않음: " .. npc)
	end
	local warp = bot:request(resp.warp, req.dialog { dialog_type = DIALOG_YES_NO, next = true }, function(p)
		return p.character.map == map_id
	end, 10000)
	if warp == false then
		return ctx:fail(bot:name() .. " 맵 이동 실패: " .. map_id)
	end
	return true
end

local function kill_boss(ctx, bot, mob_id)
	local mob = bot:mobs(mob_id)[1]
	if mob == nil then
		local spawn = bot:request(resp.spawn_mob, nil, function(p)
			return p.mob.mob_id == mob_id
		end, 10000)
		if spawn == false then
			return ctx:fail("보스가 보이지 않음: " .. mob_id)
		end
		mob = { oid = spawn.mob.oid }
	end
	if bot:kill(mob.oid) == false then
		return ctx:fail("보스 처치 실패: " .. mob_id)
	end
	return true
end

test_suite {
	name = "Boss: 대왕지네·포장마차",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 120 100 0 - 4103:1,4014:2,4013:1", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		return true
	end,

	scenarios = {
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(701010321) == false then
				return ctx:fail("검은 숲으로 이동 실패")
			end
			if click_warp(ctx, bot, 9310005, 701010322) == false then
				return false
			end
			if click_next(ctx, bot, 9310006, 701010323) == false then
				return false
			end
			if kill_boss(ctx, bot, 9600009) == false then
				return false
			end
			return confirm_warp(ctx, bot, 9310007, 701010320)
		end,
		function(ctx)
			local bot = ctx:bot(0)
			if bot:map_move(741020100) == false then
				return ctx:fail("야시장 뒷골목으로 이동 실패")
			end
			if click_next(ctx, bot, 9330028, 741020101) == false then
				return false
			end
			if kill_boss(ctx, bot, 9410014) == false then
				return false
			end
			return confirm_warp(ctx, bot, 9330032, 741020100)
		end,
	},
}
