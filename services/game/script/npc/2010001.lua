-- NPC name (String.wz/Npc.img.xml): 미노

local HAIR_COUPON = 5150005
local COLOR_COUPON = 5151005
local MESSY_COUPON = 5154000
local MESSY_HAIR = 31240

return {
	on_click = function(me, npc)
		local sel = me:dialog_list(npc, "오르비스의 멋진 헤어숍 입니다. 여러분의 헤어 스타일을 책임져 드린답니다~ #b#i" .. HAIR_COUPON .. "# #t" .. HAIR_COUPON .. "##k 또는 #b#i" .. COLOR_COUPON .. "# #t" .. COLOR_COUPON .. "##k 아이템을 가져와야 해요~", {
			"머리 스타일 바꾸기",
			"머리 색깔 염색하기",
			"산발머리 쿠폰 사용하기",
		})
		if sel == nil then
			return
		end

		local coupon = nil
		local style = nil
		if sel == 3 then
			if math.floor(me:hair() / 10) * 10 == MESSY_HAIR then
				me:dialog(npc, "이미 머리 손질이 끝나신 것 같은데요~? 다른 헤어가 필요하시다면 다시 제게 말을 걸어주세요~")
				return
			end
			if me:gender() ~= 1 then
				me:dialog(npc, "죄송하지만 산발머리는 #r여성#k분만 이용이 가능해요. 다른 헤어를 원하신다면 다시 제게 말을 걸어주세요~")
				return
			end
			if me:dialog_yes_no(npc, "정말 #i5154000# #t5154000# 아이템을 사용하시겠어요?") == false then
				return
			end
			coupon = MESSY_COUPON
			style = MESSY_HAIR + me:hair() % 10
		else
			local styles = nil
			local prompt = nil
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
				coupon = HAIR_COUPON
				prompt = "어떤 헤어 스타일을 원하시나요?"
			else
				local base = math.floor(me:hair() / 10) * 10
				styles = { base, base + 1, base + 3, base + 4, base + 5 }
				coupon = COLOR_COUPON
				prompt = "어떤 헤어 색깔을 원하시나요?"
			end

			local pick = me:dialog_style(npc, prompt, styles)
			if pick == nil or styles[pick] == nil then
				return
			end
			style = styles[pick]
		end

		if me:exchange({ item = { [coupon] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "죄송하지만 쿠폰을 가져오시지 않으면 저희 헤어숍을 이용하실 수 없답니다.")
			return
		end
		me:hair(style)
		me:dialog(npc, "헤어 손질이 모두 끝났답니다. 다음에 또 이용해 주세요.")
	end,
}
