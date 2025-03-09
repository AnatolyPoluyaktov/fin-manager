package domain

type Category struct {
	ID          int
	Name        string
	Description string
}

type NewCategoryData struct {
	Name        string
	Description string
}

func NewCategory(data NewCategoryData) (Category, error) {
	return Category{
		Name:        data.Name,
		Description: data.Description,
	}, nil
}
