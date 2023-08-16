package crawler

import (
	"errors"
	"fmt"
)

type CrawlStatus string

const (
	CrawlStatusUnknown              CrawlStatus = "unknown"
	CrawlStatusCompleted            CrawlStatus = "completed"
	CrawlStatusFailed               CrawlStatus = "failed"
	CrawlStatusBlocked              CrawlStatus = "blocked"
	CrawlStatusTimeout              CrawlStatus = "timeout"
	CrawlStatusInternalError        CrawlStatus = "internal_error"
)

type CrawlError interface {
	error
	Status() CrawlStatus
}

// Only the status is relevant outside of the crawler package
var (
	// errCrawlTimeout is returned when the crawl times out
	errCrawlTimeout = newCrawlError(CrawlStatusTimeout, errors.New("crawl timeout"))
	// errDomainLookupFailed is returned when the domain lookup fails
	errDomainLookupFailed = newCrawlError(CrawlStatusFailed, errors.New("domain lookup failed"))
	// errNetworkError is returned when the network request fails
	errNetworkError = newCrawlError(CrawlStatusFailed, errors.New("network error"))
	// errRobotsTxtDisallowsCrawling is returned when the robots.txt disallows crawling
	errRobotsTxtDisallowsCrawling = newCrawlError(CrawlStatusBlocked, errors.New("robots.txt disallows crawling"))
	// errUnsupportedSoftware is returned when the software is not supported
	errUnsupportedSoftware = newCrawlError(CrawlStatusFailed, errors.New("unsupported software"))
	// errNodeInfoSyntax is returned when the node info syntax is invalid
	errNodeInfoSyntax = newCrawlError(CrawlStatusInternalError, errors.New("node info syntax"))
)

func newCrawlError(status CrawlStatus, err error) *crawlError {
	return &crawlError{
		status: status,
		err:    err,
	}
}

func newCrawlInternalError(err error) *crawlError {
	return &crawlError{
		status: CrawlStatusInternalError,
		err:    err,
	}
}

type crawlError struct {
	status CrawlStatus
	err    error
}

func (e *crawlError) Error() string {
	return fmt.Sprintf("%s: %s", e.status, e.err)
}

func (e *crawlError) Unwrap() error {
	return e.err
}

func (e *crawlError) Status() CrawlStatus {
	return e.status
}

func (e *crawlError) Wrap(err error) *crawlError {
	return &crawlError{
		status: e.status,
		err:    errors.Join(e.err, err),
	}
}
