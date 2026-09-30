-- NPC name (String.wz/Npc.img.xml): 브리트니

local HAIR_COUPON = 5150000
local COLOR_COUPON = 5151000

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "헤네시스 헤어샵의 보조를 맡고 있는 브리트니 라고 합니다. #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k을 가져오시면 무작위로 머리 손질을 해 드린답니다. ", {
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
				styles = { 30030, 30020, 30000, 30480, 30440, 30470, 30410, 30330, 30310, 30210, 30200, 30150, 30140, 30120, 30060, 30670 }
			else
				styles = { 31050, 31040, 31000, 31070, 31080, 31030, 31340, 31480, 31490, 31410, 31310, 31300, 31160, 31100, 31150, 31640 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base + 1, base + 2, base + 4, base + 7 }
		end

		if me:dialog_yes_no(npc, "무작위로 머리 손질을 받으시고 싶으신가요? 원하시는것을 정확하게 선택하셨는지 확인해 주시기 바랍니다. 정말 무작위로 머리 손질을 해 드릴까요?") == false then
			return
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 머리 손질을 해드릴 수 없답니다.")
			return
		end
		me:hair(styles[math.random(1, #styles)])
		me:dialog(npc, "자~ 다 되었답니다. 어떠세요? 마음에 드셨으면 좋겠군요.")
	end,
}
