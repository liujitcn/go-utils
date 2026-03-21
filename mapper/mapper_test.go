package mapper

import (
	"testing"
	"time"

	"github.com/stretchr/testify/assert"
)

func TestCopierMapper(t *testing.T) {
	type DtoType struct {
		Name string
		Age  int
	}

	type EntityType struct {
		Name string
		Age  int
	}

	mapper := NewCopierMapper[DtoType, EntityType]()

	// 测试 ToEntity 方法
	dto := &DtoType{Name: "Alice", Age: 25}
	entity := mapper.ToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, "Alice", entity.Name)
	assert.Equal(t, 25, entity.Age)

	// 测试 ToEntity 方法，传入 nil
	entityNil := mapper.ToEntity(nil)
	assert.Nil(t, entityNil)

	// 测试 ToDTO 方法
	entity = &EntityType{Name: "Bob", Age: 30}
	dtoResult := mapper.ToDTO(entity)
	assert.NotNil(t, dtoResult)
	assert.Equal(t, "Bob", dtoResult.Name)
	assert.Equal(t, 30, dtoResult.Age)

	// 测试 ToDTO 方法，传入 nil
	dtoNil := mapper.ToDTO(nil)
	assert.Nil(t, dtoNil)
}

func TestEnumTypeConverter(t *testing.T) {
	type DtoType int32
	type EntityType string

	const (
		DtoTypeOne DtoType = 1
		DtoTypeTwo DtoType = 2
	)

	const (
		EntityTypeOne EntityType = "One"
		EntityTypeTwo EntityType = "Two"
	)

	nameMap := map[int32]string{
		1: "One",
		2: "Two",
	}
	valueMap := map[string]int32{
		"One": 1,
		"Two": 2,
	}

	converter := NewEnumTypeConverter[DtoType, EntityType](nameMap, valueMap)

	// 测试 ToEntity 方法
	entity := converter.ToEntity(new(DtoTypeOne))
	assert.NotNil(t, entity)
	assert.Equal(t, "One", string(*entity))

	// 测试 ToEntity 方法，传入不存在的值
	entityInvalid := converter.ToEntity(new(DtoType(3)))
	assert.Nil(t, entityInvalid)

	// 测试 ToDTO 方法
	entity = new(EntityTypeTwo)
	dtoResult := converter.ToDTO(entity)
	assert.NotNil(t, dtoResult)
	assert.Equal(t, DtoType(2), *dtoResult)

	// 测试 ToDTO 方法，传入不存在的值
	entityInvalid = new(EntityType("Three"))
	dtoInvalidResult := converter.ToDTO(entityInvalid)
	assert.Nil(t, dtoInvalidResult)
}

func TestCopierMapper_TimeConverter(t *testing.T) {
	type DtoType struct {
		CreatedAt string
	}

	type EntityType struct {
		CreatedAt time.Time
	}

	mapper := NewCopierMapper[DtoType, EntityType]()
	entity := &EntityType{
		CreatedAt: time.Date(2024, 1, 2, 3, 4, 5, 0, time.Local),
	}

	dto := mapper.ToDTO(entity)
	assert.NotNil(t, dto)
	assert.NotEmpty(t, dto.CreatedAt)

	entityResult := mapper.ToEntity(dto)
	assert.NotNil(t, entityResult)
	assert.False(t, entityResult.CreatedAt.IsZero())
}

func TestJSONTypeConverter(t *testing.T) {
	type Profile struct {
		Name string `json:"name"`
		Age  int    `json:"age"`
	}

	converter := NewJSONTypeConverter[Profile]()

	jsonValue := `{"name":"alice","age":18}`
	entity := converter.ToEntity(&jsonValue)
	assert.NotNil(t, entity)
	assert.Equal(t, "alice", entity.Name)
	assert.Equal(t, 18, entity.Age)

	dto := converter.ToDTO(entity)
	assert.NotNil(t, dto)
	assert.JSONEq(t, jsonValue, *dto)

	invalidJSON := "{"
	entityInvalid := converter.ToEntity(&invalidJSON)
	assert.Nil(t, entityInvalid)
}

func TestTimeTypeConverter(t *testing.T) {
	converter := NewTimeTypeConverter()
	now := time.Date(2024, 1, 2, 3, 4, 5, 0, time.Local)

	dto := converter.ToDTO(&now)
	assert.NotNil(t, dto)
	assert.NotEmpty(t, *dto)

	entity := converter.ToEntity(dto)
	assert.NotNil(t, entity)
	assert.Equal(t, now.Unix(), entity.Unix())

	invalidValue := "invalid"
	entityInvalid := converter.ToEntity(&invalidValue)
	assert.Nil(t, entityInvalid)
}

func TestInt64ArrayTypeConverter(t *testing.T) {
	converter := NewInt64ArrayTypeConverter()
	value := `[1,2,3]`

	entity := converter.ToEntity(&value)
	assert.NotNil(t, entity)
	assert.Equal(t, []int64{1, 2, 3}, *entity)

	dto := converter.ToDTO(entity)
	assert.NotNil(t, dto)
	assert.JSONEq(t, value, *dto)

	emptyJSON := ""
	entityEmpty := converter.ToEntity(&emptyJSON)
	assert.NotNil(t, entityEmpty)
	assert.Empty(t, *entityEmpty)
}

func TestStringArrayTypeConverter(t *testing.T) {
	converter := NewStringArrayTypeConverter()
	value := `["a","b","c"]`

	entity := converter.ToEntity(&value)
	assert.NotNil(t, entity)
	assert.Equal(t, []string{"a", "b", "c"}, *entity)

	dto := converter.ToDTO(entity)
	assert.NotNil(t, dto)
	assert.JSONEq(t, value, *dto)

	emptyJSON := ""
	entityEmpty := converter.ToEntity(&emptyJSON)
	assert.NotNil(t, entityEmpty)
	assert.Empty(t, *entityEmpty)
}

func TestCopierMapper_ArrayConverter(t *testing.T) {
	type DtoType struct {
		RoleIDs string
		Tags    string
	}

	type EntityType struct {
		RoleIDs []int64
		Tags    []string
	}

	mapper := NewCopierMapper[DtoType, EntityType]()
	mapper.AppendConverters(NewInt64ArrayTypeConverter().NewConverterPair())
	mapper.AppendConverters(NewStringArrayTypeConverter().NewConverterPair())

	entity := &EntityType{
		RoleIDs: []int64{1, 2, 3},
		Tags:    []string{"a", "b"},
	}

	dto := mapper.ToDTO(entity)
	assert.NotNil(t, dto)
	assert.JSONEq(t, `[1,2,3]`, dto.RoleIDs)
	assert.JSONEq(t, `["a","b"]`, dto.Tags)

	entityResult := mapper.ToEntity(dto)
	assert.NotNil(t, entityResult)
	assert.Equal(t, entity.RoleIDs, entityResult.RoleIDs)
	assert.Equal(t, entity.Tags, entityResult.Tags)
}
