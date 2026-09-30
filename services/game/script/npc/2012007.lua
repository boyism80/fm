-- NPC name (String.wz/Npc.img.xml): 린스

local HAIR_COUPON = 5150004
local COLOR_COUPON = 5151004

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "안녕하세요, 오르비스 헤어샵의 보조로 일하고 있는 린스 라고 해요. #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k을 가져오시면 무작위로 머리 손질을 해 드린답니다. ", {
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
				styles = { 30030, 30020, 30000, 30520, 30480, 30490, 30460, 30420, 30340, 30290, 30280, 30270, 30260, 30240, 30230 }
			else
				styles = { 31040, 31000, 31050, 31440, 31540, 31420, 31320, 31270, 31260, 31250, 31240, 31230, 31220, 31110, 31030, 31530 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 1, base + 7, base + 3, base + 4, base + 5 }
		end

		if me:dialog_yes_no(npc, "무작위로 머리 손질을 받으시고 싶으신가요? 원하시는것을 정확하게 선택하셨는지 확인해 주시기 바랍니다. 정말 무작위로 머리 손질을 해 드릴까요?") == false then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 머리 손질을 해드릴 수 없답니다.")
			return
		end
		me:hair(styles[math.random(1, #styles)])
		me:dialog(npc, "자~ 다 되었답니다. 원장님 못지 않은 솜씨죠? 마음에 드셨으면 좋겠어요.")
	end,
}
