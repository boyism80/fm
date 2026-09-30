-- NPC name (String.wz/Npc.img.xml): 돈 지오바네

local HAIR_COUPON = 5150003
local COLOR_COUPON = 5151003

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "커닝시티 헤어샵에 온걸 환영하네! 하하하! 나는 이 헤어샵의 원장, 돈 지오바네 라고 한다네! #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 을 가져온다면 머리를 아주 아름~답게 손질해주지!", {
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
				styles = { 30030, 30020, 30000, 30360, 30370, 30450, 30350, 30190, 30180, 30160, 30130, 30110, 30050, 30040 }
			else
				styles = { 31050, 31040, 31000, 31010, 31510, 31170, 31180, 31330, 31140, 31130, 31120, 31090, 31060, 31020, 31010 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
			prompt = "어떤 머리 스타일을 원하나?"
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 2, base + 3, base + 7, base + 5 }
			prompt = "어떤 머리 색깔을 원하나?"
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
		me:dialog(npc, "머리 손질이 끝났다네! 후후. 정말 멋지군! 역시 내 솜씨는 훌륭하다니까!")
	end,
}
