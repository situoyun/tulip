package app

import "time"

// Model ...
type Model struct {
	CreatedAt time.Time  `json:"created_at,omitempty"`
	UpdatedAt time.Time  `json:"updated_at,omitempty"`
	DeletedAt *time.Time `json:"deleted_at,omitempty"`
}

type User struct {
	ID           uint       `gorm:"primary_key" json:"id,omitempty"`
	UUID         string     `gorm:"column:uid;unique_index" sql:"default:''" json:"uuid,omitempty"`
	Name         string     `gorm:"column:name;index" sql:"default:''" json:"name,omitempty"`
	Alias        string     `gorm:"column:alias;index" sql:"default:''" json:"alias,omitempty"`
	Tel          string     `json:"tel,omitempty" gorm:"index" sql:"default:''"`
	Mail         string     `json:"mail,omitempty" gorm:"index" sql:"default:''"`
	OpenID       string     `sql:"default:''" json:"open_id,omitempty"`
	PassWord     string     `gorm:"column:password" sql:"not null" json:"password,omitempty"`
	PassWordType int        `gorm:"column:password_type" sql:"default:0" json:"pass_word_type,omitempty"` //0表示默认 1表示只用md5
	Salt         string     `gorm:"column:salt" sql:"default:''" json:"salt,omitempty"`
	Area         string     `gorm:"column:area" sql:"default:''" json:"area,omitempty"`
	LoginType    uint       `gorm:"column:login_type" json:"-,omitempty"` //1账号 2微信
	InvitedCode  string     `gorm:"column:invited_code;index" sql:"default:''" json:"invited_code,omitempty"`
	ModifyPWD    bool       `gorm:"column:modify_pwd" json:"modify_pwd,omitempty"`
	Avatar       string     `sql:"default:''" json:"avatar,omitempty"`
	Sex          string     `sql:"default:''" json:"sex,omitempty"`
	IdNumber     string     `gorm:"type:varchar(100);default:''" json:"id_number,omitempty"`
	Hobby        string     `json:"hobby,omitempty"`
	Birth        string     `json:"birth,omitempty"`
	BornDate     *time.Time `json:"born_date,omitempty"`

	Model
}
