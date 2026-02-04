package dto

import "time"

type CreateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type UpdateTagRequest struct {
	Name string `json:"name" binding:"required"`
}

type TagResponse struct {
	ID            string     `json:"id"`
	TenantID      string     `json:"tenant_id"`
	Name          string     `json:"name"`
	CreatedAt     time.Time  `json:"created_at"`
	CreatedBy     string     `json:"created_by"`
	UpdatedAt     time.Time  `json:"updated_at"`
	UpdatedBy     string     `json:"updated_by"`
	DeactivatedAt *time.Time `json:"deactivated_at,omitempty"`
}
