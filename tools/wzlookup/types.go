package main

type Meta struct {
	Generator string `yaml:"generator"`
	WzPath    string `yaml:"wz_path"`
}

type QuestsFile struct {
	Meta   Meta             `yaml:"_meta"`
	Quests map[uint32]Quest `yaml:"quests"`
}

type Quest struct {
	ID       uint32     `yaml:"id"`
	Meta     QuestMeta  `yaml:"meta"`
	Start    QuestPhase `yaml:"start"`
	Complete QuestPhase `yaml:"complete"`
}

type QuestMeta struct {
	Name            string         `yaml:"name,omitempty"`
	Parent          string         `yaml:"parent,omitempty"`
	Order           int            `yaml:"order,omitempty"`
	Descriptions    map[int]string `yaml:"descriptions,omitempty"`
	Area            int            `yaml:"area,omitempty"`
	AutoStart       bool           `yaml:"auto_start,omitempty"`
	AutoPreComplete bool           `yaml:"auto_pre_complete,omitempty"`
	AutoComplete    bool           `yaml:"auto_complete,omitempty"`
	AutoAccept      bool           `yaml:"auto_accept,omitempty"`
	Blocked         bool           `yaml:"blocked,omitempty"`
	ViewMedalItem   int            `yaml:"view_medal_item,omitempty"`
	SelectedSkillID int            `yaml:"selected_skill_id,omitempty"`
	TimeLimit       int            `yaml:"time_limit,omitempty"`
	TimeLimit2      int            `yaml:"time_limit2,omitempty"`
}

type QuestPhase struct {
	Requirements []QuestRequirement `yaml:"requirements,omitempty"`
	Actions      []QuestAction      `yaml:"actions,omitempty"`
}

type QuestRequirement struct {
	Kind        string         `yaml:"kind"`
	IntValue    int            `yaml:"int_value,omitempty"`
	StrValue    string         `yaml:"str_value,omitempty"`
	InfoStrings []string       `yaml:"info_strings,omitempty"`
	Classes     []int          `yaml:"jobs,omitempty"`
	PetIDs      []uint32       `yaml:"pet_ids,omitempty"`
	Items       map[uint32]int `yaml:"items,omitempty"`
	Mobs        map[uint32]int `yaml:"mobs,omitempty"`
	Quests      map[uint32]int `yaml:"quests,omitempty"`
	Skills      map[uint32]int `yaml:"skills,omitempty"`
}

type QuestAction struct {
	Kind              string             `yaml:"kind"`
	IntValue          int                `yaml:"int_value,omitempty"`
	StrValue          string             `yaml:"str_value,omitempty"`
	ApplicableClasses []int              `yaml:"applicable_jobs,omitempty"`
	Items             []QuestRewardItem  `yaml:"items,omitempty"`
	Skills            []QuestRewardSkill `yaml:"skills,omitempty"`
	Quests            map[uint32]int     `yaml:"quests,omitempty"`
}

type QuestRewardItem struct {
	ItemID     uint32 `yaml:"item_id"`
	Count      int    `yaml:"count"`
	Class      int    `yaml:"job,omitempty"`
	ClassEx    int    `yaml:"job_ex,omitempty"`
	Gender     int    `yaml:"gender,omitempty"`
	Period     int    `yaml:"period,omitempty"`
	Prop       int    `yaml:"prop,omitempty"`
	DateExpire string `yaml:"date_expire,omitempty"`
}

type QuestRewardSkill struct {
	SkillID     uint32 `yaml:"skill_id"`
	SkillLevel  int    `yaml:"skill_level"`
	MasterLevel int    `yaml:"master_level"`
	Classes     []int  `yaml:"jobs,omitempty"`
}

type MapsFile struct {
	Meta Meta           `yaml:"_meta"`
	Maps map[uint32]Map `yaml:"maps"`
}

type Map struct {
	ID            uint32                  `yaml:"id"`
	Name          string                  `yaml:"name,omitempty"`
	StreetName    string                  `yaml:"street_name,omitempty"`
	MapName       string                  `yaml:"map_name,omitempty"`
	ReturnMapID   int                     `yaml:"return_map_id,omitempty"`
	ForcedReturn  int                     `yaml:"forced_return,omitempty"`
	FieldLimit    int                     `yaml:"field_limit,omitempty"`
	IsTown        bool                    `yaml:"is_town,omitempty"`
	BGM           string                  `yaml:"bgm,omitempty"`
	MapMark       string                  `yaml:"map_mark,omitempty"`
	MapDesc       string                  `yaml:"map_desc,omitempty"`
	NpcSpawns     map[uint32]NpcSpawn     `yaml:"npc_spawns,omitempty"`
	MobSpawns     map[uint32]MobSpawn     `yaml:"mob_spawns,omitempty"`
	ReactorSpawns map[uint32]ReactorSpawn `yaml:"reactor_spawns,omitempty"`
	Portals       map[uint8]Portal        `yaml:"portals,omitempty"`
}

