package build

// an empty payload function
// this is here so if we redefine 'empty'
// then we can modify it in one place
func emptyPayload() string {
	return "{}"
}
