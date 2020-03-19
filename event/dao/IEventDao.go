package dao

type IEventDao interface {
	InitDb()
	FindById(id int64) Event
	Insert(event Event)
	Update(id int64, event Event)
	GetAll() []Event
	DeleteById(id int64)
}
