package smetrics

type Value struct {
	valueNext     *Value
	valuePrevious *Value

	Payload float64
}

type LinkedList struct {
	Head *Value
	Tail *Value

	NumberValues uint16
}

type Values struct {
	values        map[IdentifierEmitter]*LinkedList
	maxListLength uint16
}

func NewValues(maxListLength uint16) *Values {
	return &Values{
		maxListLength: maxListLength,
		values:        map[IdentifierEmitter]*LinkedList{},
	}
}

func (v *Values) AddValue(onList IdentifierEmitter, payload float64) {
	if linkedList, exists := v.values[onList]; exists {
		newHead := &Value{
			valueNext: linkedList.Head,
			Payload:   payload,
		}

		if linkedList.Head != nil {
			linkedList.Head.valuePrevious = newHead
		}

		linkedList.Head = newHead

		if linkedList.Tail == nil {
			linkedList.Tail = newHead
		}

		linkedList.NumberValues++

		if linkedList.NumberValues > v.maxListLength {
			newTail := linkedList.Tail.valuePrevious

			if newTail != nil {
				newTail.valueNext = nil
				linkedList.Tail = newTail
			}

			linkedList.NumberValues = v.maxListLength
		}

		return
	}

	newValue := &Value{Payload: payload}

	v.values[onList] = &LinkedList{
		Head: newValue,
		Tail: newValue,

		NumberValues: 1,
	}
}

func (v *Values) GetMetric(forList IdentifierEmitter) float64 {
	ix := v.maxListLength

	var sum float64

	head := v.values[forList].Head

	for head.valueNext != nil && ix > 0 {
		sum = sum + head.Payload
		head = head.valueNext

		ix--
	}

	if ix == v.maxListLength {
		return 0
	}

	return sum / float64(v.maxListLength-ix)
}

func (v *Values) GetNumberValues(identifier IdentifierEmitter) uint16 {
	if linkedList, exists := v.values[identifier]; exists {
		return linkedList.NumberValues
	}

	return 0
}
