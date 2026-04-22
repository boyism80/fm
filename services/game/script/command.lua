


local function string_split(s, sep)
	sep = sep or " "
	local t = {}
	local pattern = "([^" .. sep .. "]+)"
	for v in string.gmatch(s, pattern) do
		table.insert(t, v)
	end
	return t
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
	["메소초기화"] = {
		privilege = ROLE.Admin,
		usage = "- 메소 초기화",
		command = function(me, args)
			me:meso(0)
			me:notice("메소를 초기화했습니다.")
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
		usage = "<맵이름> [스폰포인트] - 맵 이동",
		command = function(me, args)
			if not args[1] or args[1] == "" then
				me:notice("사용법: /맵이동 <맵이름> [스폰포인트]")
				return true
			end
			local spawn = 1
			if args[2] then
				spawn = tonumber(args[2]) or 1
			end
			if name2map(args[1]) == nil then
				me:notice("존재하지 않는 맵입니다: " .. args[1])
				return true
			end
			me:map(args[1], spawn)
			return true
		end,
	},
	["좌표"] = {
		privilege = ROLE.Admin,
		usage = "- 현재 좌표 확인",
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
			me:notice(string.format("Map: %s (%d), Position: %d, %d", map_name, map_id, x, y))
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
			local id = name2mob(args[1])
			if id == nil then
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
	["몬스터죽이기"] = {
		privilege = ROLE.Admin,
		usage = "- 모든 몬스터 제거",
		command = function(me, args)
			local m = me:map()
			if m == nil then
				me:notice("맵 정보 없음")
				return true
			end
			local count = 0
			for oid, _ in pairs(m:mobs()) do
				m:remove_mob(oid, MobDieAnimation.FadeOut)
				count = count + 1
			end
			me:notice(string.format("몬스터 %d마리 제거", count))
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
		usage = "- script/script.lua의 on_script(me) 실행 (또는 me:script(경로, 함수, ...) 형태로 사용)",
		command = function(me, args)
			me:script("script/script.lua", "on_script")
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
