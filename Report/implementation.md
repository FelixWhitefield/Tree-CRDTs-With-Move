# Problems with Implementation

As Go's interface implementation is not explicit, there seems to be some limitations on what interfaces can define.
For example: I have an interface `Timestamp` that defines `Clone()`. I cannot create a struct that has a `Clone()` method and returns an instance of itself, and have it also implement the interface.
I would have to have the Clone() method either return an 'any' or a blank interface, which is not ideal.
I have been able to use Go's new 1.18 Generics, however this isn't ideal as the return type is not explicit.

While this solution works, any struct which has a method called `Clone()` and returns something will implement the interface, which is not ideal.
(Ideally, the Clone() method should be defined to only return an instance of the struct that implements the interface, but this is not possible in Go.)

