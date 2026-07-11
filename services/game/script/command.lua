


local function string_split(s, sep)
	sep = sep or " "
	local t = {}
	local pattern = "([^" .. sep .. "]+)"
	for v in string.gmatch(s, pattern) do
		table.insert(t, v)
	end
	return t
end

local function is_sponge_root(mob)
	local children = mob:children()
	return children ~= nil and #children > 0
end

local function resolve_skill_entry(me, skill_arg)
	if skill_arg == nil or skill_arg == "" then
		return nil, "스킬ID 또는 스킬이름이 필요합니다."
	end
	local skill_id = tonumber(skill_arg)
	if skill_id == nil then
		local wz = name2skill(skill_arg)
		if wz == nil then
			return nil, "존재하지 않는 스킬입니다: " .. tostring(skill_arg)
		end
		skill_id = wz.id
	end
	local skill = me:skill(skill_id)
	if skill == nil then
		return nil, string.format("배운 스킬이 아닙니다: %d", skill_id)
	end
	return skill, nil
end

function effect_show_skill(me, skill_arg, skill_effect_type)
	local skill, err = resolve_skill_entry(me, skill_arg)
	if skill == nil then
		me:notice(err)
		return true
	end
	if skill_effect_type == nil then
		skill_effect_type = SkillEffectType.Cast
	end
	me:show_skill_effect(skill, skill_effect_type)
	return true
end

function effect_show_basic(me, effect_type)
	if effect_type == nil then
		effect_type = EffectType.LevelUp
	end
	me:show_effect(effect_type)
	return true
end

function effect_show_dragon_blood(me, skill_arg)
	local skill, err = resolve_skill_entry(me, skill_arg)
	if skill == nil then
		me:notice(err)
		return true
	end
	me:show_dragon_blood_effect(skill)
	return true
end

function effect_show_hp_healed(me, amount)
	if amount == nil then
		amount = 1
	end
	me:show_hp_healed_effect(tonumber(amount) or 1)
	return true
end

function effect_show_reward_item_animation(me, item_id, effect_text)
	local id = tonumber(item_id)
	if id == nil or id <= 0 then
		me:notice("item_id는 1 이상의 숫자여야 합니다.")
		return true
	end
	if effect_text == nil or effect_text == "" then
		effect_text = "Effect/BasicEff.img/LevelUp"
	end
	me:show_reward_item_animation(id, effect_text)
	return true
end

function effect_show_item_maker_success(me)
	me:show_item_maker_success_effect()
	return true
end

function effect_show_crafting(me, effect_text, time_value, mode_value)
	if effect_text == nil or effect_text == "" then
		effect_text = "Effect/BasicEff.img/LevelUp"
	end
	time_value = tonumber(time_value) or 0
	mode_value = tonumber(mode_value) or 0
	me:show_crafting_effect(effect_text, time_value, mode_value)
	return true
end

function effect_show_dice(me, effect_id, skill_arg)
	local skill, err = resolve_skill_entry(me, skill_arg)
	if skill == nil then
		me:notice(err)
		return true
	end
	effect_id = tonumber(effect_id) or 1
	me:show_dice_effect(effect_id, skill)
	return true
end

local function item_count(me, item_id)
	local slots = me:item(item_id)
	if slots == nil then
		return 0
	end
	local total = 0
	for _, it in pairs(slots) do
		if it ~= nil then
			total = total + it:count()
		end
	end
	return total
end

local unsupported_quest_req_kinds = {
	pet = true,
	pettamenessmin = true,
	mbmin = true,
	mbcard = true,
	subJobFlags = true,
	dayByDay = true,
	normalAutoStart = true,
	partyQuest_S = true,
	fieldEnter = true,
	interval = true,
	start = true,
	["end"] = true,
}

local function class_matches(class_id, codes)
	if codes == nil then
		return true
	end
	for _, code in ipairs(codes) do
		if code == 0 and class_id == 0 then
			return true
		end
		if code == class_id then
			return true
		end
		if code % 100 == 0 and math.floor(class_id / 100) == math.floor(code / 100) then
			return true
		end
	end
	return false
end

local function pick_class(codes)
	for _, code in ipairs(codes) do
		if code % 100 ~= 0 then
			return code
		end
	end
	return codes[1]
end

local function ensure_quest_state(me, quest_id, state)
	local other = me:quest(quest_id)
	if other == nil then
		return false
	end
	if state == 2 then
		if other:completed() then
			return true
		end
		if not other:started() then
			if not other:start(0, true) then
				return false
			end
			other = me:quest(quest_id)
			if other == nil then
				return false
			end
		end
		return other:force_complete(0)
	end
	if state == 1 then
		if other:status() == 1 then
			return true
		end
		if other:completed() or other:started() then
			me:clear_quests(quest_id)
			other = me:quest(quest_id)
			if other == nil then
				return false
			end
		end
		return other:start(0, true)
	end
	if state == 0 then
		if other:status() == 0 then
			return true
		end
		me:clear_quests(quest_id)
		other = me:quest(quest_id)
		return other ~= nil and other:status() == 0
	end
	return false
