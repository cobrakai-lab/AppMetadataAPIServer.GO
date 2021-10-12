package storage

type AppMetadataKeySet struct{
	core map[string]map[string]struct{}
	size int
}

func (set *AppMetadataKeySet) Add(key AppMetadataKey) AppMetadataKeySet{
	if set.core==nil{
		set.core = make(map[string]map[string]struct{})
		set.size=0
	}
	if _, found:=set.core[key.Title];found {
		if _, found = set.core[key.Title][key.Version]; !found{
			set.core[key.Title][key.Version] = struct{}{}
			set.size+=1
		}
	}else{
		set.core[key.Title] = make(map[string]struct{})
		set.core[key.Title][key.Version] = struct{}{}
		set.size+=1
	}
	return *set
}

func (set *AppMetadataKeySet) Exists(key AppMetadataKey) bool{
	_, found := set.core[key.Title]
	if found{
		_,found=set.core[key.Title][key.Version]
		return found
	} else{
		return false
	}
}

func (set *AppMetadataKeySet) Delete(key AppMetadataKey){
	_, found := set.core[key.Title]
	if found {
		delete(set.core[key.Title], key.Version)
		set.size-=1
	}
}

func (set *AppMetadataKeySet) Size() int{
	return set.size
}

func (set *AppMetadataKeySet) GetAllAppMetadataKeys() []AppMetadataKey{
	allKeys := []AppMetadataKey{}
	for title,propertyIndex:=range set.core{
		for version, _:= range propertyIndex{
			allKeys= append(allKeys, AppMetadataKey{title, version})
		}
	}
	return allKeys
}
