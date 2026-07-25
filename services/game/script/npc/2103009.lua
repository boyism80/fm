-- NPC name (String.wz/Npc.img.xml): 아리안트 민가1 찬장

return {
	on_click = function(me, npc)
		local npc_id = npc
		if type(npc) ~= "number" then
			npc_id = npc:id()
		end
		local q = me:quest(3926)
		if q == nil or not q:started() then
			return
		end
		local info = q:record()
		if info == nil or info == "" then
			info = "0000"
			q:record("0000")
			q:sync_progress()
		end
		if info == "3333" then
			return
		end
		local nCount = npc_id - 2103009
		local nInvite = info:sub(nCount + 1, nCount + 1)
		if nInvite == "3" then
			return
		end
		local count = 0
		for _, it in pairs(me:item(4031579)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "내려놓을 보물이 없다.")
			return
		end
		if me:exchange({ item = { [4031579] = 1 } }, nil) ~= ExchangeResult.OK then
			return
		end
		me:dialog(npc, "가져온 보물을 살며시 내려놓았다.")
		local z = ""
		for i = 0, 3 do
			if nCount == i then
				z = z .. "3"
			else
				z = z .. info:sub(i + 1, i + 1)
			end
		end
		q:record(z)
		q:sync_progress()
		if z == "3333" then
			me:show_quest_completion(3926)
		end
	end
}
