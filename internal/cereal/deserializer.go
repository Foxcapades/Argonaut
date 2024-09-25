package cereal

type Deserializer[T any] interface {
	Deserialize(raw string, prev *T) (T, error)
}
