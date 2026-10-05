local pq = require("script/integration/lib/party_quest")

local WAITING_ROOM = 123256780
local BOSSES = {
	{ npc = 2038, name = "힐라", map = 123356780, mob = 8870000, x = 164, y = 196 },
	{ npc = 2022, name = "매그너스", map = 123356785, mob = 8880000, x = 2365, y = -1347 },
	{ npc = 2023, name = "스우", map = 123356782, mob = 9801028, x = -110, y = -16 },
	{ npc = 2024, name = "카오스 파풀라투스", map = 123356781, mob = 8820119, x = -393, y = -386 },
	{ npc = 2025, name = "데미안", map = 123356786, mob = 9300890, x = 1222, y = 16 },
	{ npc = 2026, name = "검은 마법사", map = 123356788, mob = 8880500, x = -3, y = 85 },
	{ npc = 2039, name = "루시드", map = 450004150, mob = 8880166, x = 1008, y = 48 },
}

local function enter_and_kill(ctx, bot, boss)
	if bot:map_move(WAITING_ROOM) == false then
		return ctx:fail("보스 대기실로 이동 실패")
	end
	local oid = pq.npc(ctx, bot, boss.npc)
	if oid == false then
		return false
	end
	local dlg = bot:npc_click(oid)
	if dlg == false then
		return ctx:fail(boss.name .. " 입장 대화가 오지 않음")
	end
	local warp = bot:request(resp.warp, req.dialog { dialog_type = dlg.enable_escape and 11 or 12, next = true }, function(p)
		return p.character.map == boss.map
	end, 10000)
	if warp == false then
		return ctx:fail(boss.name .. " 보스 맵으로 이동하지 않음")
	end
	if pq.move(bot, boss.x, boss.y) == false then
		return ctx:fail(boss.name .. " 보스 위치로 이동 실패")
	end
	local mob = bot:mobs(boss.mob)[1]
	if mob == nil then
		local spawn = bot:request(resp.spawn_mob, nil, function(p)
			return p.mob.mob_id == boss.mob
		end, 5000)
		if spawn == false then
			return ctx:fail(boss.name .. " 보스가 소환되지 않음: " .. boss.mob)
		end
		mob = { oid = spawn.mob.oid }
	end
	if bot:kill(mob.oid) == false then
		return ctx:fail(boss.name .. " 처치 실패: " .. boss.mob)
	end
	if bot:mobs(boss.mob)[1] ~= nil then
		return ctx:fail(boss.name .. " 처치 후에도 보스가 남아 있음")
	end
	return true
end

test_suite {
	name = "Boss: 보스 처치 7종",
	bot_count = 1,

	on_initialize = function(ctx)
		local bot = ctx:bot(0)
		if pq.command(bot, "/봇초기화 120 100 0", "봇초기화 완료") == false then
			return ctx:fail("봇 초기화 실패")
		end
		if pq.command(bot, "/플레이어모드", "플레이어 모드: enabled") == false then
			return ctx:fail("플레이어 모드 설정 실패")
		end
		if pq.command(bot, "/무적", "무적 상태: enabled") == false then
			return ctx:fail("무적 설정 실패")
		end
		if pq.command(bot, "/즉사", "즉사 상태: enabled") == false then
			return ctx:fail("즉사 설정 실패")
		end
		return bot:map_move(WAITING_ROOM)
	end,

	on_finished = function(ctx)
		ctx:bot(0):request(resp.party_update_disband, req.party_operation { operation = PARTY.Leave }, nil, 5000)
	end,

	scenarios = {
		pq.form_party,
		function(ctx)
			local passed = true
			for _, boss in ipairs(BOSSES) do
				if enter_and_kill(ctx, ctx:bot(0), boss) == false then
					passed = false
				end
			end
			if ctx:bot(0):map_move(WAITING_ROOM) == false then
				return false
			end
			return passed
		end,
	},
}
