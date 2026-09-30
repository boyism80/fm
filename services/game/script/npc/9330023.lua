-- NPC name (String.wz/Npc.img.xml): 허치희

local COUPON = 5152024

return {
	on_click = function(me, npc)
		if me:dialog(npc, "시술을 받으러 여기까지 온건가요? #b#t" .. COUPON .. "##k 을 가져온다면 아주 멋진 얼굴로 바꿔드릴수 있답니다.", false, true) == false then
			return
		end

		local face = me:face()
		local styles = nil
		if me:gender() == 0 then
			styles = { 20000, 20001, 20002, 20003, 20004, 20005, 20006, 20007, 20008, 20012, 20014, 20016, 20020, 20017 }
		else
			styles = { 21000, 21001, 21002, 21003, 21004, 21005, 21006, 21007, 21008, 21012, 21014, 21016, 21020, 21017 }
		end
		for i = 1, #styles do
			styles[i] = styles[i] + face % 1000 - face % 100
		end

		local pick = me:dialog_style(npc, "어떤 얼굴로 바꿔보고 싶은가?", styles)
		if pick == nil then
			return
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.")
			return
		end
		me:face(styles[pick])
		me:dialog(npc, "자아~ 다 됐다네. 어떤가? 맘에 들지? 다음에 또 이곳을 찾아주게나~")
	end,
}
