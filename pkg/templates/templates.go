package templates

import (
	"context"
	"io"
	"strings"

	"github.com/nikoksr/simplog"

	"github.com/pkg/errors"
	"github.com/valyala/fasttemplate"
)

// MissingKeyFn defines the behavior of a function that is called when a key is not found in the template.
type MissingKeyFn func(key string) (string, error)

var defaultMissingKeyFn = func(key string) (string, error) {
	return "", errors.Errorf("value for template key %q missing", key)
}

// TemplateEngine is a template engine that can be used to render templates. It uses fasttemplate to parse and render
// templates. Start- and end-tags are used to mark keys. The default start- and end-tags are '{{' and '}}'.
// MissingKeyFn is a function that is called when a key is not found in the template. It is expected to return
// the value for the key.
type TemplateEngine struct {
	StartTag, EndTag string
	MissingKeyFn     MissingKeyFn
}

const (
	defaultStartTag = "%{{"
	defaultEndTag   = "}}%"
)

// NewEngine creates a new template engine.
func NewEngine(startTag, endTag string) *TemplateEngine {
	if startTag == "" {
		startTag = defaultStartTag
	}
	if endTag == "" {
		endTag = defaultEndTag
	}

	return &TemplateEngine{
		StartTag:     startTag,
		EndTag:       endTag,
		MissingKeyFn: defaultMissingKeyFn,
	}
}

// normalizeKey normalizes a key in order to avoid as many duplicate key entries as possible.
// For example, we don't want the map of replaced placeholders to hold an entries for 'project-name', 'projectname' and
// 'Project-Name'. Not only would this result in unnecessarily allocated memory but also in duplicate user input prompts;
// asking three times for the project name would be pretty annoying.
func normalizeKey(key string) string {
	key = strings.ReplaceAll(key, " ", "")
	key = strings.ReplaceAll(key, "-", "")
	key = strings.ReplaceAll(key, "_", "")
	key = strings.ToLower(key)

	return key
}

func humanReadableKey(key string) string {
	key = strings.TrimSpace(key)
	key = strings.ReplaceAll(key, "-", " ")
	key = strings.ReplaceAll(key, "_", " ")
	key = strings.Title(key)

	return key
}

// render the template and writes the result to the writer.
func (t *TemplateEngine) render(ctx context.Context, w io.Writer, tmpl *fasttemplate.Template) error {
	logger := simplog.FromContext(ctx)

	var err error
	keys := make(map[string]string)

	written, err := tmpl.ExecuteFunc(w, func(w io.Writer, key string) (int, error) {
		printableKey := humanReadableKey(key)
		key = normalizeKey(key)

		logger.Debugf("checking value for template key: %q", key)
		value, exists := keys[key]
		if !exists {
			logger.Debugf("value for template key %q not previously defined", key)
			if value, err = t.MissingKeyFn(printableKey); err != nil {
				return 0, err
			}
		}

		logger.Debugf("using value %q for template key %q", value, key)
		keys[key] = value

		return w.Write([]byte(value))
	})

	logger.Debugf("wrote %d bytes to template", written)
	if err != nil {
		return err
	}

	return nil
}

// parse the template file and renders it to the writer.
func (t *TemplateEngine) parse(ctx context.Context, w io.Writer, data []byte) error {
	logger := simplog.FromContext(ctx)

	if t.MissingKeyFn == nil {
		logger.Debugf("using default missing key function")
		t.MissingKeyFn = defaultMissingKeyFn
	}

	// Parse the template
	logger.Debug("parsing template")
	tmpl, err := fasttemplate.NewTemplate(string(data), t.StartTag, t.EndTag)
	if err != nil {
		return errors.Wrap(err, "parse template")
	}

	// Render the template
	logger.Debugf("rendering template")

	return t.render(ctx, w, tmpl)
}

// Parse the template file and renders it to the writer. If the template file is not found, an error is returned. If the
// template file is found but cannot be parsed, an error is returned. If the template file is found and parsed, the
// template is rendered to the writer.
func (t *TemplateEngine) Parse(ctx context.Context, w io.Writer, data []byte) error {
	return t.parse(ctx, w, data)
}

// ParseString is similar to Parse but accepts a string as input instead of a byte slice. It returns an error if the
// template cannot be parsed or rendered. If the template is parsed and rendered successfully, the result is returned.
func (t *TemplateEngine) ParseString(ctx context.Context, data string) (string, error) {
	var b strings.Builder
	if err := t.parse(ctx, &b, []byte(data)); err != nil {
		return "", err
	}

	return b.String(), nil
}
