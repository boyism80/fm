-- NPC name (String.wz/Npc.img.xml): 니키

local HAIR_COUPON = 5150022
local COLOR_COUPON = 5151015

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "하하! 야시장의 멋진 헤어샵에 잘 오셨습니다~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 만 있다면 멋진 헤어스타일로 바꿔드리지~", {
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
				styles = { 30030, 30020, 30000, 30560, 30420, 30310, 30270, 30260, 30240, 30230, 30200, 30150 }
			else
				styles = { 31000, 31040, 31050, 31560, 31320, 31240, 31230, 31150, 31140, 31110, 31100, 31030, 31010 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
			prompt = "어떤 머리로 손질해주길 원하나요~?"
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 3, base + 4, base + 7, base + 5 }
			prompt = "어떤 머리 색으로 원하나요~?"
		end

		local pick = me:dialog_style(npc, prompt, styles)
		if pick == nil or styles[pick] == nil then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없습니다.")
			return
		end
		me:hair(styles[pick])
		me:dialog(npc, "머리 손질이 끝났답니다~")
	end,
}
