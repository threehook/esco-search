package resolvers

const (
	allowedFileContentType = "text/plain; charset=utf-8" // mime type of json is text/plain
)

//nolint:gochecknoglobals
var (
	errContentType  = "contenttype niet toegestaan, toegestane contenttype: application/json"
	errOrgEmptyKvKs = "minimaal één kvk nummer verplicht"
	errInternal     = "interne fout"
)
