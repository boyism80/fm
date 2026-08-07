local pq = require("script/lib/party_quest")

local PQ_ITEMS = { 4001130, 4001131, 4001132, 4001133, 4001134, 4001135 }
local RANKING_QUEST = 1205
local DURATION_MS = 45 * 60000
local SCALE_LEVEL = 200
local MIN_PARTY = 4
local MIN_LEVEL = 71

local function strip_items(me)
	for _, id in ipairs(PQ_ITEMS) do
		pq.remove_all(id, me)
	end
end

local function shuffle_lab(map)
	if map == nil then
		return
	end
	local reactors = {}
	local positions = {}
	for _, reactor in pairs(map:reactors()) do
		local name = reactor:name()
		if name == nil or not string.find(name, "out", 1, true) then
			local x, y = reactor:position()
			reactors[#reactors + 1] = reactor
			positions[#positions + 1] = { x, y }
		end
	end
	for i = #positions, 2, -1 do
		local j = math.random(i)
		positions[i], positions[j] = positions[j], positions[i]
	end
	for i, reactor in ipairs(reactors) do
		reactor:position(positions[i][1], positions[i][2])
	end
end

local function hit_all_reactors(map)
	if map == nil then
		return
	end
	for _, reactor in pairs(map:reactors()) do
		if reactor ~= nil and reactor:state() == 0 then
			reactor:hit(1)
		end
	end
end

local function clear_fx(map)
	if map ~= nil then
		map:clear_effect()
	end
end

local function init_stage6(sm)
	for y = 0, 3 do
		sm:set_property("stage6_" .. tostring(y), "0")
		for b = 0, 9 do
			for c = 0, 3 do
				sm:set_property(string.format("stage6_%d_%d_%d", y, b, c), "0")
			end
		end
	end
	for y = 0, 3 do
		for i = 0, 9 do
			local found = false
			while not found do
				for x = 0, 3 do
					if found then
						break
					end
					local taken = false
					for z = 0, 3 do
						if sm:get_property(string.format("stage6_%d_%d_%d", z, i, x)) == "1" then
							taken = true
							break
						end
					end
					if not taken and math.random() < 0.25 then
						sm:set_property(string.format("stage6_%d_%d_%d", y, i, x), "1")
						found = true
					end
				end
			end
		end
	end
end

local function pick_investigate(sm, map, npc_id)
	local oids = {}
	for _, n in pairs(map:npcs()) do
		if n ~= nil and n:id() == npc_id then
			oids[#oids + 1] = n:oid()
		end
	end
	if #oids < 2 then
		return
	end
	local a = oids[math.random(#oids)]
	local b = a
	while b == a do
		b = oids[math.random(#oids)]
	end
	sm:set_property("stage1_" .. tostring(a), "2")
	sm:set_property("stage1_" .. tostring(b), "3")
end

local function create(cfg)
	local maps = cfg.maps
	local start_map = cfg.start_map
	local exit_map = cfg.exit_map
	local town_map = cfg.town_map
	local investigate_npc = cfg.investigate_npc
	local letter_item = cfg.letter_item
	local protect_mob = cfg.protect_mob
	local protect_fail_msg = cfg.protect_fail_msg or "보호에 실패했습니다."
	local protect_ok_msg = cfg.protect_ok_msg or "보호에 성공했습니다."
	local marble = cfg.marble
	local door_prefix = cfg.door_prefix
	local stage1_exp = cfg.stage1_exp
	local stage4_exp = cfg.stage4_exp
	local urete_office = cfg.urete_office
	local boss_map = cfg.boss_map
	local urete_map = cfg.urete_map
	local reward_map = cfg.reward_map
	local hallway = cfg.hallway
	local beaker_map = cfg.beaker_map
	local hub_map = cfg.hub_map
	local lab1 = cfg.lab1
	local lab2 = cfg.lab2

	local function end_run(sm)
		sm:finish(exit_map)
	end

	local function clear_props(sm)
		sm:set_property("stage", "0")
		sm:set_property("stage1", "0")
		sm:set_property("stage1_way_clear", "")
		sm:set_property("stage3", "0")
		sm:set_property("stage4", "0")
		sm:set_property("stage5", "0")
		sm:set_property("stage7", "0")
		sm:set_property("summoned", "0")
		sm:set_property("clear_protect", "")
		sm:set_property("urete_boss", "")
		sm:set_property("persuade_urete", "")
	end

	return {
		config = cfg,
		strip_items = strip_items,
		letter_item = letter_item,
		marble = marble,
		door_prefix = door_prefix,
		stage1_exp = stage1_exp,
		stage4_exp = stage4_exp,
		town_map = town_map,
		exit_map = exit_map,
		boss_map = boss_map,
		urete_map = urete_map,
		reward_map = reward_map,
		hub_map = hub_map,
		beaker_map = beaker_map,
		hallway = hallway,
		hit_all_reactors = hit_all_reactors,
		clear_fx = clear_fx,

		on_init = function(group)
			group:set_property("state", "0")
			group:declare_min_players(1)
			group:declare_exit_map(exit_map)
		end,

		on_create = function(sm)
			local group = sm:group()
			group:set_property("state", "1")
			clear_props(sm)
			init_stage6(sm)
			for _, map_id in ipairs(maps) do
				local map = group:map(map_id)
				if map ~= nil then
					map:reset()
					map:respawn(true)
				end
			end
			local s1 = group:map(start_map)
			if s1 ~= nil then
				pick_investigate(sm, s1, investigate_npc)
			end
			shuffle_lab(group:map(lab1))
			shuffle_lab(group:map(lab2))
			local office = group:map(urete_office)
			if office ~= nil then
				office:spawn_npc(2112000, 200, 188)
			end
			local boss = group:map(boss_map)
			if boss ~= nil then
				boss:block_gen(true)
				boss:kill_all_mobs()
			end
			return maps
		end,

		on_start = function(sm)
			sm:start_timer(DURATION_MS)
		end,

		on_player_enter = function(sm, player)
			strip_items(player)
			player:map(start_map)
			player:try_party_quest(RANKING_QUEST)
		end,

		on_changed_map = function(sm, player, map_id)
			if map_id == boss_map and sm:get_property("urete_boss") == "" then
				local group = sm:group()
				local boss = group ~= nil and group:map(boss_map) or nil
				if boss ~= nil then
					boss:spawn_npc(2112010, 242, 150)
					sm:set_property("urete_boss", "1")
				end
			end
		end,

		on_mob_die = function(sm, mob)
			if mob == nil then
				return
			end
			local group = sm:group()
			if group == nil then
				return
			end
			local id = mob:id()
			local hall = group:map(hallway)
			if hall ~= nil and pq.mob_count(hall) == 0 and sm:get_property("stage1_way_clear") == "" then
				clear_fx(hall)
				sm:set_property("stage1_way_clear", "1")
			end
			if id >= 9300142 and id <= 9300146 then
				local office = group:map(urete_office)
				if office ~= nil and pq.mob_count(office) == 0 and sm:get_property("stage5") == "1" then
					sm:set_property("stage5", "2")
					clear_fx(office)
					hit_all_reactors(office)
				end
				return
			end
			if id == protect_mob then
				if sm:get_property("clear_protect") == "" then
					sm:set_property("stage7", "1")
					sm:notice(protect_fail_msg, Msg.PinkText)
				end
				return
			end
			if id == 9300139 or id == 9300140 then
				sm:set_property("clear_protect", "1")
				local boss = group:map(boss_map)
				if boss ~= nil then
					clear_fx(boss)
					boss:block_gen(true)
					boss:kill_all_mobs()
				end
				local ok = sm:get_property("stage7") ~= "1"
				if ok then
					sm:notice(protect_ok_msg)
					pq.party_exp(sm, 45000)
				else
					pq.party_exp(sm, 30000)
				end
				local urete = group:map(urete_map)
				local reward = group:map(reward_map)
				if ok then
					if urete ~= nil then
						urete:spawn_npc(2112002, 232, 150)
					end
					if reward ~= nil then
						reward:spawn_npc(2112003, 157, 128)
						reward:spawn_npc(2112004, 107, 128)
						reward:spawn_npc(2112002, 320, 128)
					end
					if boss ~= nil then
						boss:spawn_npc(2112004, -416, -116)
						boss:spawn_npc(2112003, -300, -126)
					end
				else
					if urete ~= nil then
						urete:spawn_npc(2112001, 232, 150)
					end
					if reward ~= nil then
						reward:spawn_npc(2112009, 111, 128)
						reward:spawn_npc(2112008, 211, 128)
					end
					if boss ~= nil then
						boss:spawn_npc(2112009, -416, -116)
						boss:spawn_npc(2112008, -300, -126)
					end
				end
			end
		end,

		on_left_party = function(sm, player)
			end_run(sm)
		end,

		on_disband_party = function(sm)
			end_run(sm)
		end,

		on_scheduled_timeout = function(sm)
			end_run(sm)
		end,

		on_clear = function(sm)
			end_run(sm)
		end,

		on_finish = function(sm)
			sm:group():set_property("state", "0")
			clear_props(sm)
		end
	}
end

local function exchange_marbles(me, npc)
	local sel = me:dialog_list(npc, "제뉴미스트의 구슬과 알카드노의 구슬을 각각 25개씩 가져오면 호루스의 눈을, 둘 중 어느쪽의 구슬이라도 10개 가져오면 목걸이에 새로운 힘을 부여할 수 있는 지혜의 돌을 만들어 주겠네. 자, 어떤 아이템을 만들겠는가?\r\n#b", {
		"호루스의 눈을 만들어 주세요.",
		"제뉴미스트 구슬로 지혜의 돌을 만들어 주세요.",
		"알카드노 구슬로 지혜의 돌을 만들어 주세요.",
	})
	if sel == nil then
		return false
	end
	local material
	local reward_item
	if sel == 1 then
		material = { [4001159] = 25, [4001160] = 25 }
		reward_item = 1122010
	elseif sel == 2 then
		material = { [4001159] = 10 }
		reward_item = 2041212
	elseif sel == 3 then
		material = { [4001160] = 10 }
		reward_item = 2041212
	else
		return false
	end
	local str = "교환할 아이템을 확인하게.\r\n\r\n"
	if sel == 1 then
		str = str .. "#i4001159# #t4001159# 25개\r\n"
		str = str .. "#i4001160# #t4001160# 25개\r\n"
	elseif sel == 2 then
		str = str .. "#i4001159# #t4001159# 10개\r\n"
	else
		str = str .. "#i4001160# #t4001160# 10개\r\n"
	end
	str = str .. "\r\n위 아이템으로 \r\n"
	str = str .. string.format("#i%d# #t%d# 1개로 교환하겠네. 계속하겠는가?", reward_item, reward_item)
	if not me:dialog_yes_no(npc, str) then
		me:dialog(npc, "흐음, 잘못 선택한건가? 마음이 바뀌면 다시 오게. 아마 오래 만나지 못할게야.", false, false)
		return false
	end
	local code = me:exchange({ item = material }, { item = { [reward_item] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 재료를 분명 제대로 갖고 계신건지 확인해 주시게.", false, false)
		return false
	end
	if code ~= ExchangeResult.OK then
		me:dialog(npc, "인벤토리 공간이 부족한건 아닌지, 혹은 재료를 분명 제대로 갖고 계신건지 확인해 주시게.", false, false)
		return false
	end
	if reward_item == 1122010 then
		local q1205 = me:quest(1205)
		if q1205 ~= nil and q1205:record_ex("have") == nil then
			q1205:record_ex("have", "1")
		end
	end
	return true
end

local M = {
	create = create,
	strip_items = strip_items,
	pq_items = PQ_ITEMS,
	ranking_quest = RANKING_QUEST,
	min_party = MIN_PARTY,
	min_level = MIN_LEVEL,
	scale_level = SCALE_LEVEL,
	hit_all_reactors = hit_all_reactors,
	clear_fx = clear_fx,
	exchange_marbles = exchange_marbles,
}

function M.path_for_map(map_id)
	if map_id >= 926100000 and map_id <= 926100700 then
		return "romeo"
	end
	if map_id >= 926110000 and map_id <= 926110700 then
		return "juliet"
	end
	return nil
end

function M.group_name(map_id)
	local path = M.path_for_map(map_id)
	if path == "romeo" then
		return "romeo_party_quest"
	end
	if path == "juliet" then
		return "juliet_party_quest"
	end
	return nil
end

return M
