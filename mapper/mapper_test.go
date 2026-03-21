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
