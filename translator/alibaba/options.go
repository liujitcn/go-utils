package alibaba

// Option 配置阿里云翻译器。
type Option func(*Translator)

// WithRegionID 设置阿里云区域，默认使用 cn-hangzhou。
func WithRegionID(regionID string) Option {
	return func(translator *Translator) {
		translator.regionID = regionID
	}
}
