// Code generated by github.com/atombender/go-jsonschema, DO NOT EDIT.

package v20

import "encoding/json"
import "fmt"
import "reflect"

// Free form key value pairs for software specific values. Clients should not rely
// on any specific key present.
type NodeinfoMetadata map[string]interface{}

type NodeinfoProtocolsElem string

var enumValues_NodeinfoProtocolsElem = []interface{}{
	"activitypub",
	"buddycloud",
	"dfrn",
	"diaspora",
	"libertree",
	"ostatus",
	"pumpio",
	"tent",
	"xmpp",
	"zot",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoProtocolsElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_NodeinfoProtocolsElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_NodeinfoProtocolsElem, v)
	}
	*j = NodeinfoProtocolsElem(v)
	return nil
}

const NodeinfoProtocolsElemActivitypub NodeinfoProtocolsElem = "activitypub"
const NodeinfoProtocolsElemBuddycloud NodeinfoProtocolsElem = "buddycloud"
const NodeinfoProtocolsElemDfrn NodeinfoProtocolsElem = "dfrn"
const NodeinfoProtocolsElemDiaspora NodeinfoProtocolsElem = "diaspora"
const NodeinfoProtocolsElemLibertree NodeinfoProtocolsElem = "libertree"
const NodeinfoProtocolsElemOstatus NodeinfoProtocolsElem = "ostatus"
const NodeinfoProtocolsElemPumpio NodeinfoProtocolsElem = "pumpio"
const NodeinfoProtocolsElemTent NodeinfoProtocolsElem = "tent"
const NodeinfoProtocolsElemXmpp NodeinfoProtocolsElem = "xmpp"
const NodeinfoProtocolsElemZot NodeinfoProtocolsElem = "zot"

type NodeinfoServicesInboundElem string

var enumValues_NodeinfoServicesInboundElem = []interface{}{
	"atom1.0",
	"gnusocial",
	"imap",
	"pnut",
	"pop3",
	"pumpio",
	"rss2.0",
	"twitter",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoServicesInboundElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_NodeinfoServicesInboundElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_NodeinfoServicesInboundElem, v)
	}
	*j = NodeinfoServicesInboundElem(v)
	return nil
}

const NodeinfoServicesInboundElemAtom10 NodeinfoServicesInboundElem = "atom1.0"
const NodeinfoServicesInboundElemGnusocial NodeinfoServicesInboundElem = "gnusocial"
const NodeinfoServicesInboundElemImap NodeinfoServicesInboundElem = "imap"
const NodeinfoServicesInboundElemPnut NodeinfoServicesInboundElem = "pnut"
const NodeinfoServicesInboundElemPop3 NodeinfoServicesInboundElem = "pop3"
const NodeinfoServicesInboundElemPumpio NodeinfoServicesInboundElem = "pumpio"
const NodeinfoServicesInboundElemRss20 NodeinfoServicesInboundElem = "rss2.0"
const NodeinfoServicesInboundElemTwitter NodeinfoServicesInboundElem = "twitter"

type NodeinfoServicesOutboundElem string

var enumValues_NodeinfoServicesOutboundElem = []interface{}{
	"atom1.0",
	"blogger",
	"buddycloud",
	"diaspora",
	"dreamwidth",
	"drupal",
	"facebook",
	"friendica",
	"gnusocial",
	"google",
	"insanejournal",
	"libertree",
	"linkedin",
	"livejournal",
	"mediagoblin",
	"myspace",
	"pinterest",
	"pnut",
	"posterous",
	"pumpio",
	"redmatrix",
	"rss2.0",
	"smtp",
	"tent",
	"tumblr",
	"twitter",
	"wordpress",
	"xmpp",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoServicesOutboundElem) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_NodeinfoServicesOutboundElem {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_NodeinfoServicesOutboundElem, v)
	}
	*j = NodeinfoServicesOutboundElem(v)
	return nil
}

// The third party sites this server can connect to via their application API.
type NodeinfoServices struct {
	// The third party sites this server can retrieve messages from for combined
	// display with regular traffic.
	Inbound []NodeinfoServicesInboundElem `json:"inbound" yaml:"inbound" mapstructure:"inbound"`

	// The third party sites this server can publish messages to on the behalf of a
	// user.
	Outbound []NodeinfoServicesOutboundElem `json:"outbound" yaml:"outbound" mapstructure:"outbound"`
}

