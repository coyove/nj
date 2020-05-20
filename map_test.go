package potatolang

// var N = 16
//
// func TestMapBatch(t *testing.T) {
// 	rand.Seed(time.Now().Unix())
//
// 	m := treeMap{}
// 	m2 := map[Value]bool{}
//
// 	args := make([]Value, N*2)
// 	for i := 0; i < N; i++ {
// 		x := Value{unsafe.Pointer(uintptr(rand.Int()))}
// 		m2[x] = true
// 		args[i] = x
// 		args[i+N] = x
// 	}
//
// 	m.BatchSet(args)
//
// 	for k := range m2 {
// 		if v, _ := m.Get(k); v != k {
// 			t.Fatal(m)
// 		}
// 	}
// }
//
// func TestMapAdd(t *testing.T) {
// 	m := treeMap{}
// 	m2 := map[Value]bool{}
// 	for i := 0; i < N; i++ {
// 		x := Value{unsafe.Pointer(uintptr(rand.Int()))}
// 		m.Add(true, x, x)
// 		m2[x] = true
// 	}
//
// 	for k := range m2 {
// 		if v, _ := m.Get(k); v != k {
// 			t.Fatal(m)
// 		}
// 	}
// }
