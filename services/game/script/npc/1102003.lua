-- NPC name (String.wz/Npc.img.xml): 키단

local NEXT_CLASSES = {
	[0] = {
		{ 1100, "#r소울마스터 : 1차 #b전직하기#k" },
		{ 1200, "#r플레임위자드 : 1차 #b전직하기#k" },
		{ 1300, "#r윈드브레이커 : 1차 #b전직하기#k" },
		{ 1400, "#r나이트워커 : 1차 #b전직하기#k" },
		{ 1500, "#r스트라이커 : 1차 #b전직하기#k" },
	},
	[1100] = { { 1110, "#r소울마스터 : 2차 #b전직하기#k" } },
	[1200] = { { 1210, "#r플레임위자드 : 2차 #b전직하기#k" } },
	[1300] = { { 1310, "#r윈드브레이커 : 2차 #b전직하기#k" } },
	[1400] = { { 1410, "#r나이트워커 : 2차 #b전직하기#k" } },
	[1500] = { { 1510, "#r스트라이커 : 2차 #b전직하기#k" } },
	[1110] = { { 1111, "#r소울마스터 : 3차 #b전직하기#k" } },
	[1210] = { { 1211, "#r플레임위자드 : 3차 #b전직하기#k" } },
	[1310] = { { 1311, "#r윈드브레이커 : 3차 #b전직하기#k" } },
	[1410] = { { 1411, "#r나이트워커 : 3차 #b전직하기#k" } },
	[1510] = { { 1511, "#r스트라이커 #b전직하기#k" } },
	[1111] = { { 1112, "#r소울마스터 : 4차 #b전직하기#k" } },
	[1211] = { { 1212, "#r플레임위자드 : 4차 #b전직하기#k" } },
	[1311] = { { 1312, "#r윈드브레이커 : 4차 #b전직하기#k" } },
	[1411] = { { 1412, "#r나이트워커 : 4차 #b전직하기#k" } },
	[1511] = { { 1512, "#r스트라이커 : 4차 #b전직하기#k" } },
}

return {
	on_click = function(me, npc)
		local class = me:class()
		local level = me:level()
		local message = "#b#h ##k님의 전직을 도와드리고 있습니다.\r\n"
		if class == Class.Beginner then
			message = message .. "1차 전직은 10 부터 전직이 가능합니다.\r\n"
		elseif class % 100 == 0 then
			message = message .. "2차 전직은 30 부터 전직이 가능합니다.\r\n"
		elseif class % 10 == 0 then
			message = message .. "3차 전직은 70 부터 전직이 가능합니다.\r\n"
		elseif class % 10 ~= 2 then
			message = message .. "4차 전직은 120 부터 전직이 가능합니다.\r\n"
		end

		local next_classes = NEXT_CLASSES[class]
		if next_classes == nil then
			me:dialog(npc, message .. "#r하지만 #b#h ##k님은 이미 모든 전직을 하셨습니다.#k")
			return
		end
		local options = {}
		for i, entry in ipairs(next_classes) do
			options[i] = entry[2]
		end
		local sel = me:dialog_list(npc, message, options)
		if sel == nil then
			return
		end

		local required = 120
		if class == Class.Beginner then
			required = 10
		elseif class % 100 == 0 then
			required = 30
		elseif class % 10 == 0 then
			required = 70
		end
		if level < required then
			me:dialog(npc, "전직할 레벨이 부족합니다.")
			return
		end
		me:class(next_classes[sel][1])
		if class == Class.Beginner then
			me:exchange(nil, { meso = 3000 })
		end
	end
}
