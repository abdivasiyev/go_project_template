package models

type DriverFileJobStatus int

type GetFileRequest struct {
	ID string `json:"id" uri:"id" binding:"required,uuid4"`
}

type GetFileResponse struct {
	FileID     string `json:"file_id,omitempty" example:"973cb235-bdc7-4ffc-94f8-bf4eaf23b778"`
	FileName   string `json:"file_name,omitempty" example:"Driver License.pdf"`
	FileURL    string `json:"file_url,omitempty" example:"https://cdn.example.com/files/973cb235-bdc7-4ffc-94f8-bf4eaf23b778.pdf"`
	FileStatus string `json:"file_status,omitempty" binding:"oneof=init in_process error finish"`
	FileError  string `json:"file_error,omitempty"`
}

type GenerateZipInternalRequest struct {
	EntityID        string
	EntityType      string
	RequestID       string
	GeneratedFileID string
	Status          DriverFileJobStatus
	JobError        error
}

type GenerateZipFileRequest struct {
	EntityID   string `json:"entity_id" binding:"required"`
	EntityType string `json:"entity_type" binding:"required,oneof=driver truck trailer"`
}

type GenerateZipFileResponse struct {
	RequestID string `json:"request_id"`
}

type FilesForZippingResponse struct {
	CategoryName string
	FileNames    []string
}