const NodeinfoServicesOutboundElemAtom10 NodeinfoServicesOutboundElem = "atom1.0"
const NodeinfoServicesOutboundElemBlogger NodeinfoServicesOutboundElem = "blogger"
const NodeinfoServicesOutboundElemBuddycloud NodeinfoServicesOutboundElem = "buddycloud"
const NodeinfoServicesOutboundElemDiaspora NodeinfoServicesOutboundElem = "diaspora"
const NodeinfoServicesOutboundElemDreamwidth NodeinfoServicesOutboundElem = "dreamwidth"
const NodeinfoServicesOutboundElemDrupal NodeinfoServicesOutboundElem = "drupal"
const NodeinfoServicesOutboundElemFacebook NodeinfoServicesOutboundElem = "facebook"
const NodeinfoServicesOutboundElemFriendica NodeinfoServicesOutboundElem = "friendica"
const NodeinfoServicesOutboundElemGnusocial NodeinfoServicesOutboundElem = "gnusocial"
const NodeinfoServicesOutboundElemGoogle NodeinfoServicesOutboundElem = "google"
const NodeinfoServicesOutboundElemInsanejournal NodeinfoServicesOutboundElem = "insanejournal"
const NodeinfoServicesOutboundElemLibertree NodeinfoServicesOutboundElem = "libertree"
const NodeinfoServicesOutboundElemLinkedin NodeinfoServicesOutboundElem = "linkedin"
const NodeinfoServicesOutboundElemLivejournal NodeinfoServicesOutboundElem = "livejournal"
const NodeinfoServicesOutboundElemMediagoblin NodeinfoServicesOutboundElem = "mediagoblin"
const NodeinfoServicesOutboundElemMyspace NodeinfoServicesOutboundElem = "myspace"
const NodeinfoServicesOutboundElemPinterest NodeinfoServicesOutboundElem = "pinterest"
const NodeinfoServicesOutboundElemPnut NodeinfoServicesOutboundElem = "pnut"
const NodeinfoServicesOutboundElemPosterous NodeinfoServicesOutboundElem = "posterous"
const NodeinfoServicesOutboundElemPumpio NodeinfoServicesOutboundElem = "pumpio"
const NodeinfoServicesOutboundElemRedmatrix NodeinfoServicesOutboundElem = "redmatrix"
const NodeinfoServicesOutboundElemRss20 NodeinfoServicesOutboundElem = "rss2.0"
const NodeinfoServicesOutboundElemSmtp NodeinfoServicesOutboundElem = "smtp"
const NodeinfoServicesOutboundElemTent NodeinfoServicesOutboundElem = "tent"
const NodeinfoServicesOutboundElemTumblr NodeinfoServicesOutboundElem = "tumblr"
const NodeinfoServicesOutboundElemTwitter NodeinfoServicesOutboundElem = "twitter"
const NodeinfoServicesOutboundElemWordpress NodeinfoServicesOutboundElem = "wordpress"
const NodeinfoServicesOutboundElemXmpp NodeinfoServicesOutboundElem = "xmpp"

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoServices) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["inbound"]; !ok || v == nil {
		return fmt.Errorf("field inbound in NodeinfoServices: required")
	}
	if v, ok := raw["outbound"]; !ok || v == nil {
		return fmt.Errorf("field outbound in NodeinfoServices: required")
	}
	type Plain NodeinfoServices
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = NodeinfoServices(plain)
	return nil
}