type NpcSpawn struct {
	SpawnID    uint32 `yaml:"spawn_id"`
	NpcID      uint32 `yaml:"npc_id"`
	NpcName    string `yaml:"npc_name,omitempty"`
	X          int16  `yaml:"x"`
	Y          int16  `yaml:"y"`
	Foothold   int16  `yaml:"foothold,omitempty"`
	Hide       bool   `yaml:"hide,omitempty"`
	MobTimeSec int    `yaml:"mob_time_sec,omitempty"`
}

type MobSpawn struct {
	SpawnID    uint32 `yaml:"spawn_id"`
	MobID      uint32 `yaml:"mob_id"`
	MobName    string `yaml:"mob_name,omitempty"`
	X          int16  `yaml:"x"`
	Y          int16  `yaml:"y"`
	Foothold   int16  `yaml:"foothold,omitempty"`
	Hide       bool   `yaml:"hide,omitempty"`
	MobTimeSec int    `yaml:"mob_time_sec,omitempty"`
}

type ReactorSpawn struct {
	SpawnID    uint32 `yaml:"spawn_id"`
	ReactorID  uint32 `yaml:"reactor_id"`
	X          int16  `yaml:"x"`
	Y          int16  `yaml:"y"`
	RespawnSec int    `yaml:"respawn_sec,omitempty"`
	Name       string `yaml:"name,omitempty"`
}

type Portal struct {
	ID          uint8  `yaml:"id"`
	Name        string `yaml:"name,omitempty"`
	TargetMapID int32  `yaml:"target_map_id,omitempty"`
	Target      string `yaml:"target,omitempty"`
	X           int16  `yaml:"x"`
	Y           int16  `yaml:"y"`
	ScriptName  string `yaml:"script_name,omitempty"`
	Type        uint8  `yaml:"type,omitempty"`
}

type MobsFile struct {
	Meta Meta           `yaml:"_meta"`
	Mobs map[uint32]Mob `yaml:"mobs"`
}

type Mob struct {
	ID                    uint32         `yaml:"id"`
	Name                  string         `yaml:"name,omitempty"`
	BodyAttack            int            `yaml:"body_attack,omitempty"`
	Level                 uint8          `yaml:"level,omitempty"`
	MaxHP                 int            `yaml:"max_hp,omitempty"`
	MaxMP                 int            `yaml:"max_mp,omitempty"`
	Speed                 int16          `yaml:"speed,omitempty"`
	PADamage              int            `yaml:"pa_damage,omitempty"`
	PDDamage              int            `yaml:"pd_damage,omitempty"`
	MADamage              int            `yaml:"ma_damage,omitempty"`
	MDDamage              int            `yaml:"md_damage,omitempty"`
	ACC                   int            `yaml:"acc,omitempty"`
	EVA                   int            `yaml:"eva,omitempty"`
	EXP                   uint32         `yaml:"exp,omitempty"`
	Undead                bool           `yaml:"undead,omitempty"`
	Pushed                bool           `yaml:"pushed,omitempty"`
	Boss                  bool           `yaml:"boss,omitempty"`
	FfaLoot               bool           `yaml:"ffa_loot,omitempty"`
	ExplosiveReward       bool           `yaml:"explosive_reward,omitempty"`
	FS                    float32        `yaml:"fs,omitempty"`
	SummonType            uint8          `yaml:"summon_type,omitempty"`
	MobType               uint8          `yaml:"mob_type,omitempty"`
	Link                  string         `yaml:"link,omitempty"`
	ElemResist            map[string]int `yaml:"elem_resist,omitempty"`
	Skills                []MobSkillSlot `yaml:"skills,omitempty"`
	Attacks               []MobAttack    `yaml:"attacks,omitempty"`
	Banish                *MobBanish     `yaml:"banish,omitempty"`
	Revives               []uint32       `yaml:"revives,omitempty"`
	RemoveAfter           int            `yaml:"remove_after,omitempty"`
	SelfDestructionAction int8           `yaml:"self_destruction_action,omitempty"`
	HpTagColor            uint8          `yaml:"hp_tag_color,omitempty"`
	HpTagBgColor          uint8          `yaml:"hp_tag_bg_color,omitempty"`
}

