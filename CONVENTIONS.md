<img src="https://craterdog.com/images/CraterDog.png" width="50%">

### Naming Conventions
This short document contains the coding conventions and idioms used in this Go
module. Since clarity is very important, each major Go concept is named using a
different part of speech:
 * [Type](#Types) names are always _noun_ phrases.
 * [Method](#Methods) names are always _verb_ phrases.
 * [Interface](#Interfaces) names are always _adjective_ phrases.

This simple guide provides immediate context to any piece of code in this
repository.

### Packages
The Go best practices suggest that we keep package names short so that they are
easy to use in our code. The problem is that the import section can be rather
cryptic. Instead, we use longer descriptive package names and assign each to a
short three character variable. This makes package references in the code terse
(and recognizable) while making the import section informative. Here is an
example:
```go
package elements

import (
	fmt "fmt"
	col "github.com/craterdog/go-collection-framework"
	mat "math"
	cmp "math/cmplx"
	str "strings"
)
```

### Types
Since all types should represent well defined "things" they are named using
_noun_ phrases.

#### Structured Types
A type that contains only _read-writeable_ attributes is structural in nature
and should be made public (Capitalized) and should not define any methods. Its
sole purpose is to represent the structure and give full access to its
attributes. An example of a structural type is the following:

```go
type Association[K any, V any] struct {
	Key K
	Value V
}
```

Note that an encapsulated type may be chosen over a pure structured type even if
read-update access to all attributes is allowed. This leads us to constrained
types.

#### Constrained Types
The constrained type approach allows constraints to be enforced on specific
attributes. For example, the actual implementation of the `Association` type is
something more like this:

```go
// ASSOCIATION INTERFACE

// This interface defines the methods supported by all association-like types.
// An association associates a key with an setable value.
type AssociationLike[K Key, V Value] interface {
	GetKey() K
	GetValue() V
	SetValue(value V)
}

// ASSOCIATION IMPLEMENTATION

// This constructor creates a new association with the specified key and value.
func Association[K Key, V Value](key K, value V) AssociationLike[K, V] {
	return &association[K, V]{key, value}
}

// This type defines the structure and methods associated with a key-value
// pair. This type is parameterized as follows:
//   - K is any type of key.
//   - V is any type of value.
//
// This structure is used by the catalog type to maintain its associations.
type association[K Key, V Value] struct {
	key   K
	value V
}

// This method returns the key for this association.
func (v *association[K, V]) GetKey() K {
	return v.key
}

// This method returns the value for this association.
func (v *association[K, V]) GetValue() V {
	return v.value
}

// This method sets the value of this association to a new value.
func (v *association[K, V]) SetValue(value V) {
	v.value = value
}
```

Notice that the interface defining the Association type has a "Like" suffix on
it turning a noun phrase into an adjective phrase. This approach is used for all
types that implement one or more interfaces.

#### Encapsulated Types
A type that fully encapsulates information should restrict all access to the
information via its methods. Since encapsulated types have methods, it is
important to isolate the implementation of the type from other code that
depends on its interface. This can be done by making the type and its attributes
private (lowercase) and its methods public (Capitalized), and then adding one or
more public constructor functions that create an instance of the private type
and return it **via an interface**. For example:

```go
// LIST INTERFACE

// This interface consolidates all the interfaces supported by list-like
// collections.
type ListLike[V Value] interface {
	Sequential[V]
	Indexed[V]
	...
}

// This constructor creates a new empty list and returns it via an interface.
func List[V Value]() ListLike[V] {
	values := make([]V, 0, 4)
	return &list[V]{values}
}

// LIST IMPLEMENTATION

// This type defines the structure and methods associated with a list of values.
type list[V Value] struct {
	values   []V
}

// SEQUENTIAL INTERFACE

// This method determines whether or not this list is empty.
func (v *list[V]) IsEmpty() bool {
	return len(v.values) == 0
}

// This method returns the number of values contained in this list.
func (v *list[V]) GetSize() int {
	return len(v.values)
}

// This method returns all the values in this list as an array (slice).
func (v *list[V]) AsArray() []V {
	length := len(v.values)
	result := make([]V, length)
	copy(result, v.values)
	return result
}

// INDEXED INTERFACE
...

```

Notice that everything in the implementation references interfaces rather than
specific types. That means that the only place where a specific type must be
referenced is when creating a new instance of that type. This approach makes it
easy to swap in new implementations without changing the interfaces or
implementation code.

For coding conventions associated with interfaces see [Interfaces](#Interfaces).

#### Extended Primitive Types
The Go primitive types (`int`, `string`, `[]byte`, etc.) can be extended to
support methods while preserving their primitive type compatibility with the
built-in operations for the Go language. Again it is important to isolate the
implementation of the methods from the code that depends on the extended type.
This can be done in a similar manner to encapsulated types as follows:

```go
// ARRAY IMPLEMENTATION

// This type defines the structure and methods associated with a native array of
// values. Each value is associated with an implicit positive integer index. The
// array uses ORDINAL based indexing rather than ZERO based indexing (see the
// description of what this means in the Sequential interface definition).
// This type is parameterized as follows:
//   - V is any type of value.
type Array[V Value] []V

// STRINGER INTERFACE

// This method implements the Go Stringer interface using the canonical
// formatter.
func (v Array[V]) String() string {
	return FormatValue(v)
}

// SEQUENTIAL INTERFACE

// This method determines whether or not this array is empty.
func (v Array[V]) IsEmpty() bool {
	return len(v) == 0
}

// This method returns the number of values contained in this array.
func (v Array[V]) GetSize() int {
	return len(v)
}

// This method returns all the values in this array. The values retrieved are in
// the same order as they are in the array.
func (v Array[V]) AsArray() []V {
	var length = len(v)
	var result = make([]V, length)
	copy(result, v)
	return result
}

// INDEXED INTERFACE
...

```

Again, for coding conventions associated with interfaces see [Interfaces](#Interfaces).

#### Function Types
The signature of a function can be used as a type. This makes it easy to pass
different implementations of a function as a parameter to another function.
Since all types are named using noun phrases function type names end with the
word "Function". Here are some examples:

```go
// COLLATION FUNCTION TYPES

// This type defines the function signature for any function that can determine
// whether or not two specified values are equal to each other.
type ComparisonFunction func(first any, second any) bool

// This type defines the function signature for any function that can determine
// the relative ordering of two specified values. The result must be one of
// the following:
//   - -1: The first value is less than the second value.
//   - 0: The first value is equal to the second value.
//   - 1: The first value is more than the second value.
type RankingFunction func(first any, second any) int

```

### Methods
Since methods act on instances of types they are named using _verb_ phrases.
This naming convention differs slightly from the convention recommended by the
Gopher community concerning "getters". We think that clarity and consistency are
more important than brevity.

To make it easier to spot the target variable in a method implementation we use
the single character variable name "v" for the target "value" of every method.
This also makes it easier to use a method from one type as a template when
creating a similar method for a different type. When copying or moving
statements from one method to another there is no need to change the target
variable name.

Therefore the declaration of each method has the same form:

```go
func (v *TargetType) MethodName(...) ... {
	...
}
```

#### Transformer Methods
A method that transforms the state of an object into a different type of object
is a transformer method. This type of method is most common when transforming a
custom type into a Go primitive type (`int`, `string`, `[]byte`, etc.). The
of the method begins with "As" followed by the target type. For example:
 * `AsInteger() int`
 * `AsReal() float64`
 * `AsImaginary() float64`
 * `AsString() string`
 * `AsArray() []T`

#### Getter Methods
A method that returns the value of one of the attributes (or characteristics) of
a structured type is a getter method. The Go naming conventions say not to
include the "Get" in front of the attribute name in the method. But our methods
are verb phrases and leaving off the "Get" would result in a noun phrase which
is reserved for types and variables. Therefore we prefix getter methods with
"Get" followed by the attribute name. For example:
 * `GetSize() int`
 * `GetPath() string`
 * `GetIndex(value T) int`
 * `GetValues(first, last int) []T`

#### Setter Methods
A method that sets the value of one of the attributes of a structured type is a
setter method. Consistent with the getter methods each setter method is prefixed
with "Set" followed by the attribute name. For example:
 * `SetPath(path string)`
 * `SetPassphrase(passphrase []rune)`
 * `SetValues(first, last int, values []T)`

Note that just because a type has a getter method for an attribute doesn't mean
that it necessarily should have a setter method for the same attribute. Some
attributes are designed to be read-only.

#### Question Methods
A method that determines whether or not something is true is a question method.
A question method generally (though not always) begins with a _to be verb_ (i.e.
"Is", "Am", "Are", "Was", "Were", "Been", "Being") followed by the condition
that is being checked. Often question methods are used as the pseudo-getter
method for a boolean attribute.  Here are some examples of question methods:

 * `IsActive() bool`
 * `IsEmpty() bool`
 * `HasFailed() bool`
 * `WasCancelled() bool`
 * `MatchesText(text string) bool`

#### Action Methods
All other methods should perform some _action_ on the target value. If it is
performing an action that involves more than one value of the target type, it
should be implemented as a library function rather than a method. All action
methods should be named with a verb phrase denoting the action being performed.
Here are some examples of action methods:

 * `AddValue(value T)`
 * `ShuffleValues()`
 * `RemoveTop() T`
 * `EmitToken(token int)`

### Interfaces
Since an interface describes only one aspect of a type or function library we
use an _adjective_ phrase to name each interface. This is different from the
approach often used within the Go packages which just adds "er" to the
end of a method name within the interface. In general (not always) an interface
with a single method doesn't make much sense. Some of the Go packages contain
interfaces that are named using _noun_ phrases. But since an interface generally
describes only one aspect of a type an adjective phrase is more appropriate.

In general the methods associated with a type, and the library functions
associated with a library, can be grouped according to their cohesiveness. A
well designed interface conveys that cohesiveness quickly and concisely.

#### Type Interfaces
One important goal of a type interface is to isolate the code that uses the type
from the implementation details of the type. This allows the developers of the
type to change its implementation without breaking the dependent code. With that
in mind, a good interface should not depend on anything but native Go types and
other interfaces. It is also fine for an interface to depend on a structured
type or a function type since neither type has a hidden implementation.

Here are a few examples of simple type interfaces:

```go
// This interface defines the methods supported by all complex types.
type Complex interface {
	GetReal() float64
	GetImaginary() float64
	GetMagnitude() float64
	GetPhase() float64
}

// This interface defines the methods supported by all temporal elements.
type Temporal interface {
	// Return the entire time in specific units.
	AsMilliseconds() float64
	AsSeconds() float64
	AsMinutes() float64
	AsHours() float64
	AsDays() float64
	AsWeeks() float64
	AsMonths() float64
	AsYears() float64

	// Return a specific part of the entire time.
	GetMilliseconds() int
	GetSeconds() int
	GetMinutes() int
	GetHours() int
	GetDays() int
	GetWeeks() int
	GetMonths() int
	GetYears() int
}

// This interface defines the methods supported by all associative collections
// whose values consist of key-value pair associations.
type Associative[K any, V any] interface {
	AddAssociation(association AssociationLike[K, V])
	AddAssociations(associations []AssociationLike[K, V])
	GetKeys() []K
	GetValue(key K) V
	GetValues(keys []K) []V
	SetValue(key K, value V) V
	RemoveValue(key K) V
	RemoveValues(keys []K) []V
	RemoveAll()
	SortAssociations()
	SortAssociationsWithRanker(ranker RankingFunction)
	ReverseAssociations()
}
```
