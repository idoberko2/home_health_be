package app

type App interface {
	Run()
}

func New() App {
	return &app{}
}

type app struct{}

func (a *app) Run() {

}
