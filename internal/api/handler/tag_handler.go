package handler

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/igoventura/fintrack-core/domain"
	"github.com/igoventura/fintrack-core/internal/api/dto"
	"github.com/igoventura/fintrack-core/internal/service"
)

type TagHandler struct {
	service *service.TagService
}

func NewTagHandler(service *service.TagService) *TagHandler {
	return &TagHandler{service: service}
}

// CreateTag creates a new tag
// @Summary Create tag
// @Description Create a new tag for the authenticated user's tenant
// @Tags tags
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param request body dto.CreateTagRequest true "Create Tag Request"
// @Success 201 {object} dto.TagResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags [post]
func (h *TagHandler) CreateTag(c *gin.Context) {
	var req dto.CreateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := domain.GetUserID(c.Request.Context())
	tag := &domain.Tag{
		Name:      req.Name,
		CreatedBy: userID,
		UpdatedBy: userID,
	}

	if err := h.service.CreateTag(c.Request.Context(), tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusCreated, dto.TagResponse{
		ID:            tag.ID,
		TenantID:      tag.TenantID,
		Name:          tag.Name,
		CreatedAt:     tag.CreatedAt,
		CreatedBy:     tag.CreatedBy,
		UpdatedAt:     tag.UpdatedAt,
		UpdatedBy:     tag.UpdatedBy,
		DeactivatedAt: tag.DeactivatedAt,
	})
}

// GetTag returns a tag by ID
// @Summary Get tag
// @Description Get a tag by its ID
// @Tags tags
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Tag ID"
// @Success 200 {object} dto.TagResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/{id} [get]
func (h *TagHandler) GetTag(c *gin.Context) {
	id := c.Param("id")
	tag, err := h.service.GetTag(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TagResponse{
		ID:            tag.ID,
		TenantID:      tag.TenantID,
		Name:          tag.Name,
		CreatedAt:     tag.CreatedAt,
		CreatedBy:     tag.CreatedBy,
		UpdatedAt:     tag.UpdatedAt,
		UpdatedBy:     tag.UpdatedBy,
		DeactivatedAt: tag.DeactivatedAt,
	})
}

// ListTags returns all tags for the tenant
// @Summary List tags
// @Description Get all tags for the authenticated user's tenant
// @Tags tags
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Success 200 {object} []dto.TagResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags [get]
func (h *TagHandler) ListTags(c *gin.Context) {
	tags, err := h.service.ListTags(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	var resp []*dto.TagResponse
	for _, tag := range tags {
		resp = append(resp, &dto.TagResponse{
			ID:            tag.ID,
			TenantID:      tag.TenantID,
			Name:          tag.Name,
			CreatedAt:     tag.CreatedAt,
			CreatedBy:     tag.CreatedBy,
			UpdatedAt:     tag.UpdatedAt,
			UpdatedBy:     tag.UpdatedBy,
			DeactivatedAt: tag.DeactivatedAt,
		})
	}

	if resp == nil {
		resp = []*dto.TagResponse{}
	}
	c.JSON(http.StatusOK, resp)
}

// UpdateTag updates a tag
// @Summary Update tag
// @Description Update an existing tag
// @Tags tags
// @Accept json
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Tag ID"
// @Param request body dto.UpdateTagRequest true "Update Tag Request"
// @Success 200 {object} dto.TagResponse
// @Failure 400 {object} ErrorResponse
// @Failure 500 {object} ErrorResponse
// @Router /tags/{id} [put]
func (h *TagHandler) UpdateTag(c *gin.Context) {
	id := c.Param("id")
	var req dto.UpdateTagRequest
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, ErrorResponse{Error: err.Error()})
		return
	}

	userID := domain.GetUserID(c.Request.Context())
	tag := &domain.Tag{
		ID:        id,
		Name:      req.Name,
		UpdatedBy: userID,
	}

	if err := h.service.UpdateTag(c.Request.Context(), tag); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	updatedTag, err := h.service.GetTag(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.JSON(http.StatusOK, dto.TagResponse{
		ID:            updatedTag.ID,
		TenantID:      updatedTag.TenantID,
		Name:          updatedTag.Name,
		CreatedAt:     updatedTag.CreatedAt,
		CreatedBy:     updatedTag.CreatedBy,
		UpdatedAt:     updatedTag.UpdatedAt,
		UpdatedBy:     updatedTag.UpdatedBy,
		DeactivatedAt: updatedTag.DeactivatedAt,
	})
}

// DeleteTag deletes a tag
// @Summary Delete tag
// @Description Soft Delete a tag
// @Tags tags
// @Produce json
// @Security AuthPassword
// @Param X-Tenant-ID header string true "Tenant ID"
// @Param id path string true "Tag ID"
// @Success 204 "No Content"
// @Failure 500 {object} ErrorResponse
// @Router /tags/{id} [delete]
func (h *TagHandler) DeleteTag(c *gin.Context) {
	id := c.Param("id")
	userID := domain.GetUserID(c.Request.Context())

	if err := h.service.DeleteTag(c.Request.Context(), id, userID); err != nil {
		c.JSON(http.StatusInternalServerError, ErrorResponse{Error: err.Error()})
		return
	}

	c.Status(http.StatusNoContent)
}
