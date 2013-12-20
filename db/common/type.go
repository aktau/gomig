package common

type TypeModifier uint

const (
	TypeSmall TypeModifier = iota
	TypeNormal
	TypeLarge
	TypeHuge
)

type Type struct {
	/* a base type, possible values:
	 * - text,
	 * - char,
	 * - boolean,
	 * - blob
	 * - float
	 * - double
	 * - numeric (floating point type with scale/precision)
	 * - integer */
	Name string

	/* modifiers that might be important for some RDBMs: small, large, varying... */
	Modifier TypeModifier

	/* parameters that are optionally added to base types */
	Max       uint
	Min       uint
	Scale     uint
	Precision uint
}

func (t *Type) HasMax() bool {
	return t.Max != 0
}

func (t *Type) HasMin() bool {
	return t.Min != 0
}

/* for internal usage (shorter) */
func simple(name string) *Type {
	return &Type{Name: name}
}

func simplem(name string, modifier TypeModifier) *Type {
	return &Type{Name: name, Modifier: modifier}
}

/* for external usage */
func SimpleType(name string) *Type {
	return simple(name)
}

func FloatType() *Type {
	return simple("float")
}

func DoubleType() *Type {
	return simple("double")
}

func NumericType(precision, scale uint) *Type {
	return &Type{Name: "numeric", Precision: precision, Scale: scale}
}

func IntType(modifier TypeModifier) *Type {
	return simplem("integer", modifier)
}

func BlobType() *Type {
	return simple("blob")
}

func BoolType() *Type {
	return simple("boolean")
}

/* padded text */
func PaddedTextType() *Type {
	return simple("char")
}

func TextType() *Type {
	return simple("text")
}

func DateType() *Type {
	return simple("date")
}

func TimeType() *Type {
	return simple("time")
}

func TimestampType() *Type {
	return simple("timestamp")
}

func SetType() *Type {
	return simple("set")
}

func BitType(max uint) *Type {
	t := simple("text")
	t.Max = max
	return t
}
