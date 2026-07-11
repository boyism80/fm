-- NPC name (String.wz/Npc.img.xml): 세번째 마법진

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
	local count = 0
	for _, it in pairs(me:item(4031741)) do
		count = count + it:count()
	end
	if info == "2" and count >= 1 then
		q:record("3")
		q:sync_progress()
		show_npc_effect(me, npc, "act33453")
		local map = me:map()
		if map ~= nil then
			map:message("세번째 마법진이 반응했습니다.")
		end
		me:exchange({ item = { [4031741] = 1 } }, nil)
	end
end
