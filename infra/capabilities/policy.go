package capabilities

// Capability identifies an abstract infra capability.
// The kit does not define debug-surface capabilities; services implement their own policy.
type Capability string

// Allowed returns whether the given capability is allowed by infra policy.
// Default is deny (fail-closed). The kit does not define any capabilities;
// this function exists for future use by services that add their own capability types.
func Allowed(_ Config, _ Capability) bool {
	return false
}
