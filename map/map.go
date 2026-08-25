package _map

import "maps"

// Merge 合并多个 map，后出现的键会覆盖前面的值。
func Merge[K comparable, V any](mapInstances ...map[K]V) map[K]V {
	var mergedMapSize int

	for _, mapInstance := range mapInstances {
		mergedMapSize += len(mapInstance)
	}

	mergedMap := make(map[K]V, mergedMapSize)

	for _, mapInstance := range mapInstances {
		maps.Copy(mergedMap, mapInstance)
	}

	return mergedMap
}

// ForEach 遍历 map 中的每个键值对。
func ForEach[K comparable, V any](mapInstance map[K]V, function func(key K, value V)) {
	for key, value := range mapInstance {
		function(key, value)
	}
}

// Drop 从 map 中删除指定键并返回原 map。
func Drop[K comparable, V any](mapInstance map[K]V, keys []K) map[K]V {
	for _, key := range keys {
		delete(mapInstance, key)
	}

	return mapInstance
}

// Filter 按条件过滤 map 中的键值对。
func Filter[K comparable, V any](mapInstance map[K]V, function func(key K, value V) bool) map[K]V {
	mapCopy := make(map[K]V, len(mapInstance))

	for key, value := range mapInstance {
		if function(key, value) {
			mapCopy[key] = value
		}
	}

	return mapCopy
}
