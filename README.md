# config

Go package for app configuration. Supports chained configuration sources for multiple levels of defaults.
Includes APIs for loading Linux style configuration files (name/value pairs) or INI files, map based properties,
or easily create new configuration sources (e.g. load from database).

Support monitoring configuration sources for changes, hot loading properties, and notifying listeners of changes.

## Usage

```Go
cfg := &config.Config{}
defer cfg.Shutdown() // stops monitoring

// load file via filespec string, os.File
src, err := config.NewSrcFileFromFilespec("./myfile.conf")
if err != nil {
    return err
}
// add src to top of chain, meaning first searched
cfg.PrependSource(src)

// fetch prop 'retries', default to 3 if not found
val := cfg.Int("retries", 3)
```

