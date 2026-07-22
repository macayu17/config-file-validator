package filetype

import (
	"github.com/Boeing/config-file-validator/v3/pkg/formatter"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/envfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/hclfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/inifmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/jsoncfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/jsonfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/propfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/tomlfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/xmlfmt"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/yamlfmt"
)

// init registers formatters with their corresponding FileTypes. This must run
// after the main init() in file_type.go that builds the FileTypes slice (Go
// processes init() functions in filename-sorted order within a package).
//
// We update the slice entries directly because FileTypes holds value copies —
// updating the package-level vars (JSONFileType etc.) has no effect on the
// already-copied slice.
func init() {
	for i, ft := range FileTypes {
		switch ft.Name {
		case "json":
			FileTypes[i].Formatter = jsonWithComments{}
		case "jsonc":
			FileTypes[i].Formatter = jsoncfmt.Formatter{}
		case "yaml":
			FileTypes[i].Formatter = yamlfmt.Formatter{}
		case "hcl":
			FileTypes[i].Formatter = hclfmt.Formatter{}
		case "xml":
			FileTypes[i].Formatter = xmlfmt.Formatter{}
		case "toml":
			FileTypes[i].Formatter = tomlfmt.Formatter{}
		case "ini":
			FileTypes[i].Formatter = inifmt.Formatter{}
		case "env":
			FileTypes[i].Formatter = envfmt.Formatter{}
		case "properties":
			FileTypes[i].Formatter = propfmt.Formatter{}
		default:
			// no formatter registered for this type yet
		}
	}
}

// jsonWithComments formats .json files, falling back to the JSONC formatter
// when the file uses comments or trailing commas. Editors and toolchains write
// those into .json files (.vscode/settings.json, tsconfig.json), and a strict
// JSON formatter can only skip them — leaving the file with, for example, its
// original tab indentation. Trailing commas are dropped so the fallback still
// produces valid JSON.
type jsonWithComments struct{}

var _ formatter.Formatter = jsonWithComments{}

func (jsonWithComments) Format(src []byte, opts formatter.Options) ([]byte, error) {
	out, err := jsonfmt.Formatter{}.Format(src, opts)
	if err == nil {
		return out, nil
	}

	opts.TrailingCommas = formatter.TrailingCommasNone
	lenient, lenientErr := jsoncfmt.Formatter{}.Format(src, opts)
	if lenientErr != nil {
		// Not JSONC either — the file is a syntax error, so report it as the
		// JSON error the caller expects.
		return nil, err
	}
	return lenient, nil
}
