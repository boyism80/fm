-- NPC name (String.wz/Npc.img.xml): 바드르

local COUPON = 5152030

return {
	on_click = function(me, npc)
		if me:dialog(npc, "성형수술이 필요한가? #b#t" .. COUPON .. "##k 을 가져온다면 자네의 얼굴을 바꿔주겠네.", false, true) == false then
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

		local pick = me:dialog_style(npc, "원하는 얼굴을 골라보게.", styles)
		if pick == nil or styles[pick] == nil then
			return
		end

		if me:exchange({ item = { [COUPON] = 1 } }) ~= ExchangeResult.OK then
			me:dialog(npc, "미안하지만 쿠폰 없이는 성형수술을 해 줄수 없다네.")
			return
		end
		me:face(styles[pick])
		me:dialog(npc, "자, 다 되었다네.")
	end,
}
