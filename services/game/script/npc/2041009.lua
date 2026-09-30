-- NPC name (String.wz/Npc.img.xml): 미니

local HAIR_COUPON = 5150006
local COLOR_COUPON = 5151006

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "어서오세요~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k을 가져오시면 무작위로 머리 손질을 해 드린답니다. ", {
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
				styles = { 30030, 30020, 30000, 30510, 30340, 30710, 30300, 30050, 30160, 30190, 30280, 30240, 30150, 30650 }
			else
				styles = { 31040, 31050, 31000, 31520, 31460, 31290, 31280, 31270, 31230, 31160, 31120, 31150, 31010, 31030, 31650 }
			end
			for i = 1, #styles do
				styles[i] = styles[i] + color
			end
		else
			local base = math.floor(me:hair() / 10) * 10
			coupon = COLOR_COUPON
			styles = { base, base + 2, base + 3, base + 4, base + 5 }
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
