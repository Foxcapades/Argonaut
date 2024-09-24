package cereal

type Deserializer[T any] interface {
	Deserialize(raw string) (T, error)
}
