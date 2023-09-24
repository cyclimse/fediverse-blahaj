package orchestrator

import (
	"context"
	"strings"
	"time"

	"log/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/crawler"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"golang.org/x/sync/errgroup"
)

const (
	startingCrawlCapacity = 100
)

func New(config OrchestratorConfig) *Orchestrator {
	return &Orchestrator{
		crawledDomains: make(map[string]struct{}),
		config:         config,
	}
}

type OrchestratorConfig struct {
	NumCrawlers      int
	BlockedDomains   []string
	SeedDomains      []string
	CrawlTimeout     time.Duration
	CrawlerUserAgent string
}

type Orchestrator struct {
	crawledDomains map[string]struct{}
	config         OrchestratorConfig
}

// crawlerIdKey is the key for the crawler id in the context.
type crawlerIdKey struct{}

// Crawl crawls the fediverse and streams the results to the results channel.
// It returns an error if the context is exceeded.
func (o *Orchestrator) Crawl(ctx context.Context, results chan models.Crawl) error {
	crawlers := make([]*crawler.Crawler, o.config.NumCrawlers)

	for i := 0; i < o.config.NumCrawlers; i++ {
		crawlers[i] = crawler.New(o.config.CrawlerUserAgent)
	}

	// channels for the crawl
	requested := make(chan string, startingCrawlCapacity)
	processed := make(chan crawler.CrawlResult, startingCrawlCapacity)

	// start the crawl
	for _, domain := range o.config.SeedDomains {
		requested <- domain
	}

	for i := range crawlers {
		c := crawlers[i]
		// capture as argument to avoid loopclosure issues
		go func(i int) {
			for {
				select {
				case <-ctx.Done():
					return
				case url := <-requested:
					crawlCtx := context.WithValue(ctx, crawlerIdKey{}, i)
					crawlCtx, cancel := context.WithTimeout(crawlCtx, o.config.CrawlTimeout)
					// this could hang if the processed channel is full
					processed <- *c.Crawl(crawlCtx, url)
					cancel()
				}
			}
		}(i)
	}

	g, ctx := errgroup.WithContext(ctx)
	g.SetLimit(-1)

	g.Go(func() error {
		for {
			select {
			case <-ctx.Done():
				return ctx.Err()
			case res := <-processed:
				if res.Err != nil {
					slog.ErrorContext(ctx, "failed to crawl", "domain", res.Domain, "error", res.Err)
				}
				// send the peer to the results channel
				results <- crawler.CrawlFromResult(res)

				// mark the peer as crawled
				o.crawledDomains[res.Domain] = struct{}{}

				// parallelized because otherwise
				// it fails to keep up with the crawler
				// on servers with many peers
				for i := range res.Peers {
					peer := res.Peers[i]
					g.Go(func() error {
						// check if the peer was already this session
						if o.wasCrawled(peer) {
							slog.InfoContext(ctx, "peer was already crawled", "peer", peer)
							return nil
						}

						if o.isBlocked(peer) {
							slog.InfoContext(ctx, "peer is blocked", "peer", peer)
							return nil
						}

						// if the context is exceeded, the workers have exited and we can stop
						// to avoid a deadlock
						select {
						case <-ctx.Done():
							return ctx.Err()
						case requested <- peer:
						}

						return nil
					})
				}

			}
		}
	})

	// this will hang until the context is exceeded
	// this is expected
	return g.Wait()
}

func (o *Orchestrator) wasCrawled(domain string) bool {
	// parallel access to the map is safe
	_, ok := o.crawledDomains[domain]
	return ok
}

// isBlocked returns true if the domain is blocked.
func (o *Orchestrator) isBlocked(domain string) bool {
	// we need to block domains and subdomains
	// e.g. blocking ngrok.io should also block a.ngrok.io
	for _, blockedDomain := range o.config.BlockedDomains {
		if domain == blockedDomain {
			return true
		}

		if strings.HasSuffix(domain, blockedDomain) {
			return true
		}
	}

	return false
}
