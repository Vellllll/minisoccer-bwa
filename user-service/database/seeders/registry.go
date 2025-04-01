package seeders

import "gorm.io/gorm"

type Registry struct {
	db *gorm.DB
}

type InterfaceSeederRegistry interface {
	Run()
}

func NewSeederRegistry(db *gorm.DB) InterfaceSeederRegistry {
	return &Registry{db: db}
}

func (s *Registry) Run() {
	SeedRoles(s.db)
	SeedUsers(s.db)
}
