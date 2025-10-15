package handler

import (
	"errors"
	"net/http"
	"reflect"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/ilhamgepe/gepay/internal/domain"
	"github.com/ilhamgepe/gepay/internal/server/middleware"
	"github.com/ilhamgepe/gepay/pkg/logger"
	"github.com/ilhamgepe/gepay/pkg/utils"
)

type merchantHandler struct {
	merchantService domain.MerchantService
	mw              *middleware.Middleware
	log             logger.Logger
}

func NewMerchantHandler(srv *gin.RouterGroup, merchantService domain.MerchantService, mw *middleware.Middleware, log logger.Logger) *merchantHandler {
	handlers := &merchantHandler{
		merchantService: merchantService,
		mw:              mw,
		log:             log,
	}
	srv.POST("/kyc", mw.WithAuth, handlers.kyc)

	return handlers
}

func (h *merchantHandler) kyc(c *gin.Context) {
	userIDCtx, ok := c.Get(domain.CtxUserIDKey)
	if !ok {
		c.JSON(400, domain.ErrorResponse{
			Errors: map[string]any{
				"error": domain.ErrInternalServerError.Error(),
			},
		})
		return
	}
	userID, err := uuid.Parse(userIDCtx.(string))
	if err != nil {
		c.JSON(400, domain.ErrorResponse{Errors: map[string]any{"error": err.Error()}})
		return
	}

	var req domain.MerchantKYCReq
	if err := c.ShouldBind(&req); err != nil {
		h.log.Error(err, "TEST")
		result := utils.GenerateMessage(err, reflect.TypeOf(req))
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Message: "invalid request",
			Errors:  result,
		})
		return
	}

	if req.KYCType == "business" {
		messages, error := businessKYCValidation(req)
		if error {
			c.JSON(http.StatusBadRequest, domain.ErrorResponse{
				Message: "invalid request",
				Errors:  messages,
			})
			return
		}
	}

	req.UserID = userID
	err = h.merchantService.CreateMerchantKYC(c.Request.Context(), req)
	if err != nil {
		if errors.Is(err, domain.ErrInternalServerError) {
			c.JSON(http.StatusInternalServerError, domain.ErrorResponse{
				Errors: map[string]any{
					"error": err.Error(),
				},
			})
			return
		}
		c.JSON(http.StatusBadRequest, domain.ErrorResponse{
			Errors: map[string]any{
				"error": err.Error(),
			},
		})
		return
	}

	c.JSON(http.StatusOK, domain.SuccessResponse{
		Message: "success",
	})
}

func businessKYCValidation(req domain.MerchantKYCReq) (map[string]any, bool) {
	messages := make(map[string]any)
	error := false
	//  business license number / NIB
	if req.BusinessLicenseNumber == "" || len(req.BusinessLicenseNumber) > 20 || len(req.BusinessLicenseNumber) < 13 {
		messages["business_license_number"] = "business_license_number is required and must be 13-20 characters long"
		error = true
	}

	// NPWP / Tax ID
	if req.TaxIDNumber == "" || len(req.TaxIDNumber) > 16 || len(req.TaxIDNumber) < 15 {
		messages["tax_id_number"] = "tax_id_number is required and must be 15-16 characters long"
		error = true
	}

	// deed number / akta pendirian
	if req.DeedNumber == "" || len(req.DeedNumber) > 20 || len(req.DeedNumber) < 13 {
		messages["deed_number"] = "deed_number is required and must be 13-20 characters long"
		error = true
	}

	// website url
	if req.WebsiteURL == "" {
		messages["website_url"] = "website_url is required"
		error = true
	}

	return messages, error
}
