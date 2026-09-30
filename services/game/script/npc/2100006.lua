-- NPC name (String.wz/Npc.img.xml): 마즈라

local HAIR_COUPON = 5150027
local COLOR_COUPON = 5151022

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "원하는 헤어라도 있는가~? 아리안트 미용실에 잘 오셨네~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 만 있다면 멋진 헤어스타일로 바꿔드리지~", {
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
				styles = { 30030, 30020, 30000, 30470, 30490, 30230, 30200, 30260, 30290, 30160, 30050, 30630 }
			else
				styles = { 31000, 31040, 31050, 31490, 31130, 31420, 31260, 31110, 31160, 31300, 31610, 31030 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
			prompt = "어떤 머리로 손질해주길 원하나~?"
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 3, base + 2, base + 4 }
			prompt = "어떤 머리 색으로 원하나?"
		end

		local pick = me:dialog_style(npc, prompt, styles)
		if pick == nil or styles[pick] == nil then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "미안하지만 쿠폰 없이는 머리를 손질해줄 수 없다네.")
			return
		end
		me:hair(styles[pick])
		me:dialog(npc, "머리 손질이 끝났다네!")
	end,
}
