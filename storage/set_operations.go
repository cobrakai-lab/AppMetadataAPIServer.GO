package storage

func IntersectAll( keysFromQuery [][]AppMetadataKey) []AppMetadataKey {
	if len(keysFromQuery) == 0{
		return []AppMetadataKey{}
	}
	result := keysFromQuery[0]
	for _,keys :=range keysFromQuery[1:]{
		result = intersect(result, keys)
	}
	return result
}


func intersect(keys1 []AppMetadataKey, keys2[]AppMetadataKey) []AppMetadataKey{
	result := []AppMetadataKey{}
	if len(keys1) == 0 || len(keys2) == 0{
		return result
	}
	keyMap := make(map[AppMetadataKey]bool)
	for _, key := range keys1{
		keyMap[key] = true
	}

	for _, key:=range keys2{
		if _,found:=keyMap[key];found{
			result = append(result,key )
		}
	}
	return result
}