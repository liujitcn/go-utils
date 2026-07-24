package mapper

import (
	"github.com/jinzhu/copier"
)

// Mapper 定义 DTO 与实体之间的双向转换能力。
type Mapper[DTO any, ENTITY any] interface {
	// ToEntity 将 DTO 转换为实体对象。
	ToEntity(*DTO) *ENTITY

	// ToDTO 将实体对象转换为 DTO。
	ToDTO(*ENTITY) *DTO
}

// CopierMapper 基于 copier 执行 DTO 与实体的双向转换。
type CopierMapper[DTO any, ENTITY any] struct {
	copierOption copier.Option
}

// NewCopierMapper 创建带默认转换器的映射器。
func NewCopierMapper[DTO any, ENTITY any]() *CopierMapper[DTO, ENTITY] {
	mapper := &CopierMapper[DTO, ENTITY]{
		copierOption: copier.Option{
			Converters: []copier.TypeConverter{},
		},
	}
	// 默认注册时间类型转换器，统一时间字段的双向复制行为。
	mapper.AppendConverters(NewTimeTypeConverter().NewConverterPair())
	// 默认注册基本类型的JSON转换
	mapper.AppendConverters(NewJSONTypeConverter[[]float64]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]float32]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]uint64]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]uint32]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]int32]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]int64]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]bool]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[]string]().NewConverterPair())
	mapper.AppendConverters(NewJSONTypeConverter[[][]byte]().NewConverterPair())

	return mapper
}

// AppendConverter 追加单个类型转换器。
func (m *CopierMapper[DTO, ENTITY]) AppendConverter(converter copier.TypeConverter) {
	m.copierOption.Converters = append(m.copierOption.Converters, converter)
}

// AppendConverters 批量追加类型转换器。
func (m *CopierMapper[DTO, ENTITY]) AppendConverters(converters []copier.TypeConverter) {
	m.copierOption.Converters = append(m.copierOption.Converters, converters...)
}

// ToEntity 将 DTO 转换为实体对象。
func (m *CopierMapper[DTO, ENTITY]) ToEntity(dto *DTO) *ENTITY {
	if dto == nil {
		return nil
	}

	var entity ENTITY
	if err := copier.CopyWithOption(&entity, dto, m.copierOption); err != nil {
		panic(err) // 转换失败属于配置或类型映射错误，调用方应在开发阶段修正。
	}

	return &entity
}

// ToDTO 将实体对象转换为 DTO。
func (m *CopierMapper[DTO, ENTITY]) ToDTO(entity *ENTITY) *DTO {
	if entity == nil {
		return nil
	}

	var dto DTO
	if err := copier.CopyWithOption(&dto, entity, m.copierOption); err != nil {
		panic(err) // 转换失败属于配置或类型映射错误，调用方应在开发阶段修正。
	}

	return &dto
}
