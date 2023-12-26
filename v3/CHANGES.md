<img src="https://craterdog.com/images/CraterDog.png" width="50%">

## Go Collection Framework v3 Changes
The `v2` version of the collection framework showed weaknesses in the following
areas:
 1. The exported package level functions were generally associated with a single
    class and caused potential namespace collisions.
 1. The `Array[V]` and `Map[K, V]` types extended the `[]V` and `map[K]V` Go
    data types but were public with no constructors to constrain the possible
	values created.
 1. The `Token` type exposed its attributes in a way that could not be
    protected.
 1. The private constants were defined at the package level—the only place
    possible in Go—even though each constant was associated with a specific
	class.
 1. Since Go does not allow any non-primitive data type constants, variables
    were used to represent these constants without the ability to make them
	read-only.

The functionality of each of the classes in the collection framework worked well
and was not changed in any significant ways.  That said, the following changes
to the package made in `v3`:

### Consolidated Package Interface
The public interface to a package is now captured in a `package.go` file.  This
file defines all abstract interfaces and types that are exported by the classes
defined in the package.  If a new abstract interface or type is added to this
file, the _minor_ version number of the package must be incremented.  If an
existing abstract interface or type is changed or deleted, the _major_ version
number of the package must be incremented.

In doing this, the following class types are no longer part of the public
interface for the package:
 * Token
 * Scanner
 * Parser
 * Collator
 * Sorter[Value]

These classes are still available from the package but any dependencies on them
are not guaranteed to be stable.

### Added CDCN Notation Class
Various notations are possible for importing and exporting collections including
XML, JSON and CDCN (Crater Dog Collection Notation™).  The package currently
only provides a class that implements the latter.  But the exported abstract
interfaces and types would allow classes to be implemented that support the
other two formats mentioned above.

### Class Level Namespaces for Each Class
Each exported class defined in the collection framework now has an associated
class _namespace_.  All _constants_, _constructors_ and _class functions_ are now
accessed via its class namespace functions.  See the latest
[Crater Dog Technologies™ Go Coding Conventions](https://github.com/craterdog/go-coding-conventions/wiki#class-namespaces)
for more details on class namespaces.

### Fully Encapsulated `Array[V]` and `Map[K, V]` Classes 
These two classes which were type extensions before are now fully encapsulated
classes with their own namespaces just like the rest of the collection classes.

### `Stringent` Interface
This interface was added to codify the methods that must be supported by all
lexical parsing agents.

### `Standardized` Interface
This interface was added to codify the methods that must be supported by all
canonical formatting agents.

### `Malleable[V]` Interface Refactored and Renamed
The `Malleable[V]` interface has been renamed to `Extendable[V]` and the `AddValue()`
and `AddValues()` methods renamed to `AppendValue()` and `AppendValues()` for
clarity since they only append values to the end of a sequence.

<H5 align="center"> Copyright © 2009 - 2024  Crater Dog Technologies™. All rights reserved. </H5>
