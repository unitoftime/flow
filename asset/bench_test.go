package asset

// BenchmarkRegularPointer-12    	1000000000	         0.3453 ns/op
// BenchmarkAtomicPointer-12     	1000000000	         0.4445 ns/op

// type ptrHandle[T any] struct {
// 	ptr *T
// 	name string
// }
// func (h *ptrHandle[T]) Get() *T {
// 	return h.ptr
// }
// func BenchmarkRegularPointer(b *testing.B) {
// 	handle := &ptrHandle[MyAsset]{
// 		ptr: &MyAsset{1},
// 		name: "myfilenamepath.png",
// 	}
// 	val := 0
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		asset := handle.Get()
// 		val += asset.Health
// 	}
// 	b.Log(val)
// }
// func BenchmarkAtomicPointer(b *testing.B) {
// 	handle := &Handle[MyAsset]{
// 		name: "myfilenamepath.png",
// 	}
// 	handle.Set(&MyAsset{1})
// 	val := 0
// 	b.ResetTimer()

// 	for i := 0; i < b.N; i++ {
// 		asset := handle.Get()
// 		val += asset.Health
// 	}
// 	b.Log(val)
// }
