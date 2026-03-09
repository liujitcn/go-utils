package _map

func Keys[K comparable, V any](mapInstance map[K]V) []K {
	keys := make([]K, len(mapInstance))

	i := 0
	for k := range mapInstance {
		keys[i] = k
		i++
	}

	return keys
}

func Values[K comparable, V any](mapInstance map[K]V) []V {
	values := make([]V, len(mapInstance))

	i := 0
	for _, v := range mapInstance {
		values[i] = v
		i++
	}

	return values
}

func Merge[K comparable, V any](mapInstances ...map[K]V) map[K]V {
	var mergedMapSize int

	for _, mapInstance := range mapInstances {
		mergedMapSize += len(mapInstance)
	}

	mergedMap := make(map[K]V, mergedMapSize)

	for _, mapInstance := range mapInstances {
		for k, v := range mapInstance {
			mergedMap[k] = v
		}
	}

	return mergedMap
}

func ForEach[K comparable, V any](mapInstance map[K]V, function func(key K, value V)) {
	for key, value := range mapInstance {
		function(key, value)
	}
}

func Drop[K comparable, V any](mapInstance map[K]V, keys []K) map[K]V {
	for _, key := range keys {
		delete(mapInstance, key)
	}

	return mapInstance
}

func Copy[K comparable, V any](mapInstance map[K]V) map[K]V {
	mapCopy := make(map[K]V, len(mapInstance))

	for key, value := range mapInstance {
		mapCopy[key] = value
	}

	return mapCopy
}

func Filter[K comparable, V any](mapInstance map[K]V, function func(key K, value V) bool) map[K]V {
	mapCopy := make(map[K]V, len(mapInstance))

	for key, value := range mapInstance {
		if function(key, value) {
			mapCopy[key] = value
		}
	}

	return mapCopy
}
