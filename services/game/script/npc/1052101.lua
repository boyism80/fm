-- NPC name (String.wz/Npc.img.xml): 안드레아

local HAIR_COUPON = 5150002
local COLOR_COUPON = 5151002

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "이몸은~ 세상 최고의 헤어 디자이너를 꿈꾸는~ 안드레아 라고 하지~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k을 가져오시면 무작위로 머리 손질을 해 준다구~ ", {
			"무작위로 머리 스타일 바꾸기",
			"무작위로 머리 색깔 염색하기",
		})
		if sel == nil then
			return
		end

		local coupon = HAIR_COUPON
		local styles = nil
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
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 1, base + 2, base + 4, base + 6, base + 7 }
		end

		if me:dialog_yes_no(npc, "정말로 무작위로 머리 손질을 받고 싶어~? 무슨 결과가 나와도 장담은 못한다구. 뭐 이 안드레아 님의 솜씨라면 뭐든지 아름다울 테지만~") == false then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "쿠폰이 없으면 아무리 이 몸이 관대하다고 해도 머리 손질을 해줄 수 없다구~")
			return
		end
		me:hair(styles[math.random(1, #styles)])
		me:dialog(npc, "이 몸의 판타스틱 하고 엘레강스한 ~ 솜씨가 어때? 반해버리겠지?")
	end,
}
