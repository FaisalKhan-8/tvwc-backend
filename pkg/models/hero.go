package models

type HeroSection struct {
	ID            string `json:"id" bson:"_id,omitempty"`
	HeadingText   string `json:"heading_text" bson:"heading_text"`
	SubHeadingText string `json:"sub_heading_text" bson:"sub_heading_text"`
	ToolTipName   string `json:"tool_tip_name" bson:"tool_tip_name"`
	Image         string `json:"image" bson:"image"`
	Designation   string `json:"designation" bson:"designation"`
}