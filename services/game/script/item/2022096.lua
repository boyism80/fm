-- Item name (String.wz/Consume.img.xml): 닭튀김

return {
	on_active_item = function(me, item)
		me:exp(me:exp() + 10000)
		return false
	end,
}
