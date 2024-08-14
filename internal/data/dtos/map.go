package dtos

// MapFunc provides a generic function for mapping data types from one to another.
// Utilizing this type allows the user to define a single mapping function and then
// utilize receiver functions to easily map slices of the types, or simplify returns
// when an error is present.
type MapFunc[T any, U any] func(T) U

func (a MapFunc[T, U]) Map(v T) U {
	return a(v)
}

func (a MapFunc[T, U]) Slice(v []T) []U {
	result := make([]U, len(v))
	for i := range v {
		result[i] = a(v[i])
	}
	return result
}

func (a MapFunc[T, U]) Err(v T, err error) (U, error) {
	if err != nil {
		var zero U
		return zero, err
	}

	return a(v), nil
}

func (a MapFunc[T, U]) SliceErr(v []T, err error) ([]U, error) {
	if err != nil {
		return nil, err
	}

	return a.Slice(v), nil
}
