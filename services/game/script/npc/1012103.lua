-- NPC name (String.wz/Npc.img.xml): 나탈리

local HAIR_COUPON = 5150001
local COLOR_COUPON = 5151001

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "어디가도 멋진 헤어스타일을 만들 곳은 여기 뿐이란 말이지. 서비스를 이용하려면 #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 아이템을 가져와야 해요~", {
			"머리 스타일 바꾸기",
			"머리 색깔 염색하기",
		})
		if sel == nil then
			return
		end

		local coupon = HAIR_COUPON
		local styles = nil
		local prompt = nil
		if sel == 1 then
			local color = me:hair() % 10
			if me:gender() == 0 then
				styles = { 30030, 30020, 30000, 30480, 30440, 30470, 30410, 30330, 30310, 30210, 30200, 30150, 30140, 30120, 30060, 30670 }
			else
				styles = { 31050, 31040, 31000, 31070, 31080, 31030, 31340, 31480, 31490, 31410, 31310, 31300, 31160, 31100, 31150, 31640 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
			prompt = "오호호~ 어떤 헤어스타일을 원하시나요 고객님?"
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base + 1, base + 2, base + 4, base + 7 }
			prompt = "머리를 염색하고 싶으시군요? 원하시는 색을 골라보세요."
		end

		local pick = me:dialog_style(npc, prompt, styles)
		if pick == nil then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 머리 손질을 해드릴 수 없답니다.")
			return
		end
		me:hair(styles[pick])
		me:dialog(npc, "자~ 다 되었답니다. 어떠세요? 저희 미용실만의 최고의 솜씨를 발휘해 보았답니다.")
	end,
}
