package business

//go:generate go run github.com/sqlc-dev/sqlc/cmd/sqlc generate -f ../../sqlc.yaml

import (
	"context"
	"math/rand"
	"strings"
	"time"

	"log/slog"

	"github.com/cyclimse/fediverse-blahaj/internal/db"
	"github.com/cyclimse/fediverse-blahaj/internal/models"
	"github.com/jackc/pgx/v5"
	"github.com/jackc/pgx/v5/pgxpool"
)

var (
	initialSeedDomains = []string{
		"mastodon.social",
		"mastodon.online",
		"mastodon.xyz",
		"mastodon.art",
	}
)

func New(conn *pgxpool.Pool) *Business {
	return &Business{
		conn:                  conn,
		queries:               db.New(conn),
		errorCodeDescriptions: newCachedErrorCodeDescriptions(),
	}
}

type Business struct {
	conn    *pgxpool.Pool
	queries *db.Queries

	errorCodeDescriptions cachedErrorCodeDescriptions
}

func newCachedErrorCodeDescriptions() cachedErrorCodeDescriptions {
	expireAfter := time.Hour
	return cachedErrorCodeDescriptions{
		Descriptions: make(map[models.CrawlErrCode]string),
		LastUpdate:   time.Now().Add(-expireAfter),
		// in practice, the descriptions should not change
		// however, we still want to be able to change them
		// without having to restart the server
		ExpireAfter: expireAfter,
	}
}

type cachedErrorCodeDescriptions struct {
	Descriptions map[models.CrawlErrCode]string
	LastUpdate   time.Time
	ExpireAfter  time.Duration
}

// Run runs the business logic.
// Handles the results from the crawler, but is not aware of the crawler logic.
func (b *Business) Run(ctx context.Context, crawls chan models.Crawl) error {
	for {
		select {
		case <-ctx.Done():
			return ctx.Err()
		case crawl, ok := <-crawls:
			if !ok {
				return nil
			}
			slog.InfoContext(ctx, "received crawl", "crawl", crawl.Domain, "status", crawl.Status)
			// check if the instance is already in the db
			var instance db.Instance
			instance, err := b.queries.GetInstanceByDomain(ctx, crawl.Domain)
			if err != nil {
				if err != pgx.ErrNoRows {
					return err
				}
				// instance is not in the db, add it
				instance, err = b.AddInstance(ctx, crawl.Domain, crawl.SoftwareName)
				if err != nil {
					slog.ErrorContext(ctx, "failed to add instance", "error", err)
					return err
				}
			}
			// add the crawl to the instance
			err = b.AddCrawlToFediverseInstance(ctx, crawl, instance)
			if err != nil {
				slog.ErrorContext(ctx, "failed to add crawl to instance", "error", err)
				return err
			}
		}
	}
}

func (b *Business) GetCrawlerSeedDomains(ctx context.Context, count int) ([]string, error) {
	domains := make([]string, 0, count)

	offset := 0

	for len(domains) < count {
		newDomains, err := b.queries.GetCrawlerSeedDomains(ctx, db.GetCrawlerSeedDomainsParams{
			Offset: int32(offset),
			Limit:  int32(count - len(domains)),
		})
		if err != nil {
			return nil, err
		}

		if len(newDomains) == 0 {
			break
		}

		for _, domain := range newDomains {
			if !isBlocked(domain) {
				domains = append(domains, domain)
			}
		}
	}

	if len(domains) == 0 {
		domains = append(domains, initialSeedDomains...)
	}

	// For now until we have a better way to handle this
	// We shuffle to make it more random

	// (no need to seed in go1.20+)
	rand.Shuffle(len(domains), func(i, j int) {
		domains[i], domains[j] = domains[j], domains[i]
	})

	return domains, nil
}

// isBlocked returns true if the domain is blocked.
func isBlocked(domain string) bool {
	// we need to block domains and subdomains
	// e.g. blocking ngrok.io should also block a.ngrok.io
	for _, blockedDomain := range BlockedDomains {
		if domain == blockedDomain {
			return true
		}

		if strings.HasSuffix(domain, blockedDomain) {
			return true
		}
	}

	return false
}
