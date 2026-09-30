-- NPC name (String.wz/Npc.img.xml): 클리프

return {
	on_click = function(me, npc)
		local pages = {
			{ "Do you see a bunch of snowmen standing around there? Go talk to one of them, and it'll take you to the famous Christmas tree here that is just humongous. The tree can be decorated using various kinds of ornaments. What do you think? Sounds fun, right?", false, true },
			{ "Only 6 can be at the place where the tree is at once, and you can't #btrade or open store#k there. The ornaments that you drop can only be picked back up by yourself, so don't worry about losing your ornaments here.", true, true },
			{ "Of course, the items that are dropped in there will never disappear. Once you get out of there through the snowman that's inside, all the items you've dropped at that map will come back to you, so you won't have to pick all those items up before leaving the place. Isn't that sweet?", true, true },
			{ "Well then, go see #p2002001#, buy some Christmas ornaments there, and then decorate the tree with those~ Oh yeah! The biggest, the most beautiful ornament cannot be bought from him. It's probably ... taken by a monster ... huh huh ..", true, false },
		}
		local i = 1
		while i >= 1 and i <= #pages do
			local page = pages[i]
			local forward = me:dialog(npc, page[1], page[2], page[3])
			if forward then
				i = i + 1
			elseif i == 1 then
				return
			else
				i = i - 1
			end
		end
	end
}
