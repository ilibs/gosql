package gosql

type ModelWrapper struct {
	dbList map[string]*DB
	model  interface{}
}

type ModelWrapperFactory func(m interface{}) *ModelWrapper

func NewModelWrapper(dbList map[string]*DB, model interface{}) *ModelWrapper {
	return &ModelWrapper{dbList: dbList, model: model}
}

func (m *ModelWrapper) GetRelationDB(connect string) *DB {
	return m.dbList[connect]
}

func (m *ModelWrapper) UnWrap() interface{} {
	return m.model
}
