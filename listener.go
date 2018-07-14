package config

// ChangedListener interface is for receiving notifications
// when one or more properties within any config sources
// have changed values.
type ChangedListener interface {

	// Changed is called when one or more properties in the Source has a
	// changed value.
	Changed(cfg *Config, src *Source)
}

// ChangedPropListener interface is for receiving notifications
// for each property value change.
type ChangedPropListener interface {

	// ChangedProp is called for each property whose value has changed.
	ChangedProp(cfg *Config, src *Source, name string)
}
