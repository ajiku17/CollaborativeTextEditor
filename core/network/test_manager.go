package network

import (
	"github.com/ajiku17/CollaborativeTextEditor/utils"
)

type DummyManager struct {
	Id utils.UUID
}

func (d DummyManager) GetId() utils.UUID {
	return d.Id
}

func (d DummyManager) Start() {

}

func (d DummyManager) Stop() {

}

func (d DummyManager) Kill() {

}

func NewDummyManager(id string) Manager {
	manager := new (DummyManager)
	manager.Id = utils.UUID(id)
	return manager
}