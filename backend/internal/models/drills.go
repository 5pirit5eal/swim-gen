package models

type Drill struct {
	Slag             string   `json:"slug"`
	Targets          []string `json:"targets"`
	ShortDescription string   `json:"short_description"`
	ImgName          string   `json:"img_name"`
	ImgDescription   string   `json:"img_description"`
	Title            string   `json:"title"`
	Description      []string `json:"description"`
	VideoURL         []string `json:"video_url"`
	Styles           []string `json:"styles"`
	Difficulty       string   `json:"difficulty"`
	TargetGroups     []string `json:"target_groups"`
}
