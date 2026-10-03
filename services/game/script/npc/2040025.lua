-- NPC name (String.wz/Npc.img.xml): 두번째 에오스의 돌

local SCROLL = 4001020

return {
	on_click = function(me, npc)
		local npc_id = npc:id()

		local count = 0
		for _, it in pairs(me:item(SCROLL)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "#p" .. npc_id .. "#이다. 하지만 #b#t4001020##k가 없어 활성화 시킬 수 없을 것 같다.")
			return
		end

		local sel = me:dialog_list(npc, "#b#t4001020##k로 #p" .. npc_id .. "#을 활성화 시켰습니다. 어디로 가시겠습니까?\r\n#b\r\n", {
			"에오스 탑 100층",
			"에오스 탑 41층",
		})
		if sel == nil then
			return
		end

		local map_id = 221024400
		local text = "#b#t4001020##k를 사용하여 #p2040024##k이 있는 100층으로 이동하시겠습니까?"
		if sel == 2 then
			map_id = 221021700
			text = "#b#t4001020##k를 사용하여 #p2040026##k이 있는 41층으로 이동하시겠습니까?"
		end
		if not me:dialog_yes_no(npc, text) then
			return
		end
		local code = me:exchange({ item = { [SCROLL] = 1 } }, nil)
		if code ~= ExchangeResult.OK then
			return
		end
		me:map(map_id, 3)
	end
}
