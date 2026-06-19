package property

// ini kontrak ke databasee
// Repository kontrak yang dibutuhin service buat nyentuh data properti dan implementasinya di folder repository, bukan di sini.
type Repository interface {
	Create(p *Property) (*Property, error)
	FindAll(f Filter) ([]Property, error)
	FindByID(id string) (*Property, error)
	Update(p *Property) (*Property, error)
	Delete(id string) error
	FindByOwner(userID string) ([]Property, error)
}
