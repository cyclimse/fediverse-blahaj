package crawler

import nodeinfo "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/unversioned"

type CrawlResult struct {
	Domain   string
	Peers    []string
	NodeInfo nodeinfo.Nodeinfo
	Err      error
}

func resultFromError(domain string, err error) CrawlResult {
	return CrawlResult{
		Domain: domain,
		Err:    err,
	}
}
