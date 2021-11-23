package particle

type LifetimeList []Lifetime
func (t *LifetimeList) ComponentSet(val interface{}) { *t = *val.(*LifetimeList) }
func (t *LifetimeList) InternalRead(index int, val interface{}) { *val.(*Lifetime) = (*t)[index] }
func (t *LifetimeList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Lifetime) }
func (t *LifetimeList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Lifetime)) }
func (t *LifetimeList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *LifetimeList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *LifetimeList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *LifetimeList) Len() int { return len(*t) }


type ColorList []Color
func (t *ColorList) ComponentSet(val interface{}) { *t = *val.(*ColorList) }
func (t *ColorList) InternalRead(index int, val interface{}) { *val.(*Color) = (*t)[index] }
func (t *ColorList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Color) }
func (t *ColorList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Color)) }
func (t *ColorList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *ColorList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *ColorList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *ColorList) Len() int { return len(*t) }


type SizeList []Size
func (t *SizeList) ComponentSet(val interface{}) { *t = *val.(*SizeList) }
func (t *SizeList) InternalRead(index int, val interface{}) { *val.(*Size) = (*t)[index] }
func (t *SizeList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Size) }
func (t *SizeList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Size)) }
func (t *SizeList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *SizeList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *SizeList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *SizeList) Len() int { return len(*t) }
