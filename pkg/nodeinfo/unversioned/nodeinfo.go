package unversioned

import (
	"fmt"

	v10 "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v10"
	v11 "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v11"
	v20 "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v20"
	v21 "github.com/cyclimse/fediverse-blahaj/pkg/nodeinfo/v21"
)

const schemaPrefix = "http://nodeinfo.diaspora.software/ns/schema/"

var knownVersions = []string{
	"2.1",
	"2.0",
	"1.1",
	"1.0",
}

type Nodeinfo interface {
	SchemaVersion() string
	SoftwareName() string
	SoftwareVersion() string
	IsRegistrationOpen() bool
	TotalUsers() *int
	ActiveUsersHalfyear() *int
	ActiveUsersMonth() *int
	LocalPosts() *int
	LocalComments() *int
}

type WellKnownLink struct {
	Rel  string `json:"rel"`
	Href string `json:"href"`
}

type WellKnown struct {
	Links []WellKnownLink `json:"links"`
}

// HighestSupported returns the highest supported nodeinfo version
// and the corresponding nodeinfo struct to decode into.
// If no supported version is found, an error is returned.
func HighestSupported(w WellKnown) (url string, n Nodeinfo, err error) {
	hist := make(map[string]string, len(w.Links))
	for _, link := range w.Links {
		// verify that the link is a nodeinfo link
		if len(link.Rel) < len(schemaPrefix) || link.Rel[:len(schemaPrefix)] != schemaPrefix {
			continue
		}
		hist[link.Rel] = link.Href
	}

	// the order of knownVersions is important,
	// because we want to return the highest supported version

	for _, v := range knownVersions {
		if href, ok := hist[schemaPrefix+v]; ok {
			return href, nodeInfoForVersion(v), nil
		}
	}

	return "", nil, fmt.Errorf("no supported nodeinfo version found")
}

func nodeInfoForVersion(version string) Nodeinfo {
	switch version {
	case "2.1":
		return &v21.Nodeinfo{}
	case "2.0":
		return &v20.Nodeinfo{}
	case "1.1":
		return &v11.Nodeinfo{}
	case "1.0":
		return &v10.Nodeinfo{}
	}
	return nil
}