type MobSkillSlot struct {
	Slot    int    `yaml:"slot"`
	SkillID uint32 `yaml:"skill_id"`
	Level   uint8  `yaml:"level"`
	Action  int    `yaml:"action,omitempty"`
}

type MobAttack struct {
	Index        int    `yaml:"index"`
	DeadlyAttack bool   `yaml:"deadly_attack,omitempty"`
	MpBurn       int    `yaml:"mp_burn,omitempty"`
	MpCon        int    `yaml:"mp_con,omitempty"`
	DiseaseSkill uint32 `yaml:"disease_skill,omitempty"`
	DiseaseLevel uint8  `yaml:"disease_level,omitempty"`
	AttackAfter  int    `yaml:"attack_after,omitempty"`
	PADamage     int    `yaml:"pa_damage,omitempty"`
	MADamage     int    `yaml:"ma_damage,omitempty"`
	Magic        bool   `yaml:"magic,omitempty"`
	RangeR       int    `yaml:"range_r,omitempty"`
}

type MobBanish struct {
	Message string `yaml:"message,omitempty"`
	MapID   int32  `yaml:"map_id"`
	Portal  string `yaml:"portal,omitempty"`
}

type ReactorsFile struct {
	Meta     Meta               `yaml:"_meta"`
	Reactors map[uint32]Reactor `yaml:"reactors"`
}

type Reactor struct {
	ID              uint32                 `yaml:"id"`
	Link            uint32                 `yaml:"link,omitempty"`
	ActivateByTouch int                    `yaml:"activate_by_touch,omitempty"`
	Action          string                 `yaml:"action,omitempty"`
	States          map[byte]*ReactorEvent `yaml:"states,omitempty"`
}

type ReactorEvent struct {
	Type         int   `yaml:"type"`
	NextState    byte  `yaml:"next_state"`
	TimeOut      int   `yaml:"time_out,omitempty"`
	ItemID       int   `yaml:"item_id,omitempty"`
	ItemQuantity int   `yaml:"item_quantity,omitempty"`
	LTX          int32 `yaml:"lt_x,omitempty"`
	LTY          int32 `yaml:"lt_y,omitempty"`
	RBX          int32 `yaml:"rb_x,omitempty"`
	RBY          int32 `yaml:"rb_y,omitempty"`
	TouchFlag    int   `yaml:"touch_flag,omitempty"`
	HasClickArea bool  `yaml:"has_click_area,omitempty"`
}

type ShopsFile struct {
	Meta  Meta            `yaml:"_meta"`
	Shops map[uint32]Shop `yaml:"shops"`
}

type Shop struct {
	NpcID   uint32     `yaml:"npc_id"`
	NpcName string     `yaml:"npc_name,omitempty"`
	Items   []ShopItem `yaml:"items"`
}

type ShopItem struct {
	ItemID    uint32  `yaml:"item_id"`
	Price     int     `yaml:"price"`
	Period    int     `yaml:"period,omitempty"`
	Stock     int     `yaml:"stock,omitempty"`
	UnitPrice float64 `yaml:"unit_price,omitempty"`
}

type StringsFile struct {
	Meta Meta               `yaml:"_meta"`
	Npcs map[uint32]Name    `yaml:"npcs,omitempty"`
	Mobs map[uint32]Name    `yaml:"mobs,omitempty"`
	Maps map[uint32]MapName `yaml:"maps,omitempty"`
}

type Name struct {
	Name string `yaml:"name"`
}

type MapName struct {
	StreetName string `yaml:"street_name,omitempty"`
	MapName    string `yaml:"map_name,omitempty"`
}

type DropsFile struct {
	Meta         Meta              `yaml:"_meta"`
	MobDrops     map[uint32][]Drop `yaml:"mob_drops,omitempty"`
	ReactorDrops map[uint32][]Drop `yaml:"reactor_drops,omitempty"`
}

type Drop struct {
	ItemID  uint32  `yaml:"item_id,omitempty"`
	Money   uint32  `yaml:"money,omitempty"`
	Prob    float32 `yaml:"prob"`
	Min     uint16  `yaml:"min"`
	Max     uint16  `yaml:"max"`
	QuestID uint32  `yaml:"quest_id,omitempty"`
}
