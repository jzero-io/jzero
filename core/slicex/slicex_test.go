package slicex

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

type Data struct {
	Id   int
	Name string
}

var datas = []Data{
	{Id: 1, Name: "a"},
	{Id: 2, Name: "b"},
	{Id: 3, Name: "c"},
}

func TestPaginate(t *testing.T) {
	paginate := Paginate[Data](datas, 1, 3)
	assert.Equal(t, 3, len(paginate))
	assert.Equal(t, "a", paginate[0].Name)
	assert.Equal(t, "b", paginate[1].Name)
	assert.Equal(t, "c", paginate[2].Name)

	paginate = Paginate[Data](datas, 2, 3)
	assert.Equal(t, 0, len(paginate))
}

func TestToMap(t *testing.T) {
	toMap := ToMap(datas, func(row Data) int {
		return row.Id
	})

	assert.Equal(t, 3, len(toMap))
	assert.Equal(t, "a", toMap[1].Name)
	assert.Equal(t, "b", toMap[2].Name)
	assert.Equal(t, "c", toMap[3].Name)
}
