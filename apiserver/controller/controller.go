package controller

type Controller interface {
	List()
	Create()
	Get()
	Update()
	Delete()
}
