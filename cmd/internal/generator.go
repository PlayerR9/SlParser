package internal

// GenData is the generator data.
type GenData struct {
	// PackageName is the package name.
	PackageName string
}

// SetPackageName implements generator.PackageNameSetter interface.
func (gd *GenData) SetPackageName(name string) bool {
	if gd == nil {
		return false
	}

	gd.PackageName = name

	return true
}
