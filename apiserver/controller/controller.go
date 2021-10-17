package controller

type Controller interface {
	List()
	Create()
	Get()
	Update()
	Delete()
}

type ModelController interface {
	Name() string
	Controller() Controller
}

type modelController struct {
	name       string
	controller Controller
}

func NewModelController(name string, controller Controller) ModelController {
	c := &modelController{name: name, controller: controller}
	return c
}

func (c *modelController) Name() string {
	return c.name
}

func (c *modelController) Controller() Controller {
	return c.controller
}
