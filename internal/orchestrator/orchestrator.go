package orchestrator

import (
	"context"
	"time"

	"github.com/cyclimse/fediverse-blahaj/internal/crawler"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"golang.org/x/exp/slog"
	"golang.org/x/sync/errgroup"
)

const (
	startingCrawlCapacity = 100
)

func New() *Orchestrator {
	return &Orchestrator{
		numCrawlers:   2,
		seedDomain:    "mastodon.social",
		domainCrawled: make(map[string]struct{}),
		crawlTimeout:  5 * time.Second,
	}
}

type Orchestrator struct {
	numCrawlers   int
	seedDomain    string
	domainCrawled map[string]struct{}
	crawlTimeout  time.Duration
}

// crawlerIdKey is the key for the crawler id in the context.
type crawlerIdKey struct{}

// Crawl crawls the fediverse and streams the results to the results channel.
// It returns an error if the context is exceeded.
func (o *Orchestrator) Crawl(ctx context.Context, results chan models.FediverseServer) error {
	crawlers := make([]*crawler.Crawler, o.numCrawlers)

	for i := 0; i < o.numCrawlers; i++ {
		crawlers[i] = crawler.New()
	}

	// channels for the crawl
	requested := make(chan string, startingCrawlCapacity)
	processed := make(chan crawler.CrawlResult, startingCrawlCapacity)

	// start the crawl
	requested <- o.seedDomain

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
					crawlCtx, cancel := context.WithTimeout(crawlCtx, o.crawlTimeout)
					// this could hang if the processed channel is full
					processed <- c.Crawl(crawlCtx, url)
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
					slog.ErrorCtx(ctx, "failed to crawl", "error", res.Err)
					continue
				}

				// send the peer to the results channel
				results <- models.ServerFromNodeInfo(res.Domain, res.Peers, res.NodeInfo)

				// mark the peer as crawled
				o.domainCrawled[res.Domain] = struct{}{}

				// parallelized because otherwise
				// it fails to keep up with the crawler
				// on servers with many peers
				for i := range res.Peers {
					peer := res.Peers[i]
					g.Go(func() error {
						// check if the peer was already this session
						// parallel access to the map is safe
						if _, ok := o.domainCrawled[peer]; ok {
							slog.InfoCtx(ctx, "peer was already crawled", "peer", peer)
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
