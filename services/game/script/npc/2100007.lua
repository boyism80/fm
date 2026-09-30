-- NPC name (String.wz/Npc.img.xml): 라일라

local COUPON = 5153007

return {
	on_click = function(me, npc)
		if me:dialog(npc, "아름다운 피부로 가꾸시고 싶으신가요? 제가 도와드리도록 하지요. 서비스를 이용하려면 #b#t" .. COUPON .. "##k 아이템을 가져오셔야 합니다.", false, true) == false then
			return
		end

		local styles = { 0, 1, 2, 3, 4 }
		local pick = me:dialog_style(npc, "원하시는 피부를 골라보세요~", styles)
		if pick == nil or styles[pick] == nil then
			return
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 피부관리를 해드릴 수 없답니다.")
			return
		end
		me:skin(styles[pick])
		me:dialog(npc, "자~ 다 되었답니다. 아름다운 피부는 미의 기본.. 앞으로도 아름다운 피부를 가꾸시고 싶으시다면 저희 피부관리실을 찾아주세요.")
	end,
}
