package entity

type Life struct {
	Object
	Hp         uint16
	MaxHp      uint16
	Mp         uint16
	MaxMp      uint16
	Stance     uint8
	Invincible bool
}
