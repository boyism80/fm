-- NPC name (String.wz/Npc.img.xml): 첫번째 마법진

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

return {
	on_click = function(me, npc)
		local q = me:quest(3345)
		if q == nil or not q:started() then
			return
		end
		local info = q:record()
		if info == nil or info == "" then
			info = "0"
		end
		local count = 0
		for _, it in pairs(me:item(4031739)) do
			count = count + it:count()
		end
		if info == "0" and count >= 1 then
			q:record("1")
			q:sync_progress()
			show_npc_effect(me, npc, "act33451")
			local map = me:map()
			if map ~= nil then
				map:message("첫번째 마법진이 반응했습니다.")
			end
			me:exchange({ item = { [4031739] = 1 } }, nil)
		end
	end
}
