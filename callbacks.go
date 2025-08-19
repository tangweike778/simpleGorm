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
	replace   bool
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
	var err error
	if p.fns, err = sortCallbacks(p.callbacks); err != nil {
		return err
	}
	return nil
}

func sortCallbacks(c []*callback) (fns []func(*DB), err error) {
	for _, item := range c {
		fns = append(fns, item.handler)
	}
	return
}

func (c *callback) Register(name string, fn func(*DB)) error {
	c.name = name
	c.handler = fn
	c.processor.callbacks = append(c.processor.callbacks, c)
	return c.processor.complie()
}
