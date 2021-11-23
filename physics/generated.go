package physics

type TransformList []Transform
func (t *TransformList) ComponentSet(val interface{}) { *t = *val.(*TransformList) }
func (t *TransformList) InternalRead(index int, val interface{}) { *val.(*Transform) = (*t)[index] }
func (t *TransformList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Transform) }
func (t *TransformList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Transform)) }
func (t *TransformList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *TransformList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *TransformList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *TransformList) Len() int { return len(*t) }

type RigidbodyList []Rigidbody
func (t *RigidbodyList) ComponentSet(val interface{}) { *t = *val.(*RigidbodyList) }
func (t *RigidbodyList) InternalRead(index int, val interface{}) { *val.(*Rigidbody) = (*t)[index] }
func (t *RigidbodyList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Rigidbody) }
func (t *RigidbodyList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Rigidbody)) }
func (t *RigidbodyList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *RigidbodyList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *RigidbodyList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *RigidbodyList) Len() int { return len(*t) }

type InputList []Input
func (t *InputList) ComponentSet(val interface{}) { *t = *val.(*InputList) }
func (t *InputList) InternalRead(index int, val interface{}) { *val.(*Input) = (*t)[index] }
func (t *InputList) InternalWrite(index int, val interface{}) { (*t)[index] = val.(Input) }
func (t *InputList) InternalAppend(val interface{}) { (*t) = append((*t), val.(Input)) }
func (t *InputList) InternalPointer(index int) interface{} { return &(*t)[index] }
func (t *InputList) InternalReadVal(index int) interface{} { return (*t)[index] }
func (t *InputList) Delete(index int) {
	lastVal := (*t)[len(*t)-1]
	(*t)[index] = lastVal
	(*t) = (*t)[:len(*t)-1]
}
func (t *InputList) Len() int { return len(*t) }
