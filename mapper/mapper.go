package mapper

import (
	"time"

	"github.com/jinzhu/copier"
)

const defaultTimeFormat = "2006-01-02 15:04:05"

// Mapper defines the interface for converting between Data Transfer Objects (DTOs) and Database Entities.
type Mapper[DTO any, ENTITY any] interface {
	// ToEntity converts a DTO to a Database Entity.
	ToEntity(*DTO) *ENTITY

	// ToDTO converts a Database Entity to a DTO.
	ToDTO(*ENTITY) *DTO
}

type CopierMapper[DTO any, ENTITY any] struct {
	copierOption copier.Option
}

func NewCopierMapper[DTO any, ENTITY any]() *CopierMapper[DTO, ENTITY] {
	mapper := &CopierMapper[DTO, ENTITY]{
		copierOption: copier.Option{
			Converters: []copier.TypeConverter{},
		},
	}
	mapper.AppendConverters(NewGenericTypeConverterPair(
		time.Time{},
		"",
		func(src time.Time) string {
			return src.Format(defaultTimeFormat)
		},
		func(src string) time.Time {
			timeValue, err := time.ParseInLocation(defaultTimeFormat, src, time.Local)
			if err != nil {
				return time.Time{}
			}
			return timeValue
		},
	))
	return mapper
}

func (m *CopierMapper[DTO, ENTITY]) AppendConverter(converter copier.TypeConverter) {
	m.copierOption.Converters = append(m.copierOption.Converters, converter)
}

func (m *CopierMapper[DTO, ENTITY]) AppendConverters(converters []copier.TypeConverter) {
	m.copierOption.Converters = append(m.copierOption.Converters, converters...)
}

func (m *CopierMapper[DTO, ENTITY]) ToEntity(dto *DTO) *ENTITY {
	if dto == nil {
		return nil
	}

	var entity ENTITY
	if err := copier.CopyWithOption(&entity, dto, m.copierOption); err != nil {
		panic(err) // Handle error appropriately in production code
	}

	return &entity
}

func (m *CopierMapper[DTO, ENTITY]) ToDTO(entity *ENTITY) *DTO {
	if entity == nil {
		return nil
	}

	var dto DTO
	if err := copier.CopyWithOption(&dto, entity, m.copierOption); err != nil {
		panic(err) // Handle error appropriately in production code
	}

	return &dto
}
