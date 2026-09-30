-- NPC name (String.wz/Npc.img.xml): 미유

local HAIR_COUPON = 5150007
local COLOR_COUPON = 5151007

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "안녕하세요~ 루디브리엄 헤어숍입니다! 어머나~ 머리 스타일이 너무 촌스러우신걸요? 여기서 머리 손질을 좀 받으셔야겠어요~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 아이템을 가져와야 해요~", {
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
				styles = { 30030, 30020, 30000, 30510, 30340, 30710, 30300, 30050, 30160, 30190, 30280, 30240, 30150, 30650 }
			else
				styles = { 31040, 31050, 31000, 31520, 31460, 31290, 31280, 31270, 31230, 31160, 31120, 31150, 31010, 31030, 31650 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
			prompt = "원하시는 헤어 스타일을 선택해 주세요."
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 2, base + 3, base + 4, base + 5 }
			prompt = "헤어 컬러를 변경하시고 싶으신가요? 원하시는 컬러를 선택해주세요."
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
		me:dialog(npc, "와아~ 너무 멋져요! 원하시는 헤어 스타일이 또 생기실 경우 다시 찾아와 주세요~")
	end,
}