end

local function prepare_quest_requirement(me, quest, req)
	local kind = req.kind
	if kind == "item" then
		local items = req.items
		if items == nil then
			return true
		end
		for _, entry in ipairs(items) do
			local item_id = entry.id
			local need = entry.count or 0
			if item_id ~= nil and need > 0 then
				local have = item_count(me, item_id)
				if have < need then
					me:mkitem(item_id, need - have)
				end
			end
		end
		return true
	end
	if kind == "mob" then
		local mobs = req.mobs
		if mobs == nil then
			return true
		end
		local changed = false
		for _, entry in ipairs(mobs) do
			local mob_id = entry.id
			local need = entry.count or 0
			if mob_id ~= nil and need > 0 then
				local have = quest:mob_kills(mob_id)
				if have < need then
					quest:set_mob_kills(mob_id, need)
					changed = true
				end
			end
		end
		if changed then
			quest:sync_progress()
		end
		return true
	end
	if kind == "pop" then
		local need = req.value or 0
		if me:population() < need then
			me:population(need)
		end
		return true
	end
	if kind == "lvmin" then
		local need = req.value or 0
		if me:level() < need then
			me:level(need)
		end
		return true
	end
	if kind == "lvmax" then
		local max_level = req.value or 0
		if max_level > 0 and me:level() > max_level then
			me:level(max_level)
		end
		return true
	end
	if kind == "npc" or kind == "startscript" or kind == "endscript" then
		return true
	end
	if kind == "quest" then
		local quests = req.quests
		if quests == nil then
			return true
		end
		for _, entry in ipairs(quests) do
			if not ensure_quest_state(me, entry.id, entry.state or 0) then
				return false, kind
			end
		end
		return true
	end
	if kind == "class" then
		local classes = req.classes
		if classes == nil or #classes == 0 then
			return true
		end
		if not class_matches(me:class(), classes) then
			local target = pick_class(classes)
			if target == nil then
				return false, kind
			end
			me:class(target)
		end
		return true
	end
	if unsupported_quest_req_kinds[kind] then
		return false, kind
	end
	if req.value ~= nil and req.value ~= 0 then
		return false, kind
	end
	return true
end

local function prepare_quest_complete(me, quest)
	local wz = quest:wz()
	if wz == nil then
		return false, "wz"
	end
	local reqs = wz.complete_requirements
	if reqs == nil then
		return true
	end
	for _, req in ipairs(reqs) do
		local ok, reason = prepare_quest_requirement(me, quest, req)
		if not ok then
			return false, reason or req.kind
		end
	end
	return true
end

local function prepare_quest_start(me, quest)
	local wz = quest:wz()
	if wz == nil then
		return false, "wz"
	end
	local reqs = wz.start_requirements
	if reqs == nil then
		return true
	end
	for _, req in ipairs(reqs) do
		local ok, reason = prepare_quest_requirement(me, quest, req)
		if not ok then
			return false, reason or req.kind
		end
	end
	return true
end

local function start_npc_id(wz)
	if wz == nil or wz.start_requirements == nil then
		return nil
	end
	for _, req in ipairs(wz.start_requirements) do
		if req.kind == "npc" then
			local npc_id = req.value or 0
			if npc_id ~= 0 then
				return npc_id
			end
		end
	end
	return nil
end

