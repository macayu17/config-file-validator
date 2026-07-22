package filetype

import (
	"testing"

	"github.com/stretchr/testify/require"

	"github.com/Boeing/config-file-validator/v3/pkg/formatter"
	"github.com/Boeing/config-file-validator/v3/pkg/formatter/jsonfmt"
)

// TestJSONFormatterHandlesComments verifies that .json files carrying comments
// or trailing commas are formatted through the JSONC fallback instead of being
// skipped, and that genuinely broken files still error.
func TestJSONFormatterHandlesComments(t *testing.T) {
	t.Parallel()
	cases := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "plain json",
			src:  "{\n\t\"name\": \"my-app\"\n}\n",
			want: "{\n  \"name\": \"my-app\"\n}\n",
		},
		{
			name: "line comment",
			src:  "{\n\t// name of the app\n\t\"name\": \"my-app\"\n}\n",
			want: "{\n  // name of the app\n  \"name\": \"my-app\"\n}\n",
		},
		{
			name: "trailing comma dropped",
			src:  "{\n\t\"name\": \"my-app\",\n}\n",
			want: "{\n  \"name\": \"my-app\"\n}\n",
		},
	}

	f := jsonWithComments{}
	opts := jsonfmt.DefaultOptions()
	opts.LineEnding = formatter.LineEndingLF

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := f.Format([]byte(tc.src), opts)
			require.NoError(t, err)
			require.Equal(t, tc.want, string(got))
		})
	}

	t.Run("broken json errors", func(t *testing.T) {
		t.Parallel()
		_, err := f.Format([]byte(`{"key": "value"`), opts)
		require.Error(t, err)
	})
}

// TestJSONFormatterRespectsTabs verifies that an explicit tab indent style
// still wins over the space default, for both the strict and fallback paths.
func TestJSONFormatterRespectsTabs(t *testing.T) {
	t.Parallel()
	opts := jsonfmt.DefaultOptions()
	opts.IndentStyle = formatter.IndentTabs
	opts.LineEnding = formatter.LineEndingLF

	cases := []struct {
		name string
		src  string
		want string
	}{
		{
			name: "plain json",
			src:  "{\n  \"name\": \"my-app\"\n}\n",
			want: "{\n\t\"name\": \"my-app\"\n}\n",
		},
		{
			name: "line comment",
			src:  "{\n  // c\n  \"name\": \"my-app\"\n}\n",
			want: "{\n\t// c\n\t\"name\": \"my-app\"\n}\n",
		},
	}

	f := jsonWithComments{}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			t.Parallel()
			got, err := f.Format([]byte(tc.src), opts)
			require.NoError(t, err)
			require.Equal(t, tc.want, string(got))
		})
	}
}
