package handlers

import (
	"errors"
	"github.com/celestiaorg/validator-da-tracker/pkg/repository"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
	"net/http"
)

type ValidatorHandler struct {
	Repo repository.ValidatorRepository
}

// NewValidatorHandler creates a new instance of ValidatorHandler
func NewValidatorHandler(repo repository.ValidatorRepository) *ValidatorHandler {
	return &ValidatorHandler{Repo: repo}
}

// GetAllValidators returns a list of all validators of the validator table
func (h *ValidatorHandler) GetAllValidators() gin.HandlerFunc {
	return func(c *gin.Context) {
		validators, err := h.Repo.GetAllValidators(c.Request.Context())
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		c.JSON(http.StatusOK, validators)
	}
}

// GetValidatorByID returns a validator by id of the validator table
func (h *ValidatorHandler) GetValidatorByID() gin.HandlerFunc {
	return func(c *gin.Context) {
		id := c.Param("id")

		validator, err := h.Repo.GetValidatorByID(c.Request.Context(), id)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Validator not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, validator)
	}
}

// GetValidatorByEmail returns a validator by email of the validator table
func (h *ValidatorHandler) GetValidatorByEmail() gin.HandlerFunc {
	return func(c *gin.Context) {
		email := c.Query("email")

		validator, err := h.Repo.GetValidatorByEmail(c.Request.Context(), email)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Validator not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, validator)
	}
}

// GetValidatorByName returns a validator by name of the validator table
func (h *ValidatorHandler) GetValidatorByName() gin.HandlerFunc {
	return func(c *gin.Context) {
		name := c.Query("name")

		validator, err := h.Repo.GetValidatorByName(c.Request.Context(), name)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Validator not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, validator)
	}
}

// GetValidatorsByPeerIDs returns a list of validators by peerIDs of the validator table
func (h *ValidatorHandler) GetValidatorsByPeerIDs() gin.HandlerFunc {
	return func(c *gin.Context) {
		peerIDs := c.QueryArray("peerid")

		validators, err := h.Repo.GetValidatorsByPeerIDs(c.Request.Context(), peerIDs)
		if err != nil {
			if errors.Is(err, gorm.ErrRecordNotFound) {
				c.JSON(http.StatusNotFound, gin.H{"error": "Validator not found"})
				return
			}
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, validators)
	}
}
