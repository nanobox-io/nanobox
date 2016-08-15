package models

func Inspect(bucket, key string) interface{} {
	if key == "" {
		v := []interface{}{}
		getAll(bucket, &v)
		return v
	}
	var v interface{}
	get(bucket, key, &v)
	return v
}
