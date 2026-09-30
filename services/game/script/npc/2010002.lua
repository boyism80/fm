-- NPC name (String.wz/Npc.img.xml): 프란츠

local COUPON = 5152005

return {
	on_click = function(me, npc)
		if me:dialog(npc, "어서오시게. 오르비스 성형외과에 온걸 환영하네. #b#t" .. COUPON .. "##k 을 가져온다면 자네의 얼굴을 멋지게 바꿔주겠네.", false, true) == false then
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
		if pick == nil or styles[pick] == nil then
			return
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.")
			return
		end
		me:face(styles[pick])
		me:dialog(npc, "자아~ 다 됐다네. 어떤가? 멋지지 않은가? 하하하. 나중에 성형이 또 하고 싶다면 이곳을 다시 찾아주게나!")
	end,
}
