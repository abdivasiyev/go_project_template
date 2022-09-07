package models

type GetEntitiesRequest struct {
	PageRequest
	Search     string `json:"search" form:"search"`
	EntityType string `json:"entity_type" form:"entity_type" binding:"required,oneof=recruiter safety company truck trailer driver driver_type drug_test_type fuel_card_status road_test_status incab_devices external_devices need_repair_devices"`
}

type GetEntitiesResponse struct {
	Count    int                 `json:"count"`
	Entities []GetEntityResponse `json:"entities"`
}

type GetEntityResponse struct {
	ID   string `json:"id"`
	Name string `json:"name"`
}

type GetAllFormRequest struct {
	PageRequest
	Search string `json:"search"`
}

type GetAllFormResponse struct {
	Count int               `json:"count"`
	Forms []GetFormResponse `json:"forms"`
}

type CreateFormRequest struct {
	ID     string   `json:"id" swaggerignore:"true"`
	Title  string   `json:"title" binding:"required"`
	Groups []string `json:"groups" binding:"required"`
}

type GetFormResponse struct {
	ID     string                 `json:"id"`
	Title  string                 `json:"title"`
	Groups []GetFormGroupResponse `json:"groups"`
}

type CreateFormGroupRequest struct {
	ID       string                   `json:"id" swaggerignore:"true"`
	Title    string                   `json:"title" binding:"required" example:"Select driver type"`
	Sequence int                      `json:"sequence" binding:"required" example:"1"`
	Fields   []CreateFormFieldRequest `json:"fields" binding:"required"`
}

type CreateFormFieldRequest struct {
	ID            string                   `json:"id" example:"f535ef7c-2718-49ae-9fcf-65670fcad644"`
	Label         string                   `json:"label" binding:"required" example:"Select Recruiter"`
	LabelPosition string                   `json:"label_position" binding:"required,oneof=top bottom left right inside"`
	Icon          string                   `json:"icon" example:"http://example.com/pencil.svg"`
	Alias         string                   `json:"alias" example:"#select-recruiter"`
	ChildFields   []CreateFormFieldRequest `json:"child_fields,omitempty"`
	Hint          string                   `json:"hint" example:"Select recruiter from list"`
	Warning       string                   `json:"warning" example:"Give trainee handbook"`
	Sequence      int                      `json:"sequence" example:"1"`
	Type          string                   `json:"type" binding:"required,oneof=text textarea select radio checkbox file" example:"text,textarea,select,radio,checkbox,file"`
	Validation    FormValidation           `json:"validation"`
	Multiple      bool                     `json:"multiple" example:"false"`
	ListType      string                   `json:"list_url" example:"recruiter"`
	CanSearch     bool                     `json:"can_search"`
	CanAddItem    bool                     `json:"can_add_item" example:"false"`
	CanOpenPopup  bool                     `json:"can_open_popup" example:"false"`
	Grouped       bool                     `json:"grouped"`
	IsSignature   bool                     `json:"is_signature"`
	Action        string                   `json:"action" example:"recruiter/create"`
}

type GetFormGroupResponse struct {
	ID       string                 `json:"id"`
	Title    string                 `json:"title" example:"Select driver type"`
	Alias    string                 `json:"alias"`
	Sequence int                    `json:"sequence" example:"1"`
	Fields   []GetFormFieldResponse `json:"fields"`
}

type GetFormFieldResponse struct {
	ID            string                 `json:"id" example:"f535ef7c-2718-49ae-9fcf-65670fcad644"`
	Name          string                 `json:"name"`
	Label         string                 `json:"label" example:"Select Recruiter"`
	LabelPosition string                 `json:"label_position" example:"top,bottom,right,left,inside"`
	Icon          string                 `json:"icon" example:"http://example.com/pencil.svg"`
	Alias         string                 `json:"alias" example:"#select-recruiter"`
	ParentAlias   string                 `json:"parent_alias"`
	ChildFields   []GetFormFieldResponse `json:"child_fields,omitempty"`
	Hint          string                 `json:"hint" example:"Select recruiter from list"`
	Warning       string                 `json:"warning" example:"Give trainee handbook"`
	Sequence      int                    `json:"sequence" example:"1"`
	Type          string                 `json:"type" example:"text,textarea,select,radio,checkbox,file"`
	Validation    FormValidation         `json:"validation"`
	Multiple      bool                   `json:"multiple" example:"false"`
	ListType      string                 `json:"list_type" example:"recruiter"`
	CanSearch     bool                   `json:"can_search"`
	CanAddItem    bool                   `json:"can_add_item" example:"false"`
	CanOpenPopup  bool                   `json:"can_open_popup" example:"false"`
	Grouped       bool                   `json:"grouped"`
	IsSignature   bool                   `json:"is_signature"`
	Action        string                 `json:"action" example:"recruiter/create"`
	Values        []FormValue            `json:"values,omitempty"`
}

type FormValue struct {
	CreatedBy string `json:"created_by"`
	CreatedAt string `json:"created_at"`
	Value     string `json:"value"`
}

type FormValidation struct {
	Blank        bool     `json:"blank" example:"false"`
	Required     bool     `json:"required" example:"false"`
	Min          int      `json:"min" example:"1"`
	Max          int      `json:"max" example:"10"`
	Regex        string   `json:"regex" example:"[0-9]+"`
	AllowedFiles []string `json:"allowed_files" example:".jpg,.png,.pdf"`
}

type SetFormFieldValue struct {
	DriverID string               `json:"driver_id" example:"f535ef7c-2718-49ae-9fcf-65670fcad644"`
	StepID   string               `json:"step_id" example:"f535ef7c-2718-49ae-9fcf-65670fcad644"`
	Fields   []FormFieldWithValue `json:"fields"`
}

type FormFieldWithValue struct {
	FieldID string   `json:"field_id" binding:"required,uuid4" example:"f535ef7c-2718-49ae-9fcf-65670fcad644"`
	Values  []string `json:"values"`
}
