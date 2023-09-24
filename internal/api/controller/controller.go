package controller

//go:generate go run github.com/deepmap/oapi-codegen/cmd/oapi-codegen -package v1 -generate "types,server,strict-server,spec" -o ../v1/openapi.gen.go ../../../api/v1/openapi.yaml

import (
	"errors"
	"net/http"

	"log/slog"

	v1 "github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/business"
	"github.com/google/uuid"
	"github.com/labstack/echo/v4"
)

func NewAPIController(business *business.Business) v1.ServerInterface {
	return &APIController{
		Business: business,
	}
}

type APIController struct {
	Business *business.Business
}

// GetInstanceById implements v1.InstanceInterface
func (c *APIController) GetInstanceByID(ctx echo.Context, id uuid.UUID) error {
	instance, err := c.Business.GetInstanceByID(ctx.Request().Context(), id)
	if err != nil {
		if errors.Is(err, business.ErrInstanceNotFound) {
			e := v1.Error{
				Code:    http.StatusNotFound,
				Message: err.Error(),
			}
			return ctx.JSON(http.StatusNotFound, e)
		}
	}
	return ctx.JSON(http.StatusOK, instanceFromModel(instance))
}

// ListInstances implements v1.InstanceInterface
func (c *APIController) ListInstances(ctx echo.Context, params v1.ListInstancesParams) error {
	page, pageSize := validatePage(params.Page), validatePageSize(params.PerPage)

	instances, total, err := c.Business.ListInstances(ctx.Request().Context(), page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx.Request().Context(), "failed to list instances", "error", err)
		return err
	}

	var resp = v1.ListInstances200JSONResponse{
		Results: make([]v1.Instance, len(instances)),
		Page:    page,
		PerPage: pageSize,
		Total:   total,
	}

	for i, s := range instances {
		resp.Results[i] = instanceFromModel(s)
	}

	return ctx.JSON(http.StatusOK, resp)
}

// ListCrawlsForInstance implements v1.InstanceInterface
func (c *APIController) ListCrawlsForInstance(ctx echo.Context, id uuid.UUID, params v1.ListCrawlsForInstanceParams) error {
	page, pageSize := validatePage(params.Page), validatePageSize(params.PerPage)

	crawls, total, err := c.Business.ListCrawlsForInstance(ctx.Request().Context(), id, page, pageSize)
	if err != nil {
		slog.ErrorContext(ctx.Request().Context(), "failed to list crawls", "error", err, "instance_id", id)
		return err
	}

	var resp = v1.ListCrawlsForInstance200JSONResponse{
		Results: make([]v1.Crawl, len(crawls)),
		Page:    page,
		PerPage: pageSize,
		Total:   total,
	}

	for i, c := range crawls {
		resp.Results[i] = crawlFromModel(c)
	}

	return ctx.JSON(http.StatusOK, resp)
}
