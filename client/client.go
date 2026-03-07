package client

type IClient interface {
	IsRepositoryExist(owner, repository string) (bool, error)
	CreateRepository(owner, repository string) error
}
