package common

type TypeModifier uint

const (
	TypeSmall TypeModifier = iota
	TypeNormal
	TypeLarge
	TypeHuge
)

const (
	TypeFloat     = "float"
	TypeDouble    = "double"
	TypeNumeric   = "numeric"
	TypeInteger   = "integer"
	TypeBlob      = "blob"
	TypeBool      = "bool"
	TypeChar      = "char"
	TypeBit       = "bit"
	TypeText      = "text"
	TypeDate      = "date"
	TypeTime      = "time"
	TypeTimeStamp = "timestamp"
	TypeSet       = "set"
)

type Type struct {
	/* a base type, possible values: see the Type* consts */
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

func FloatType() *Type                    { return simple(TypeFloat) }
func DoubleType() *Type                   { return simple(TypeDouble) }
func IntType(modifier TypeModifier) *Type { return simplem(TypeInteger, modifier) }
func BlobType() *Type                     { return simple(TypeBlob) }
func BoolType() *Type                     { return simple(TypeBool) }
func PaddedTextType() *Type               { return simple(TypeChar) }
func TextType() *Type                     { return simple(TypeText) }
func DateType() *Type                     { return simple(TypeDate) }
func TimeType() *Type                     { return simple(TypeTime) }
func TimestampType() *Type                { return simple(TypeTimeStamp) }
func SetType() *Type                      { return simple(TypeSet) }

/* for external usage */
func SimpleType(name string) *Type {
	return simple(name)
}

func NumericType(precision, scale uint) *Type {
	return &Type{Name: TypeNumeric, Precision: precision, Scale: scale}
}

func BitType(max uint) *Type {
	t := simple(TypeBit)
	t.Max = max
	return t
}

/* for internal usage (shorter) */
func simple(name string) *Type {
	return &Type{Name: name}
}

func simplem(name string, modifier TypeModifier) *Type {
	return &Type{Name: name, Modifier: modifier}
}
