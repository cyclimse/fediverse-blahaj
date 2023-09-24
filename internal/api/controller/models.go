package controller

import (
	"encoding/json"

	"github.com/cyclimse/fediverse-blahaj/internal/api/v1"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/cyclimse/fediverse-blahaj/internal/utils"
	openapi_types "github.com/deepmap/oapi-codegen/pkg/types"
)

func instanceFromModel(instance models.FediverseInstance) v1.Instance {
	return v1.Instance{
		Id:     openapi_types.UUID(instance.ID),
		Domain: instance.Domain,
		Status: v1.InstanceStatus(instance.Status),

		Description: nil,
		Software:    instance.SoftwareName,
		Version:     instance.LastCrawl.SoftwareVersion,

		NumberOfPeers: instance.LastCrawl.NumberOfPeers,

		OpenRegistrations:   instance.LastCrawl.OpenRegistrations,
		TotalUsers:          instance.LastCrawl.TotalUsers,
		ActiveUsersHalfYear: instance.LastCrawl.ActiveHalfyear,
		ActiveUsersMonth:    instance.LastCrawl.ActiveMonth,
		LocalPosts:          instance.LastCrawl.LocalPosts,
		LocalComments:       instance.LastCrawl.LocalComments,
	}
}

func crawlFromModel(crawl models.Crawl) v1.Crawl {
	rawNodeinfo := new(map[string]interface{})
	// if an error occurs, we can simply ignore it and return a nil pointer
	_ = json.Unmarshal(crawl.RawNodeinfo, rawNodeinfo)

	c := v1.Crawl{
		ActiveUsersHalfYear: crawl.ActiveHalfyear,
		ActiveUsersMonth:    crawl.ActiveMonth,
		FinishedAt:          crawl.FinishedAt,
		StartedAt:           crawl.StartedAt,
		DurationSeconds:     crawl.FinishedAt.Sub(crawl.StartedAt).Seconds(),
		Id:                  openapi_types.UUID(crawl.ID),
		InstanceId:          openapi_types.UUID(crawl.InstanceID),
		LocalComments:       crawl.LocalComments,
		LocalPosts:          crawl.LocalPosts,
		NumberOfPeers:       crawl.NumberOfPeers,
		RawNodeinfo:         rawNodeinfo,
		Status:              v1.CrawlStatus(crawl.Status),
		TotalUsers:          crawl.TotalUsers,
	}

	if crawl.Err != nil {
		c.ErrorCode = utils.ValToPtr(string(crawl.Err.Code), true)
		c.ErrorCodeDescription = utils.ValToPtr(crawl.Err.Description, crawl.Err.Description != "")
	}

	return c
}