command_funcs = {
	["명령어"] = {
		privilege = ROLE.User,
		usage = "- 사용 가능한 명령어 목록 표시",
		command = function(me, args)
			local user_role = me:role()
			local list = {}
			for cmd_name, cmd_data in pairs(command_funcs) do
				local priv = ROLE.User
				local usage_str = ""
				if type(cmd_data) == "table" then
					priv = cmd_data.privilege or ROLE.Admin
					usage_str = cmd_data.usage or ""
				else
					priv = ROLE.Admin
					usage_str = "- (설명 없음)"
				end
				if user_role >= priv then
					table.insert(list, { name = cmd_name, usage = usage_str })
				end
			end
			table.sort(list, function(a, b) return a.name < b.name end)
			me:notice("=== 사용 가능한 명령어 목록 ===")
			for i, c in ipairs(list) do
				me:notice(string.format("%d. /%s %s", i, c.name, c.usage))
			end
			me:notice(string.format("총 %d개의 명령어가 있습니다.", #list))
			return true
		end,
	},
	["아이템생성"] = {
		privilege = ROLE.Admin,
		usage = "<아이템이름> [개수] - 아이템 생성",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /아이템생성 <아이템이름> [개수]")
				return true
			end
			local count = 1
			if args[2] then
				count = tonumber(args[2]) or 1
			end
			local item_id = tonumber(args[1])
			if item_id == nil then
				item_id = name2item(args[1])
			end
			if item_id == nil then
				me:notice("존재하지 않는 아이템입니다: " .. args[1])
				return true
			end
			local item = me:mkitem(item_id, count)
			if item == nil then
				me:notice("아이템 생성 실패: " .. args[1])
				return true
			end
			me:notice(string.format("아이템 생성: %s x%d", args[1], count))
			return true
		end,
	},
	["인벤토리초기화"] = {
		privilege = ROLE.Admin,
		usage = "- 인벤토리 아이템 전부 제거 (착용 장비 제외)",
		command = function(me, args)
			local cleared = me:clear_inventory()
			me:notice(string.format("인벤토리 %d슬롯 비움", cleared))
			return true
		end,
	},
	["메소초기화"] = {
		privilege = ROLE.Admin,
		usage = "- 메소 초기화",
		command = function(me, args)
			me:meso(0)
			me:notice("메소를 초기화했습니다.")
			return true
		end,
	},
	["인기도"] = {
		privilege = ROLE.Admin,
		usage = "<값> - 인기도 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /인기도 <값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then
				v = 0
			end
			if v > 65535 then
				v = 65535
			end
			me:population(v)
			me:notice(string.format("인기도 설정: %d", v))
			return true
		end,
	},
	["퀘스트초기화"] = {
		privilege = ROLE.Admin,
		usage = "[퀘스트ID] - 진행·완료 퀘스트 제거 (생략 시 전체)",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				local count = me:clear_quests()
				me:notice(string.format("퀘스트 %d개 초기화", count))
			else
				local quest_id = tonumber(args[1])
				if not quest_id then
					me:notice("사용법: /퀘스트초기화 [퀘스트ID]")
					return true
				end
				if me:clear_quests(quest_id) == 1 then
					me:notice(string.format("퀘스트 %d 초기화", quest_id))
				else
					me:notice(string.format("퀘스트 %d 없음", quest_id))
				end
			end
			return true
		end,
	},
	["퀘스트완료준비"] = {
		privilege = ROLE.Admin,
		usage = "[퀘스트ID] - 완료 조건 충족 (생략 시 진행 중 전체)",
		command = function(me, args)
			local targets = {}
			if args[1] and args[1] ~= "" then
				local quest_id = tonumber(args[1])
				if not quest_id then
					me:notice("사용법: /퀘스트완료준비 [퀘스트ID]")
					return true
				end
				local quest = me:quest(quest_id)
				if quest == nil or not quest:started() then
					me:notice(string.format("퀘스트 %d 진행 중 아님", quest_id))
					return true
				end
				targets[quest_id] = quest
			else
				targets = me:quests()
			end
			local prepared = 0
			local skipped = 0
			for quest_id, quest in pairs(targets) do
				local ok, reason = prepare_quest_complete(me, quest)
				if ok then
					prepared = prepared + 1
					me:notice(string.format("퀘스트 %d 준비 완료", quest_id))
				else
					skipped = skipped + 1
					me:notice(string.format("퀘스트 %d 스킵 (%s)", quest_id, tostring(reason)))
				end
			end
			if prepared == 0 and skipped == 0 then
				me:notice("진행 중인 퀘스트 없음")
			else
				me:notice(string.format("완료 준비: %d개 성공, %d개 스킵", prepared, skipped))
			end
			return true
		end,
	},
	["퀘스트시작준비"] = {
		privilege = ROLE.Admin,
		usage = "<퀘스트ID> - 시작 조건 충족 후 시작 NPC 근처로 이동",
		command = function(me, args)
			local quest_id = tonumber(args[1])
			if not quest_id then
				me:notice("사용법: /퀘스트시작준비 <퀘스트ID>")
				return true
			end
			local quest = me:quest(quest_id)
			if quest == nil then
				me:notice(string.format("퀘스트 %d 없음", quest_id))
				return true
			end
			if quest:started() then
				me:notice(string.format("퀘스트 %d 이미 진행 중", quest_id))
				return true
			end
			if quest:completed() then
				me:notice(string.format("퀘스트 %d 이미 완료됨 (필요 시 /퀘스트초기화)", quest_id))
				return true
			end
			local ok, reason = prepare_quest_start(me, quest)
			if not ok then
				me:notice(string.format("퀘스트 %d 시작 준비 실패 (%s)", quest_id, tostring(reason)))
				return true
			end
			me:notice(string.format("퀘스트 %d 시작 조건 준비 완료", quest_id))
			local wz = quest:wz()
			local npc_id = start_npc_id(wz)
			if npc_id == nil then
				me:notice("시작 NPC 없음")
				return true
			end
			local spawns = npc_spawns(npc_id)
			if spawns == nil or #spawns == 0 then
				me:notice(string.format("시작 NPC %d 스폰을 찾지 못함", npc_id))
				return true
			end
			local spawn = spawns[1]
			local spawn_point = closest_spawn(spawn.map_id, spawn.x, spawn.y)
			if spawn_point == nil then
				me:map(spawn.map_id)
			else
				me:map(spawn.map_id, spawn_point)
			end
			me:notice(string.format("NPC %d → 맵 %d (%d, %d) 근처 스폰으로 이동", npc_id, spawn.map_id, spawn.x, spawn.y))
			return true
		end,
	},
	["메소얻기"] = {
		privilege = ROLE.Admin,
		usage = "<금액> - 메소 획득",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /메소얻기 <금액>")
				return true
			end
			local amount = tonumber(args[1])
			if not amount or amount < 0 then
				me:notice("금액은 0 이상의 숫자여야 합니다.")
				return true
			end
			me:meso(me:meso() + amount)
			me:notice(string.format("메소 %d 획득.", amount))
			return true
		end,
	},
	["풀메소"] = {
		privilege = ROLE.Admin,
		usage = "- 메소 최대치로 설정",
		command = function(me, args)
			me:meso(2147483647)
			me:notice("메소를 최대치로 설정했습니다.")
			return true
		end,
	},
	["맵이동"] = {
		privilege = ROLE.Admin,
		usage = "<맵이름|맵ID> [스폰포인트] - 맵 이동",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /맵이동 <맵이름|맵ID> [스폰포인트]")
				return true
			end
			local spawn = 1
			if args[2] then
				spawn = tonumber(args[2]) or 1
			end
			local map_arg = args[1]
			local map_id = tonumber(map_arg)
			local map_info
			if map_id ~= nil then
				map_info = id2map(map_id)
			else
				map_info = name2map(map_arg)
			end
			if map_info == nil then
				me:notice("존재하지 않는 맵입니다: " .. tostring(map_arg))
				return true
			end
			me:map(map_info.id, spawn)
			return true
		end,
	},
	["좌표"] = {
		privilege = ROLE.Admin,
		usage = "- 현재 좌표·저장 스폰포인트 확인",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local x, y = me:position()
			local wz_t = m:wz()
			local map_id = 0
			local map_name = "?"
			if wz_t ~= nil then
				map_id = wz_t.id or 0
				map_name = tostring(wz_t.name or "?")
			end
			local spawn_id, spawn_name, sx, sy = me:spawn_point()
			if spawn_name ~= nil then
				me:notice(string.format(
					"Map: %s (%d), Position: %d, %d, Spawn: %d (%s) @ %d, %d",
					map_name, map_id, x, y, spawn_id, spawn_name, sx, sy
				))
			else
				me:notice(string.format(
					"Map: %s (%d), Position: %d, %d, Spawn: %d",
					map_name, map_id, x, y, spawn_id or 0
				))
			end
			return true
		end,
	},
	["위치"] = {
		privilege = ROLE.User,
		usage = "- 현재 맵·좌표를 채팅으로 표시",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:chat("맵 정보 없음")
				return true
			end
			local x, y = me:position()
			local wz_t = m:wz()
			local map_id = 0
			local map_name = "?"
			if wz_t ~= nil then
				map_id = wz_t.id or 0
				map_name = tostring(wz_t.name or "?")
			end
			me:chat(string.format("맵: %s (%d), 좌표: %d, %d", map_name, map_id, x, y))
			return true
		end,
	},
	["서버저장"] = {
		privilege = ROLE.Admin,
		usage = "- 온라인 유저 저장 요청",
		command = function(me, args)
			local ok, err = save()
			if ok then
				me:notice("서버 저장이 완료되었습니다.")
			else
				me:notice("서버 저장 실패: " .. tostring(err))
			end
			return true
		end,
	},
	["체력바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<체력값> - 체력 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /체력바꾸기 <체력값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then
				me:notice("체력은 0 이상이어야 합니다.")
				return true
			end
			if v > 32767 then v = 32767 end
			me:hp(v)
			me:base_hp(v)
			me:notice(string.format("체력 설정: %d", v))
			return true
		end,
	},
	["마력바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<마력값> - 마력 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /마력바꾸기 <마력값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then
				me:notice("마력은 0 이상이어야 합니다.")
				return true
			end
			if v > 32767 then v = 32767 end
			me:mp(v)
			me:base_mp(v)
			me:notice(string.format("마력 설정: %d", v))
			return true
		end,
	},
	["힘바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<힘값> - 힘 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /힘바꾸기 <힘값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then v = 0 end
			me:base_str(v)
			me:notice(string.format("힘 설정: %d", v))
			return true
		end,
	},
	["덱스바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<덱스값> - 덱스 설정",
		command = function(me, args)
			if not args[1] then

				me:notice("사용법: /덱스바꾸기 <덱스값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then v = 0 end
			me:base_dex(v)
			me:notice(string.format("덱스 설정: %d", v))
			return true
		end,
	},
	["인트바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<인트값> - 인트 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /인트바꾸기 <인트값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then v = 0 end
			me:base_int(v)
			me:notice(string.format("인트 설정: %d", v))
			return true
		end,
	},
	["럭바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<럭값> - 럭 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /럭바꾸기 <럭값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then v = 0 end
			me:base_luk(v)
			me:notice(string.format("럭 설정: %d", v))
			return true
		end,
	},
	["전능"] = {
		privilege = ROLE.Admin,
		usage = "<값> - 힘/덱/인트/럭 동일값 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /전능 <값>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then v = 0 end
			me:base_str(v)
			me:base_dex(v)
			me:base_int(v)
			me:base_luk(v)
			me:notice(string.format("전능 설정: %d", v))
			return true
		end,
	},
	["레벨바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<레벨> - 레벨 설정",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /레벨바꾸기 <레벨>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 1 then v = 1 end
			if v > 200 then v = 200 end
			me:level(v)
			me:notice(string.format("레벨 설정: %d", v))
			return true
		end,
	},
	["무적"] = {
		privilege = ROLE.Admin,
		usage = "- 무적 상태 토글",
		command = function(me, args)
			me:invincible(not me:invincible())
			local status = me:invincible() and "enabled" or "disabled"
			me:notice("무적 상태: " .. status)
			return true
		end,
	},
	["즉사"] = {
		privilege = ROLE.Admin,
		usage = "- 즉사 상태 토글",
		command = function(me, args)
			me:instant_kill(not me:instant_kill())
			local status = me:instant_kill() and "enabled" or "disabled"
			me:notice("즉사 상태: " .. status)
			return true
		end,
	},
	["직업바꾸기"] = {
		privilege = ROLE.Admin,
		usage = "<직업코드> - 직업 변경",
		command = function(me, args)
			if not args[1] then
				me:notice("사용법: /직업바꾸기 <직업코드>")
				return true
			end
			local v = tonumber(args[1])
			if not v or v < 0 then
				me:notice("직업코드는 0 이상이어야 합니다.")
				return true
			end
			me:class(v)
			me:notice(string.format("직업 변경: %d", v))
			return true
		end,
	},
	["몬스터생성"] = {
		privilege = ROLE.Admin,
		usage = "<몬스터ID/이름> [마리수] - 현재 위치에 몬스터 생성 (마리수 기본값 1)",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /몬스터생성 <몬스터ID/이름> [마리수]")
				return true
			end
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local x, y = me:position()
			local mob_wz = name2mob(args[1])
			local id = nil
			if mob_wz ~= nil then
				id = mob_wz.id
			else
				id = tonumber(args[1])
			end
			if id == nil then
				me:notice("잘못된 몬스터 ID 또는 이름: " .. args[1])
				return true
			end
			local count = 1
			if args[2] then
				count = tonumber(args[2]) or 1
			end
			if count < 1 then
				count = 1
			end
			local spawned = 0
			for _ = 1, count do
				local mob = m:spawn_mob(id, x, y)
				if mob ~= nil then
					spawned = spawned + 1
				end
			end
			if spawned == 0 then
				me:notice("몬스터 생성 실패")
			else
				me:notice(string.format("몬스터 생성: %s %d마리", args[1], spawned))
			end
			return true
		end,
	},
	["몬스터정보"] = {
		privilege = ROLE.Admin,
		usage = "- 현재 맵에 있는 몬스터 ID·이름 목록",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local entries = {}
			for _, mob in pairs(m:mobs()) do
				if mob:map() ~= nil and mob:hp() > 0 then
					entries[#entries + 1] = mob
				end
			end
			table.sort(entries, function(a, b)
				if a:id() ~= b:id() then
					return a:id() < b:id()
				end
				return a:oid() < b:oid()
			end)
			if #entries == 0 then
				me:notice("맵에 몬스터가 없습니다.")
				return true
			end
			me:notice(string.format("몬스터 %d마리", #entries))
			for _, mob in ipairs(entries) do
				local mob_id = mob:id()
				local mob_name = tostring(mob_id)
				local wz = id2mob(mob_id)
				if wz ~= nil and wz.name ~= nil and wz.name ~= "" then
					mob_name = wz.name
				end
				me:notice(string.format("  %d - %s", mob_id, mob_name))
			end
			return true
		end,
	},
	["몬스터죽이기"] = {
		privilege = ROLE.Admin,
		usage = "- 맵 몬스터 제거 (스펀지 본체 제외)",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local mobs = {}
			for _, mob in pairs(m:mobs()) do
				mobs[#mobs + 1] = mob
			end
			local count = 0
			local skipped_sponge = 0
			for _, mob in ipairs(mobs) do
				if mob:map() == nil then
					goto continue
				end
				if is_sponge_root(mob) then
					skipped_sponge = skipped_sponge + 1
					goto continue
				end
				local hp = mob:hp()
				if hp <= 0 then
					goto continue
				end
				if mob:damage(me, hp) then
					count = count + 1
				end
				::continue::
			end
			me:notice(string.format("몬스터 %d마리 제거", count))
			return true
		end,
	},
	["리액터초기화"] = {
		privilege = ROLE.Admin,
		usage = "- 현재 맵 리액터 초기화",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local count = m:reload_reactors()
			me:notice(string.format("리액터 %d개 초기화", count))
			return true
		end,
	},
	["몬스터체력"] = {
		privilege = ROLE.Admin,
		usage = "<체력값|체력%> - 맵 몬스터 체력 조정 (스펀지 본체 제외)",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /몬스터체력 <체력값|체력%>")
				return true
			end
			local raw = args[1]
			if args[2] == "%" then
				raw = raw .. "%"
			end
			local percent = false
			local num_str = raw
			if string.sub(raw, -1) == "%" then
				percent = true
				num_str = string.sub(raw, 1, -2)
				if num_str == "" then
					me:notice("사용법: /몬스터체력 <체력값|체력%>")
					return true
				end
			end
			local value = tonumber(num_str)
			if value == nil or value < 0 then
				me:notice("체력 값이 올바르지 않습니다.")
				return true
			end
			if percent and value > 100 then
				me:notice("체력 %는 0~100 사이여야 합니다.")
				return true
			end
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local changed = 0
			local skipped = 0
			for _, mob in pairs(m:mobs()) do
				if is_sponge_root(mob) then
					skipped = skipped + 1
					goto continue_mob_hp
				end
				local max_hp = mob:max_hp()
				if max_hp ~= nil and max_hp > 0 then
					local target_hp
					if percent then
						target_hp = math.floor(max_hp * value / 100)
					elseif value > max_hp then
						target_hp = max_hp
					else
						target_hp = math.floor(value)
					end
					local current = mob:hp()
					if current > target_hp then
						mob:damage(me, current - target_hp)
						changed = changed + 1
					elseif current < target_hp then
						mob:add_hp(target_hp - current)
						changed = changed + 1
					else
						skipped = skipped + 1
					end
				else
					skipped = skipped + 1
				end
				::continue_mob_hp::
			end
			if percent then
				me:notice(string.format("몬스터 체력 %.0f%% 적용 - %d마리 변경, %d마리 스킵", value, changed, skipped))
			else
				me:notice(string.format("몬스터 체력 %d 적용 - %d마리 변경, %d마리 스킵", value, changed, skipped))
			end
			return true
		end,
	},
	["스킬마스터"] = {
		privilege = ROLE.Admin,
		usage = "- 모든 배운 스킬의 레벨을 최대치로 올림",
		command = function(me, args)
			local class = me:class()
			if class < 100 then
				me:notice("직업이 2차 이상이어야 합니다.")
				return true
			end
			local wz_skills = class_learnable_skill_wzs(class)
			if wz_skills == nil then
				me:notice("스킬 정보를 가져올 수 없습니다.")
				return true
			end
			local added = 0
			local updated = 0
			for _, wz in pairs(wz_skills) do
				local skill = me:skill(wz.id)
				if skill == nil then
					skill = me:add_skill(wz.id)
					if skill ~= nil then
						added = added + 1
					end
				end
				if skill ~= nil then
					skill:level(wz.max_level, wz.master_level)
					updated = updated + 1
				end
			end
			me:notice(string.format("스킬마스터: %d 추가, %d 갱신", added, updated))
			return true
		end,
	},
	["스킬레벨"] = {
		privilege = ROLE.Admin,
		usage = "<스킬ID/이름> <레벨> [마스터레벨] - 스킬 레벨 설정",
		command = function(me, args)
			if #args < 2 then
				me:notice("사용법: /스킬레벨 <스킬ID/이름> <레벨> [마스터레벨]")
				return true
			end

			local master_level
			local level_arg_index = #args
			local target_end_index = #args - 1
			if #args >= 3 and tonumber(args[#args]) ~= nil then
				master_level = tonumber(args[#args])
				level_arg_index = #args - 1
				target_end_index = #args - 2
			end

			local level = tonumber(args[level_arg_index])
			if level == nil or level < 0 then
				me:notice("레벨은 0 이상의 숫자여야 합니다.")
				return true
			end

			local target = table.concat(args, " ", 1, target_end_index)
			if target == nil or target == "" then
				me:notice("사용법: /스킬레벨 <스킬ID/이름> <레벨> [마스터레벨]")
				return true
			end

			local skill_id = tonumber(target)
			local skill_wz = nil
			if skill_id == nil then
				skill_wz = name2skill(target)
				if skill_wz == nil then
					me:notice("존재하지 않는 스킬입니다: " .. target)
					return true
				end
				skill_id = skill_wz.id
			end

			if skill_wz ~= nil then
				if level > skill_wz.max_level then
					level = skill_wz.max_level
				end
			end

			if master_level == nil then
				if skill_wz ~= nil then
					master_level = skill_wz.master_level
				else
					master_level = level
				end
			end
			if master_level < level then
				master_level = level
			end
			if skill_wz ~= nil and master_level > skill_wz.master_level then
				master_level = skill_wz.master_level
			end

			local skill = me:skill(skill_id)
			if skill == nil then
				skill = me:add_skill(skill_id)
			end
			if skill == nil then
				me:notice(string.format("스킬 추가 실패: %d", skill_id))
				return true
			end

			skill:level(level, master_level)
			me:notice(string.format("스킬레벨 설정: %d -> level %d, master %d", skill_id, level, master_level))
			return true
		end,
	},
	["쿨타임초기화"] = {
		privilege = ROLE.Admin,
		usage = "[스킬ID] - 모든 스킬 또는 지정 스킬 쿨타임 초기화",
		command = function(me, args)
			if args[1] then
				local skill_id = tonumber(args[1])
				if not skill_id or skill_id <= 0 then
					me:notice("사용법: /쿨타임초기화 [스킬ID]")
					return true
				end
				local skill = me:skill(skill_id)
				if skill == nil then
					me:notice(string.format("배운 스킬이 아닙니다: %d", skill_id))
					return true
				end
				skill:cooldown(0)
				me:notice(string.format("스킬 %d 쿨타임 초기화 완료", skill_id))
				return true
			end

			local skills = me:skills()
			if skills == nil then
				me:notice("스킬 정보가 없습니다.")
				return true
			end
			local cleared = 0
			for _, skill in pairs(skills) do
				if skill ~= nil then
					skill:cooldown(0)
					cleared = cleared + 1
				end
			end
			me:notice(string.format("쿨타임 초기화 완료: %d개 스킬", cleared))
			return true
		end,
	},
	["엔피씨생성"] = {
		privilege = ROLE.Admin,
		usage = "<NPCID/이름> - 현재 위치에 NPC 생성",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /엔피씨생성 <NPCID/이름>")
				return true
			end
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local x, y = me:position()
			local id = name2npc(args[1])
			if id == nil then
				id = tonumber(args[1])
			end
			if id == nil then
				me:notice("잘못된 NPC ID 또는 이름: " .. args[1])
				return true
			end
			local npc = m:spawn_npc(id, x, y)
			if npc == nil then
				me:notice("NPC 생성 실패")
				return true
			end
			me:notice(string.format("NPC 생성: %s (OID %d)", args[1], npc:oid()))
			return true
		end,
	},
	["스크립트"] = {
		privilege = ROLE.Admin,
		usage = "- script/character/hook.lua의 on_script(me) 실행 (또는 me:script(경로, 함수, ...) 형태로 사용)",
		command = function(me, args)
			me:script("script/character/hook.lua", "on_script")
			return true
		end,
	},
	["맵리셋"] = {
		privilege = ROLE.Admin,
		usage = "- 현재 맵 리셋",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				return true
			end
			m:reset()
			return true
		end,
	},
	["패킷로그"] = {
		privilege = ROLE.Admin,
		usage = " [0|1|on|off] - 수신 패킷 로그 켜기/끄기",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				local enabled = get_packet_log()
				me:notice("packet log: " .. tostring(enabled))
				return true
			end
			local arg = string.lower(string.gsub(args[1], "^%s*(.-)%s*$", "%1"))
			local enable = (arg == "1" or arg == "on" or arg == "true")
			if not enable and arg ~= "0" and arg ~= "off" and arg ~= "false" then
				me:notice("사용법: 0|1|on|off")
				return true
			end
			set_packet_log(enable)
			me:notice("packet log set to " .. tostring(enable))
			return true
		end,
	},
	["이펙트테스트"] = {
		privilege = ROLE.Admin,
		usage = "<종류> [args] - 이펙트 패킷 테스트 (종류: skill, basic, dragon, hp, reward, maker, crafting, dice)",
		command = function(me, args)
			local kind = args[1]
			if kind == nil or kind == "" then
				me:notice("사용법: /이펙트테스트 <skill|basic|dragon|hp|reward|maker|crafting|dice> ...")
				return true
			end
			if kind == "skill" then
				local skill_arg = args[2]
				local effect_type = tonumber(args[3]) or SkillEffectType.Cast
				return effect_show_skill(me, skill_arg, effect_type)
			end
			if kind == "basic" then
				local effect_type = tonumber(args[2]) or EffectType.LevelUp
				return effect_show_basic(me, effect_type)
			end
			if kind == "dragon" then
				return effect_show_dragon_blood(me, args[2])
			end
			if kind == "hp" then
				return effect_show_hp_healed(me, args[2])
			end
			if kind == "reward" then
				return effect_show_reward_item_animation(me, args[2], args[3])
			end
			if kind == "maker" then
				return effect_show_item_maker_success(me)
			end
			if kind == "crafting" then
				return effect_show_crafting(me, args[2], args[3], args[4])
			end
			if kind == "dice" then
				return effect_show_dice(me, args[2], args[3])
			end
			me:notice("알 수 없는 종류입니다. skill, basic, dragon, hp, reward, maker, crafting, dice 중에서 선택하세요.")
			return true
		end,
	},
	["현재시간"] = {
		privilege = ROLE.Admin,
		usage = '[YYYY-MM-DD HH:MM:SS] - 현재 서버 시간 조회/설정',
		command = function(me, args)
			if #args == 0 then
				local dt = datetime()
				me:notice(string.format('현재 서버 시간: %04d-%02d-%02d %02d:%02d:%02d',
					dt.year, dt.month, dt.day, dt.hour, dt.minute, dt.second))
				return true
			end

			local value = table.concat(args, ' ')
			if not string.match(value, '^%d%d%d%d%-%d%d%-%d%d %d%d:%d%d:%d%d$') then
				me:notice('사용법: /현재시간 YYYY-MM-DD HH:MM:SS')
				return true
			end

			local success, error_message = now(value)
			if success then
				me:notice(string.format('현재 시간을 %s 로 설정 요청했습니다.', value))
			else
				me:notice(string.format('현재시간 설정 실패: %s', error_message or 'unknown error'))
			end
			return true
		end,
	},
	["현재시간초기화"] = {
		privilege = ROLE.Admin,
		usage = '- 현재 서버 시간 보정 초기화',
		command = function(me, args)
			local success, error_message = now('reset')
			if success then
				me:notice('현재 시간 보정 초기화 요청을 전송했습니다.')
			else
				me:notice(string.format('현재시간 보정 초기화 실패: %s', error_message or 'unknown error'))
			end
			return true
		end,
	},
	["시간가속"] = {
		privilege = ROLE.Admin,
		usage = '<timespan> - 시간 앞으로 이동 (예: 10.12:30:00 / 12:30:00)',
		command = function(me, args)
			local value = args[1]
			if not value then
				me:notice('사용법: /시간가속 <timespan>')
				return true
			end

			local success, error_message = time_forward(value)
			if success then
				me:notice(string.format('시간가속 적용 요청: %s', value))
			else
				me:notice(string.format('시간가속 실패: %s', error_message or 'unknown error'))
			end
			return true
		end,
	},
	["시간역전"] = {
		privilege = ROLE.Admin,
		usage = '<timespan> - 시간 뒤로 이동 (예: 10.12:30:00 / 12:30:00)',
		command = function(me, args)
			local value = args[1]
			if not value then
				me:notice('사용법: /시간역전 <timespan>')
				return true
			end

			local success, error_message = time_backward(value)
			if success then
				me:notice(string.format('시간역전 적용 요청: %s', value))
			else
				me:notice(string.format('시간역전 실패: %s', error_message or 'unknown error'))
			end
			return true
		end,
	},
}

function on_chat(me, message, shout)
	if string.sub(message, 1, 1) ~= "/" then
		return false
	end
	message = string.sub(message, 2, #message)
	local args = string_split(message, " ")
	local cmd = args[1]
	if cmd == nil or cmd == "" then
		return false
	end
	local cmd_data = command_funcs[cmd]
	if cmd_data == nil then
		return false
	end
	local cmd_func = nil
	local required_privilege = ROLE.User
	if type(cmd_data) == "table" then
		cmd_func = cmd_data.command
		required_privilege = cmd_data.privilege or ROLE.User
	else
		cmd_func = cmd_data
	end
	if me:role() < required_privilege then
		me:notice("권한이 부족합니다.")
		return true
	end
	table.remove(args, 1)
	local ok = cmd_func(me, args)
	return ok
end
