package registry

type Registry struct { 
	Lookup map[string]int
}

func (r *Registry) UpdateName(oldName, newName string) error {
	id, ok := r.Lookup[oldName]
	if !ok {
		return ErrInvalidRegistryKey
	}
	delete(r.Lookup, oldName)

	r.Lookup[newName] = id

	return nil
}
