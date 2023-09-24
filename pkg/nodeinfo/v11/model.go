package v11

func (n *Nodeinfo) SchemaVersion() string {
	return "1.1"
}

func (n *Nodeinfo) SoftwareName() string {
	return string(n.Software.Name)
}

func (n *Nodeinfo) SoftwareVersion() string {
	return n.Software.Version
}

func (n *Nodeinfo) IsRegistrationOpen() bool {
	return n.OpenRegistrations
}

func (n *Nodeinfo) TotalUsers() *int {
	return n.Usage.Users.Total
}

func (n *Nodeinfo) ActiveUsersHalfyear() *int {
	return n.Usage.Users.ActiveHalfyear
}

func (n *Nodeinfo) ActiveUsersMonth() *int {
	return n.Usage.Users.ActiveMonth
}

func (n *Nodeinfo) LocalPosts() *int {
	return n.Usage.LocalPosts
}

func (n *Nodeinfo) LocalComments() *int {
	return n.Usage.LocalComments
}
