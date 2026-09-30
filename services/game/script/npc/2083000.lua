-- NPC name (String.wz/Npc.img.xml): 결사대 암호석판

local PASS = 4001086
local DEST = 240050400

return {
	on_click = function(me, npc)
		local party = me:party()
		if party == nil then
			me:dialog(npc, "파티를 이루고 있지 않아 암호 석판을 읽을 수 없습니다.")
			return
		end
		local count = 0
		for _, it in pairs(me:item(PASS)) do
			count = count + it:count()
		end
		if count < 1 then
			me:dialog(npc, "암호석판을 읽어보려 하지만, 무슨 문자가 적혀있는지 알 수 없습니다. 혼테일 결사대원이라면 읽을 수 있을 것 같습니다.")
			return
		end
		if me:dialog_yes_no(npc, "암호석판이 빛나더니 석판 뒤로 문이 열렸습니다. 문을 이용해서 입장하시겠습니까?") == false then
			return
		end
		party = me:party()
		if party == nil then
			return
		end
		local pid = party:id()
		for _, ch in pairs(me:map():characters()) do
			local p = ch:party()
			if p ~= nil and p:id() == pid then
				ch:map(DEST, 0)
			end
		end
	end
}
