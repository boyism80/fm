-- NPC name (String.wz/Npc.img.xml): 마법진 중앙

local function show_npc_effect(me, npc_id, effect)
	local map = me:map()
	if map == nil then
		return
	end
	for _, n in pairs(map:npcs()) do
		if n:id() == npc_id then
			n:show_effect(effect)
			return
		end
	end
end

function on_click(me, npc)
	local q = me:quest(3345)
	if q == nil or not q:started() then
		return
	end
	local info = q:record()
	if info == nil or info == "" then
		info = "0"
	end
	if info == "3" then
		q:record("4")
		q:sync_progress()
		me:show_quest_completion(3345)
		show_npc_effect(me, npc, "act33454")
		local map = me:map()
		if map ~= nil then
			map:message("마법진이 빛을 발하기 시작합니다.")
		end
	end
end
