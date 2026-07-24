package mapper

import (
	"github.com/jinzhu/copier"
)

type EnumTypeConverter[DTO ~int32, ENTITY ~string] struct {
	nameMap  map[int32]string
	valueMap map[string]int32
}

// NewEnumTypeConverter 创建枚举双向转换器。
func NewEnumTypeConverter[DTO ~int32, ENTITY ~string](
	nameMap map[int32]string,
	valueMap map[string]int32,
) *EnumTypeConverter[DTO, ENTITY] {
	return &EnumTypeConverter[DTO, ENTITY]{
		valueMap: valueMap,
		nameMap:  nameMap,
	}
}

// ToEntity 将 DTO 枚举值转换为实体枚举值。
func (m *EnumTypeConverter[DTO, ENTITY]) ToEntity(dto *DTO) *ENTITY {
	if dto == nil {
		return nil
	}

	find, ok := m.nameMap[int32(*dto)]
	if !ok {
		return nil
	}

	return new(ENTITY(find))
}

// ToDTO 将实体枚举值转换为 DTO 枚举值。
func (m *EnumTypeConverter[DTO, ENTITY]) ToDTO(entity *ENTITY) *DTO {
	if entity == nil {
		return nil
	}

	find, ok := m.valueMap[string(*entity)]
	if !ok {
		return nil
	}

	return new(DTO(find))
}

// NewConverterPair 创建枚举类型的双向转换器。
func (m *EnumTypeConverter[DTO, ENTITY]) NewConverterPair() []copier.TypeConverter {
	fromFn := m.ToDTO
	toFn := m.ToEntity

	return NewGenericTypeConverterPair(new(ENTITY("")), new(DTO(0)), fromFn, toFn)
}

// NewGenericTypeConverterPair 创建通用的双向类型转换器对。
func NewGenericTypeConverterPair[A interface{}, B interface{}](
	srcType A,
	dstType B,
	fromFn func(src A) B,
	toFn func(src B) A,
) []copier.TypeConverter {
	return []copier.TypeConverter{
		{
			SrcType: srcType,
			DstType: dstType,
			Fn: func(src interface{}) (interface{}, error) {
				return fromFn(src.(A)), nil
			},
		},
		{
			SrcType: dstType,
			DstType: srcType,
			Fn: func(src interface{}) (interface{}, error) {
				return toFn(src.(B)), nil
			},
		},
	}
}
