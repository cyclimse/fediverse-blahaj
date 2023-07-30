package v10

func (n *Nodeinfo) SchemaVersion() string {
	return "1.0"
}

func (n *Nodeinfo) GetSoftwareName() string {
	return string(n.Software.Name)
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