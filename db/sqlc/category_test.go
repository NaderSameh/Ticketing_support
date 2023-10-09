package db

import (
	"context"
	"database/sql"
	"testing"

	"github.com/stretchr/testify/require"
)

func createRandomCategory() (cat Category, err error) {

	return testQueries.CreateCategory(context.Background(), "hardware bugs")

}

func TestCreateCategory(t *testing.T) {
	category, err := testQueries.CreateCategory(context.Background(), "random")

	require.NoError(t, err)
	require.NotZero(t, category.CategoryID)
	require.Equal(t, category.Name, "random")

}
func TestDeleteCategory(t *testing.T) {
	category, _ := createRandomCategory()

	err := testQueries.DeleteCategory(context.Background(), category.CategoryID)
	require.NoError(t, err)

	category1, err := testQueries.GetCategory(context.Background(), category.CategoryID)
	require.Error(t, err)
	require.EqualError(t, err, sql.ErrNoRows.Error())
	require.Empty(t, category1)

}

func TestGetCategory(t *testing.T) {
	category, _ := createRandomCategory()

	category1, err := testQueries.GetCategory(context.Background(), category.CategoryID)
	require.NoError(t, err)
	require.Equal(t, category1.Name, category.Name)

}

func TestListCategories(t *testing.T) {
	for i := 0; i < 10; i++ {
		createRandomCategory()
	}

	args := ListCategoriesParams{
		Limit:  5,
		Offset: 0,
	}

	categories, err := testQueries.ListCategories(context.Background(), args)
	require.NoError(t, err)
	require.NotEmpty(t, categories)
}
