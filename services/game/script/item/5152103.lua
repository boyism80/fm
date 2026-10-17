-- Item name (String.wz/Cash.img.xml): 1회용 컬러렌즈(초록색)

return {
	on_cash = function(me, item_id)
		local face = me:face()
		me:face(face - face % 1000 + face % 100 + item_id % 10 * 100)
		return true
	end,
}
