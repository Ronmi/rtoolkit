package sam

// Storage abstracts the way you persist current db version
type Storage interface {
	Load(app string) (int, error)
	Save(app string, ver int) error
}
