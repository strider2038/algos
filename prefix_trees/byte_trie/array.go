package byte_trie

import (
	"bytes"
	"encoding/json"
)

// Array префиксное дерево с оптимизацией памяти для хранения данных с произвольными
// ключами в виде слайса байт.
//
// Оптимизация заключается в том, что дочерние узлы хранятся в массиве переменной длины
// (в исходном алгоритме выделяется массив из N элементов). Для оптимизации
// доступа индексы дочерних узлов хранятся в отдельной битовой маске.
// Таким образом класс сложности доступа к дочерним узлам остается O(1) вместо O(N)
// в варианте когда необходимо перебирать все дочерние узлы.
type Array[V any] struct {
	root  arrayNode[V]
	count int
}

func (array *Array[V]) Count() int {
	return array.count
}

func (array *Array[V]) Get(key []byte) V {
	v, _ := array.Find(key)

	return v
}

func (array *Array[V]) Find(key []byte) (V, bool) {
	node := &array.root
	var zero V

	for _, k := range key {
		// если индекс отсутствует в маске, то такого элемента нет в дереве
		if !node.bits.isSet(k) {
			return zero, false
		}
		// по номеру символа находим индекс следующего подузла дерева
		i := node.bits.getOneNumber(k)
		node = &node.children[i]
	}

	if node.value != nil {
		return *node.value, true
	}

	return zero, false
}

func (array *Array[V]) Put(key []byte, value V) {
	node := &array.root

	for _, k := range key {
		i := 0
		// если индекс найден в маске, то находим номер следующего узла в массиве
		if node.bits.isSet(k) {
			i = node.bits.getOneNumber(k)
		} else {
			// если не найден, то устанавливаем бит
			node.bits.set(k)
			// находим его порядковый номер
			i = node.bits.getOneNumber(k)
			// расширяем массив вставляя новый элемент по указанному индексу i
			node.insertChildAt(i, k)
		}
		node = &node.children[i]
	}

	// если такого элемента еще не существовало в дереве, то
	// увеличиваем счетчик количества элементов
	if node.value == nil {
		array.count++
	}

	node.value = &value
}

// Delete удаляет значение из ассоциативного массива. Реализовано в виде
// простой версии без освобождения памяти и уменьшения количества узлов.
func (array *Array[V]) Delete(key []byte) {
	node := &array.root

	for _, k := range key {
		// если индекс отсутствует в маске, то такого элемента нет в дереве
		if !node.bits.isSet(k) {
			return
		}
		// по номеру символа находим индекс следующего подузла дерева
		i := node.bits.getOneNumber(k)
		node = &node.children[i]
	}

	// если значение установлено для узла, то уменьшаем счетчик количества элементов
	if node.value != nil {
		array.count--
	}

	// удаляем ссылку на значение
	node.value = nil
}

// Walk перебирает дерево и для каждого существующего узла вызывает функцию f.
func (array *Array[V]) Walk(f func(key []byte, value V) error) error {
	return array.root.walk(nil, f)
}

func (array Array[V]) MarshalJSON() ([]byte, error) {
	var data bytes.Buffer
	data.WriteRune('{')
	i := 0

	err := array.Walk(func(key []byte, value V) error {
		if i > 0 {
			data.WriteRune(',')
		}

		k, _ := json.Marshal(string(key))
		data.Write(k)

		data.WriteRune(':')
		v, err := json.Marshal(value)
		if err != nil {
			return err
		}

		data.Write(v)

		i++

		return nil
	})
	if err != nil {
		return nil, err
	}
	data.WriteRune('}')

	return data.Bytes(), nil
}

type arrayNode[V any] struct {
	// Символ
	k byte
	// Битовая маска для индексации массива нижележащих узлов
	bits bitIndex
	// Массив нижележащих узлов переменной длины (на основе слайса)
	children []arrayNode[V]
	// Ссылка на значение ассоциативного массива
	value *V
}

func (node *arrayNode[V]) insertChildAt(index int, k byte) {
	n := arrayNode[V]{k: k}
	if len(node.children) == index {
		// вставка в конец слайса (расширение массива)
		node.children = append(node.children, n)
		return
	}

	// вставка в середину слайса со смещением элементов > index вправо
	node.children = append(node.children[:index+1], node.children[index:]...)
	node.children[index] = n
}

func (node *arrayNode[V]) walk(key []byte, f func(key []byte, value V) error) error {
	for _, child := range node.children {
		k := append(key, child.k)
		if child.value != nil {
			if err := f(k, *child.value); err != nil {
				return err
			}
		}
		if err := child.walk(k, f); err != nil {
			return err
		}
	}

	return nil
}
