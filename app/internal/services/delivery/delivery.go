package delivery

type Updater interface {
	Update(args interface{}) error
}
