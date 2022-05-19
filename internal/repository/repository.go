package repository

import "monitor-apiserver/pkg/database"

var _ Repository = (*repository)(nil)

type Repository interface {
	Users() UserRepository
}

type repository struct {
	db database.IDataSource
}

func NewRepository(_database database.IDataSource) Repository {
	return &repository{
		db: _database,
	}
}

func (r repository) Users() UserRepository {
	return newUserRepository(r)
}
