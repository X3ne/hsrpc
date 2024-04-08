//go:generate go run generate.go

// credits: https://github.com/koddr/example-embed-static-files-go

package bundle

type embedBundle struct {
	storage map[string][]byte
}

func newEmbedBundle() *embedBundle {
	return &embedBundle{storage: make(map[string][]byte)}
}

func (e *embedBundle) Add(file string, content []byte) {
	e.storage[file] = content
}

func (e *embedBundle) Get(file string) []byte {
	if f, ok := e.storage[file]; ok {
		return f
	}
	return nil
}

var box = newEmbedBundle()

func Add(file string, content []byte) {
	box.Add(file, content)
}

func Get(file string) []byte {
	return box.Get(file)
}
