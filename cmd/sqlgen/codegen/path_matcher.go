package codegen

import (
	"io/fs"
	"os"
	"path/filepath"
	"regexp"
	"strings"

	"github.com/si3nloong/sqlgen/cmd/sqlgen/internal/fileutil"
)

type Matcher interface {
	Match(v string) bool
}

type EmptyMatcher struct{}

func (*EmptyMatcher) Match(v string) bool {
	return true
}

type FileMatcher map[string]struct{}

func (f FileMatcher) Match(v string) bool {
	_, ok := f[v]
	return ok
}

type FolderMatcher string

func (f FolderMatcher) Match(v string) bool {
	return strings.HasPrefix(v, (string)(f))
}

type RegexMatcher struct {
	*regexp.Regexp
}

func (r *RegexMatcher) Match(v string) bool {
	return r.MatchString(v)
}

func PathResolver(path string) (Matcher, error) {
	if path == "." {
		path = fileutil.Getpwd()
		// If the prefix is ".", mean it's refer to current directory
	} else if path[0] == '.' {
		path = fileutil.Getpwd() + path[1:]
	} else if path[0] != '/' {
		path = filepath.Join(fileutil.Getpwd(), path)
	}

	// If suffix is *, we will add go extension to it
	if path[len(path)-1] == '*' {
		path = path + ".go"
	}

	r := regexp.MustCompile(`(?i)((?:\/)([a-z][a-z0-9-_.]+\/)*)\w*\*\w*(?:\.go)`)
	submatches := r.FindStringSubmatch(path)
	if strings.Contains(path, "**") {
		paths := strings.SplitN(path, "**", 2)
		rootDir := strings.TrimSuffix(strings.TrimSpace(paths[0]), "/")
		suffix := `(?:[\\/]\w+\.\w+)`
		if paths[1] != "" {
			suffix = path2Regex.Replace(paths[1])
		}
		if err := filepath.WalkDir(rootDir, func(path string, d fs.DirEntry, err error) error {
			// If the directory is not exists, the "d" will be nil
			if d == nil || !d.IsDir() {
				// If it's not a folder, we skip!
				return nil
			}
			// dirs = append(dirs, strings.TrimPrefix(path, rootDir))
			return nil
		}); err != nil {
			// return fmt.Errorf(`sqlgen: failed to walk schema %s: %w`, paths[0], err)
		}
		return &RegexMatcher{regexp.MustCompile(path2Regex.Replace(rootDir) + `([\\/][a-z0-9_-]+)*` + suffix)}, nil
	} else if len(submatches) > 0 {
		// rootDir = strings.TrimSuffix(submatches[1], "/")
		// dirs = append(dirs, "")
		// slog.Info("Submatch", "rootDir", rootDir, "dir", path2Regex.Replace(path))
		return &RegexMatcher{regexp.MustCompile(path2Regex.Replace(path))}, nil
	} else {
		fi, err := os.Stat(path)
		// If the file or folder not exists, we skip!
		if os.IsNotExist(err) {
			return nil, err
			// goto nextSrc
		} else if err != nil {
			return nil, err
		}

		if fi.IsDir() {
			// If it's just a folder
			return FolderMatcher(path), nil
		}
		// If it's just a file
		return FileMatcher{filepath.Join(filepath.Dir(path), fi.Name()): struct{}{}}, nil

		// rootDir = path
		// dirs = append(dirs, "")
	}
}