// Metadata about server software in use.
type NodeinfoSoftware struct {
	// The canonical name of this server software.
	Name string `json:"name" yaml:"name" mapstructure:"name"`

	// The version of this server software.
	Version string `json:"version" yaml:"version" mapstructure:"version"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoSoftware) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["name"]; !ok || v == nil {
		return fmt.Errorf("field name in NodeinfoSoftware: required")
	}
	if v, ok := raw["version"]; !ok || v == nil {
		return fmt.Errorf("field version in NodeinfoSoftware: required")
	}
	type Plain NodeinfoSoftware
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = NodeinfoSoftware(plain)
	return nil
}

// Usage statistics for this server.
type NodeinfoUsage struct {
	// The amount of comments that were made by users that are registered on this
	// server.
	LocalComments *int `json:"localComments,omitempty" yaml:"localComments,omitempty" mapstructure:"localComments,omitempty"`

	// The amount of posts that were made by users that are registered on this server.
	LocalPosts *int `json:"localPosts,omitempty" yaml:"localPosts,omitempty" mapstructure:"localPosts,omitempty"`

	// statistics about the users of this server.
	Users NodeinfoUsageUsers `json:"users" yaml:"users" mapstructure:"users"`
}

// statistics about the users of this server.
type NodeinfoUsageUsers struct {
	// The amount of users that signed in at least once in the last 180 days.
	ActiveHalfyear *int `json:"activeHalfyear,omitempty" yaml:"activeHalfyear,omitempty" mapstructure:"activeHalfyear,omitempty"`

	// The amount of users that signed in at least once in the last 30 days.
	ActiveMonth *int `json:"activeMonth,omitempty" yaml:"activeMonth,omitempty" mapstructure:"activeMonth,omitempty"`

	// The total amount of on this server registered users.
	Total *int `json:"total,omitempty" yaml:"total,omitempty" mapstructure:"total,omitempty"`
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoUsage) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["users"]; !ok || v == nil {
		return fmt.Errorf("field users in NodeinfoUsage: required")
	}
	type Plain NodeinfoUsage
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	*j = NodeinfoUsage(plain)
	return nil
}

type NodeinfoVersion string

var enumValues_NodeinfoVersion = []interface{}{
	"2.0",
}

// UnmarshalJSON implements json.Unmarshaler.
func (j *NodeinfoVersion) UnmarshalJSON(b []byte) error {
	var v string
	if err := json.Unmarshal(b, &v); err != nil {
		return err
	}
	var ok bool
	for _, expected := range enumValues_NodeinfoVersion {
		if reflect.DeepEqual(v, expected) {
			ok = true
			break
		}
	}
	if !ok {
		return fmt.Errorf("invalid value (expected one of %#v): %#v", enumValues_NodeinfoVersion, v)
	}
	*j = NodeinfoVersion(v)
	return nil
}

// NodeInfo schema version 2.0.
type Nodeinfo struct {
	// Free form key value pairs for software specific values. Clients should not rely
	// on any specific key present.
	Metadata NodeinfoMetadata `json:"metadata" yaml:"metadata" mapstructure:"metadata"`

	// Whether this server allows open self-registration.
	OpenRegistrations bool `json:"openRegistrations" yaml:"openRegistrations" mapstructure:"openRegistrations"`

	// The protocols supported on this server.
	Protocols []NodeinfoProtocolsElem `json:"protocols" yaml:"protocols" mapstructure:"protocols"`

	// The third party sites this server can connect to via their application API.
	Services NodeinfoServices `json:"services" yaml:"services" mapstructure:"services"`

	// Metadata about server software in use.
	Software NodeinfoSoftware `json:"software" yaml:"software" mapstructure:"software"`

	// Usage statistics for this server.
	Usage NodeinfoUsage `json:"usage" yaml:"usage" mapstructure:"usage"`

	// The schema version, must be 2.0.
	Version NodeinfoVersion `json:"version" yaml:"version" mapstructure:"version"`
}

const NodeinfoVersionA20 NodeinfoVersion = "2.0"

// UnmarshalJSON implements json.Unmarshaler.
func (j *Nodeinfo) UnmarshalJSON(b []byte) error {
	var raw map[string]interface{}
	if err := json.Unmarshal(b, &raw); err != nil {
		return err
	}
	if v, ok := raw["metadata"]; !ok || v == nil {
		return fmt.Errorf("field metadata in Nodeinfo: required")
	}
	if v, ok := raw["openRegistrations"]; !ok || v == nil {
		return fmt.Errorf("field openRegistrations in Nodeinfo: required")
	}
	if v, ok := raw["protocols"]; !ok || v == nil {
		return fmt.Errorf("field protocols in Nodeinfo: required")
	}
	if v, ok := raw["services"]; !ok || v == nil {
		return fmt.Errorf("field services in Nodeinfo: required")
	}
	if v, ok := raw["software"]; !ok || v == nil {
		return fmt.Errorf("field software in Nodeinfo: required")
	}
	if v, ok := raw["usage"]; !ok || v == nil {
		return fmt.Errorf("field usage in Nodeinfo: required")
	}
	if v, ok := raw["version"]; !ok || v == nil {
		return fmt.Errorf("field version in Nodeinfo: required")
	}
	type Plain Nodeinfo
	var plain Plain
	if err := json.Unmarshal(b, &plain); err != nil {
		return err
	}
	if len(plain.Protocols) < 1 {
		return fmt.Errorf("field %s length: must be >= %d", "protocols", 1)
	}
	*j = Nodeinfo(plain)
	return nil
}