-- NPC name (String.wz/Npc.img.xml): 문서뭉치

function on_click(me, npc)
	local q = me:quest(3240)
	if q == nil or not q:started() then
		return
	end
	local count = 0
	for _, it in pairs(me:item(4031034)) do
		count = count + it:count()
	end
	if count >= 1 then
		me:dialog(npc, "이미 #b#t4031034##k를 갖고 있는 것 같다.")
		return
	end
	if not me:dialog_yes_no(npc, "종이 뭉치 사이에 뭔가 주문서 처럼 생긴게 보인다. 가져갈까?") then
		return
	end
	local code = me:exchange(nil, { item = { [4031034] = 1 } })
	if code == ExchangeResult.LackCapacity then
		me:dialog(npc, "인벤토리 공간이 부족합니다.")
	end
end
