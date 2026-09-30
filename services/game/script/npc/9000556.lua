-- NPC name (String.wz/Npc.img.xml): 고로

local BLOCKED = {
	[1122017] = true,
	[1122307] = true,
}
local MAX_DROP = 30000

return {
	on_click = function(me, npc)
		local slots = {}
		for slot in pairs(me:items(InventoryType.Cash)) do
			table.insert(slots, slot)
		end
		table.sort(slots)
		if #slots == 0 then
			me:dialog(npc, "버릴 캐시아이템을 갖고 계시지 않으셔요.")
			return
		end

		local options = {}
		for _, slot in ipairs(slots) do
			local id = me:item(InventoryType.Cash, slot):wz():id()
			table.insert(options, "#b#v" .. id .. "# #t" .. id .. "##k (코드 : " .. id .. ")")
		end
		local sel = me:dialog_list(npc, "	버릴 아이템을 선택해주세요.\r\n", options)
		if sel == nil then
			return
		end

		local slot = slots[sel]
		local item = me:item(InventoryType.Cash, slot)
		if item == nil then
			me:dialog(npc, "Error, please try again.")
			return
		end
		local id = item:wz():id()
		if BLOCKED[id] then
			me:dialog(npc, "#i" .. id .. "##b#z" .. id .. "##k은(는) 버릴 수 없어요.")
			return
		end

		local limit = math.min(item:count(), MAX_DROP)
		local count = tonumber(me:dialog_input(npc, "#i" .. id .. ":# #b#z" .. id .. "##k\r\n버릴 개수를 적어주세요."))
		if count == nil or count < 1 or count > limit then
			return
		end
		if me:rmitem(InventoryType.Cash, slot, math.floor(count)) == false then
			me:dialog(npc, "Error, please try again!")
		end
	end
}
