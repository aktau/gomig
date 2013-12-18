package common

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

func IntType() *Type {
	return simple("integer")
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
