-- Item name (String.wz/Consume.img.xml): 두손도끼 공격력 주문서 70%

function on_scroll(me, scroll_slot, target_slot)
    return me:enhance(scroll_slot, target_slot, EquipmentPart.Weapon, WeaponType.Axe2H, function(target, scroll)
        target:add_bonus_stats(scroll:wz():bonus_stats())
    end)
end
