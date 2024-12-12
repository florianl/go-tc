package tc

// Logger interface only requires a Printf method for now as that is the only method used in the package.
type Logger interface {
	Printf(format string, v ...interface{})
}
