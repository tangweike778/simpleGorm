package gorm

func initializeCallbacks(db *DB) *callbacks {
	return &callbacks{
		processors: map[string]*processor{
			"create": {db: db},
		},
	}
}

type callback struct {
	name      string
	before    string
	after     string
	remove    bool
	replace   bool
	match     func(*DB) bool
	handler   func(*DB)
	processor *processor
}

type processor struct {
	db        *DB
	Clauses   []string
	fns       []func(*DB)
	callbacks []*callback
}

type callbacks struct {
	processors map[string]*processor
}

func (cs *callbacks) Create() *processor {
	return cs.processors["create"]
}

func (p *processor) Register(name string, fn func(*DB)) error {
	return (&callback{processor: p}).Register(name, fn)
}

func (p *processor) complie() error {

}

func (c *callback) Register(name string, fn func(*DB)) error {
	c.name = name
	c.handler = fn
	c.processor.callbacks = append(c.processor.callbacks, c)
	return c.processor.complie()
}
