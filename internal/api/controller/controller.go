package controller

import (
	"errors"
	"net/http"

	"github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/business"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
	"golang.org/x/exp/slog"
)

func NewAPIController(business *business.Business) v1.ServerInterface {
	return &APIController{
		Business: business,
	}
}

type APIController struct {
	Business *business.Business
}

// GetServerById implements v1.ServerInterface
func (c *APIController) GetServerByID(ctx echo.Context, id uuid.UUID) error {
	s, err := c.Business.GetServerByID(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, business.ErrServerNotFound) {
			e := v1.Error{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}
			return ctx.JSON(http.StatusNotFound, e)
		}
	}
	return ctx.JSON(http.StatusOK, s)
}

// ListServers implements v1.ServerInterface
func (c *APIController) ListServers(ctx echo.Context, params v1.ListServersParams) error {
	page, pageSize := validatePage(params.Page), validatePageSize(params.PerPage)

	servers, total, err := c.Business.ListServers(ctx.Request().Context(), page, pageSize)
	if err != nil {
		slog.ErrorCtx(ctx.Request().Context(), "failed to list servers", "error", err)
		return err
	}

	var resp = v1.ListServers200JSONResponse{
		Results: make([]v1.Server, len(servers)),
		Page:    page,
		PerPage: pageSize,
		Total:   total,
	}

	for i, s := range servers {
		resp.Results[i] = v1.Server{
			Id:     openapi_types.UUID(s.ID),
			Domain: s.Domain,

			Description: nil,
			Software:    nil,

			NumberOfPeers: s.NumberOfPeers,

			OpenRegistrations:   s.OpenRegistrations,
			TotalUsers:          s.TotalUsers,
			ActiveUsersHalfYear: s.ActiveHalfyear,
			ActiveUsersMonth:    s.ActiveMonth,
			LocalPosts:          s.LocalPosts,
			LocalComments:       s.LocalComments,
		}
	}

	return ctx.JSON(http.StatusOK, resp)
}

// ListCrawlsForServer implements v1.ServerInterface
func (*APIController) ListCrawlsForServer(ctx echo.Context, id uuid.UUID, params v1.ListCrawlsForServerParams) error {
	panic("unimplemented")
}
