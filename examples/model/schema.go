// Auto generated. DO NOT EDIT IT.
// Auto generated. DO NOT EDIT IT.
// Auto generated. DO NOT EDIT IT.

package schema

import (
	"database/sql"

	"gorm.io/datatypes"

	"github.com/wonli/gormt/extype"
)

// AiSites [...]
type AiSites struct {
	ID            uint           `gorm:"autoIncrement:true;primaryKey;column:id;type:int unsigned;not null" json:"id"`
	Status        uint8          `gorm:"column:status;type:tinyint unsigned;not null;default:1;comment:0,无效 1,正常 2,推荐" json:"status"` // 0,无效 1,正常 2,推荐
	RawID         uint           `gorm:"column:raw_id;type:int unsigned;not null;default:0" json:"rawId"`
	Name          string         `gorm:"column:name;type:varchar(128);not null" json:"name"`
	Description   string         `gorm:"column:description;type:text;not null" json:"description"`
	AiDescription sql.NullString `gorm:"column:ai_description;type:text;default:null" json:"aiDescription"`
	URL           string         `gorm:"column:url;type:varchar(255);not null" json:"url"`
	Query         string         `gorm:"column:query;type:varchar(128);not null" json:"query"`
	SnapImg       string         `gorm:"column:snap_img;type:varchar(512);not null;comment:快照截图" json:"snapImg"` // 快照截图
	SnapAt        int64          `gorm:"column:snap_at;type:bigint;not null;comment:快照时间" json:"snapAt"`         // 快照时间
	MainTags      datatypes.JSON `gorm:"column:main_tags;type:json;not null" json:"mainTags"`
	AllTags       datatypes.JSON `gorm:"column:all_tags;type:json;not null" json:"allTags"`
	Remark        string         `gorm:"column:remark;type:varchar(255);not null" json:"remark"`
	CreatedAt     int64          `gorm:"column:created_at;type:bigint;not null" json:"createdAt"`
}

// Tags [...]
type Tags struct {
	TagID    uint          `gorm:"autoIncrement:true;primaryKey;column:tag_id;type:int unsigned;not null" json:"tagId"`
	TagName  string        `gorm:"column:tag_name;type:varchar(128);not null;comment:Tag名称" json:"tagName"` // Tag名称
	Location *extype.Point `gorm:"column:location;type:geometry;default:null" json:"location"`
}
