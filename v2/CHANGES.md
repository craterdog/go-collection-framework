<img src="https://craterdog.com/images/CraterDog.png" width="50%">

.    Copyright (c) 2009-2024 Crater Dog Technologies.  All Rights Reserved.    .
The following changes to the module interface where introduced in v2:

### Decoupling of the List and Array Types
There was a subtle problem with the implementation of the Array[V] type in
v1. Since it depends on the default `CompareValues()` function and is just a
tilde type for `[]V` it cannot depend on any other compare function. The
`Searchable` interface in turn depended on the `Array.GetIndex()` method which
is the method that depends on the default comparer function. The `List[V]` type
delegated several of its interfaces including `Searchable` to the `Array[V]`
implementation. This coupling was too tight and broke the `List[V]` type in
subtle ways that were hard to diagnose and work-around.

To address this we moved the `GetIndex()` method from the `Accessible` interface
into the `Searchable` interface and removed support for the `Searchable` interface
from the `Array` type. The `List` type was rewritten to only delegate the
`Sequential`, the `Accessible` (renamed from the old `Indexed` interface), and the
`Updatable` interface implementations to the `Array` type.

### Removal of the `AddAssociation()` and `AddAssociations()` Methods
When attempting to define a new type that implements the `Associative[K, V]`
interface we realized that the `AddAssociation()` and `AddAssociations()` method
signatures depend on the `Binding[K, V]` interface. The Go language does not yet
fully support nested generics which causes `go vet` to mark certain type
conversions as impossible when in reality the two types have identical method
signatures.

To address this the `AddAssociation()` and `AddAssociations()` methods were
removed from the `Associative[K, V]` interface and `Catalog[K, V]` and
`Map[K, V]` types.

### Simplification of the `FIFO` and `LIFO` Interfaces
The `AddValues()` (plural) method is not consistent with the typical meaning of
a FIFO or LIFO interface. It is rare that multiple values are added to a queue
or stack at one time.

Therefore, the `AddValues()` method was removed from both the `FIFO` and `LIFO`
interface definitions. The queue and stack types were updated to not implement
that method.

<H5 align="center"> Copyright © 2009 - 2024  Crater Dog Technologies™. All rights reserved. </H5>
