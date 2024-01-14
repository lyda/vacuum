// Package cui implements the console UI.
package cui

import (
	"fmt"
	"time"

	"github.com/fsnotify/fsnotify"
	"github.com/pb33f/libopenapi/datamodel"
	"github.com/pb33f/libopenapi/index"
	"gopkg.in/yaml.v3"

	"github.com/daveshanley/vacuum/model"
	"github.com/daveshanley/vacuum/motor"
	"github.com/daveshanley/vacuum/shared"
	vacuum_report "github.com/daveshanley/vacuum/vacuum-report"
)

type config struct {
	filename    string
	base        string
	skipCheck   bool
	timeout     int
	hardMode    bool
	follow      bool
	customFuncs map[string]model.RuleFunction
	ruleset     string
}

// File represents an OpenAPI file.
type File struct {
	config    config
	watcher   *fsnotify.Watcher
	resultSet *model.RuleResultSet
	specIndex *index.SpecIndex
	specInfo  *datamodel.SpecInfo
}

// NewFile creates a new OpenAPI file.
func NewFile(filename string, base string, skipCheck bool,
	timeout int, hardMode bool, follow bool,
	customFuncs map[string]model.RuleFunction, ruleset string,
) (*File, error) {

	watcher, err := fsnotify.NewWatcher()
	if err != nil {
		return nil, err
	}
	f := File{
		config: config{
			filename:    filename,
			base:        base,
			skipCheck:   skipCheck,
			timeout:     timeout,
			hardMode:    hardMode,
			follow:      follow,
			customFuncs: customFuncs,
			ruleset:     ruleset,
		},
		watcher: watcher,
	}

	if follow {
		if watcher.Add(filename) != nil {
			return nil, err
		}
	}

	return &f, nil
}

// ReadFile reads in an OpenAPI document.
func (f *File) ReadFile() error {
	var err error
	vacuumReport, specBytes, _ :=
		vacuum_report.BuildVacuumReportFromFile(f.config.filename)
	if len(specBytes) <= 0 {
		return fmt.Errorf("Failed to read specification: %s", f.config.filename)
	}

	// if we have a pre-compiled report, jump straight to the end and collect $500
	if vacuumReport == nil {
		var ruleset *motor.RuleSetExecutionResult
		f.resultSet, ruleset, err = shared.BuildResultsWithDocCheckSkip(
			false,
			f.config.hardMode,
			f.config.ruleset,
			specBytes,
			f.config.customFuncs,
			f.config.base,
			f.config.skipCheck,
			time.Duration(f.config.timeout)*time.Second,
		)
		if err != nil {
			return fmt.Errorf("Failed to render dashboard: %s", err)
		}
		f.specIndex = ruleset.Index
		f.specInfo = ruleset.SpecInfo
		f.specInfo.Generated = time.Now()

	} else {
		f.resultSet = model.NewRuleResultSetPointer(vacuumReport.ResultSet.Results)

		// TODO: refactor dashboard to hold state and rendering as separate entities.
		// dashboard will be slower because it needs an index
		var rootNode yaml.Node
		err = yaml.Unmarshal(*vacuumReport.SpecInfo.SpecBytes, &rootNode)
		if err != nil {
			return fmt.Errorf("Unable to read spec bytes from report file '%s': %s",
				f.config.filename, err)
		}
	}

	return nil
}
