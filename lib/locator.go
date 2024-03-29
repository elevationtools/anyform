
package anyform

type Locator struct {
	Config* AnyformConfig
	GomplateRunner GomplateRunner
	JsonnetRunner JsonnetRunner 
}

func NewDefaultLocator() *Locator {
	loc := &Locator{}
	loc.Config = NewDefaultAnyformConfig()
	loc.GomplateRunner = NewCliGomplateRunner(loc)
	loc.JsonnetRunner = NewCliJsonnetRunner(loc)
	return loc
}
