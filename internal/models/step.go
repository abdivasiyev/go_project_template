package models

type UpdateStepStatusRequest struct {
	DriverID string `json:"driver_id" binding:"required,uuid4"`
	StepID   string `json:"step_id" binding:"required,uuid4"`
	StatusID string `json:"status_id" binding:"required,uuid4"`
}

type GetAllStepsRequest struct {
	DepartmentID string `json:"department_id" uri:"department_id" binding:"required,uuid4" swaggerignore:"true"`
	DriverID     string `json:"driver_id" uri:"driver_id" binding:"required,uuid4" swaggerignore:"true"`
}

type GetAllStepsGroupByStatusResponse struct {
	Count int                            `json:"count"`
	Steps []GetStepGroupByStatusResponse `json:"steps"`
}

type GetStepGroupByStatusResponse struct {
	Status GetStatusResponse `json:"status"`
	Steps  []GetStepResponse `json:"steps"`
}

type GetAllStepsResponse struct {
	Count int               `json:"count"`
	Steps []GetStepResponse `json:"steps"`
}

type StepFields struct {
	ID     string                  `json:"id"`
	StepID string                  `json:"step_id"`
	Alias  string                  `json:"alias"`
	Label  string                  `json:"label"`
	Values []GetFieldValueResponse `json:"values"`
}

type GetFieldValueResponse struct {
	ID    string `json:"id"`
	Value string `json:"value"`
}

type GetStepResponse struct {
	ID       string       `json:"id"`
	Alias    string       `json:"alias"`
	Name     string       `json:"name"`
	Sequence int          `json:"sequence"`
	StatusID string       `json:"-"`
	Fields   []StepFields `json:"fields"`
}
