package byte_suffix_trie

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

	for i, k := range key {
		// если индекс отсутствует в маске, то сверяем суффикс
		if !node.bits.isSet(k) {
			// если суффикс идентичен, то узел найден
			if bytes.Equal(key[i:], node.suffix) {
				break
			}

			// такого элемента нет в дереве
			return zero, false
		}
		// по значению байта находим индекс следующего подузла дерева
		node = node.child(k)
	}

	if node.present {
		return node.value, true
	}

	return zero, false
}

func (array *Array[V]) Put(key []byte, value V) {
	node := &array.root

	for i := 0; i < len(key); i++ {
		k := key[i]

		// если индекс найден в маске, то находим номер следующего узла в массиве
		if node.bits.isSet(k) {
			node = node.child(k)

			if i == len(key)-1 && len(node.suffix) > 0 {
				node.forkSuffix()
			}
			continue
		}
		// если суффикс идентичен, то узел найден
		if bytes.Equal(key[i:], node.suffix) {
			break
		}

		if len(node.suffix) > 0 {
			var shift int
			node, shift = node.splitBranch(key, i)

			// пропускаем shift элементов
			i += shift
			// если префикс полностью совпадает, то узел уже найден
			if i >= len(key) {
				break
			}

			k = key[i]
		}

		node = node.insert(k, key[i+1:])

		break
	}

	// если такого элемента еще не существовало в дереве, то
	// увеличиваем счетчик количества элементов
	if !node.present {
		array.count++
	}

	node.present = true
	node.value = value
}

// Delete удаляет значение из ассоциативного массива. Реализовано в виде
// простой версии без освобождения памяти и уменьшения количества узлов.
func (array *Array[V]) Delete(key []byte) {
	node := &array.root

	for i, k := range key {
		// если индекс отсутствует в маске, то такого элемента нет в дереве
		if !node.bits.isSet(k) {
			// если суффикс идентичен, то узел найден
			if bytes.Equal(key[i:], node.suffix) {
				break
			}

			// такого элемента нет в дереве
			return
		}

		// по номеру символа находим индекс следующего подузла дерева
		node = node.child(k)
	}

	// если значение установлено для узла, то уменьшаем счетчик количества элементов
	if node.present {
		array.count--
	}

	// сбрасываем флаг присутствия
	node.present = false
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
	// Флаг наличия значения
	present bool
	// Символ
	k byte
	// Суффикс ключа
	suffix []byte
	// Битовая маска для индексации массива нижележащих узлов
	bits bitIndex
	// Массив нижележащих узлов переменной длины (на основе слайса)
	children []arrayNode[V]
	// Ссылка на значение ассоциативного массива
	value V
}

// child - по значению байта находим индекс следующего подузла дерева
func (node *arrayNode[V]) child(k byte) *arrayNode[V] {
	childIndex := node.bits.getOneNumber(k)

	return &node.children[childIndex]
}

func (node *arrayNode[V]) insert(k byte, suffix []byte) *arrayNode[V] {
	// если не найден, то устанавливаем бит
	node.bits.set(k)
	// находим его порядковый номер
	childIndex := node.bits.getOneNumber(k)
	// расширяем массив вставляя новый элемент по указанному индексу childIndex
	node.insertChildAt(childIndex, k, suffix)

	return &node.children[childIndex]
}

func (node *arrayNode[V]) insertChildAt(index int, k byte, suffix []byte) {
	n := arrayNode[V]{k: k, suffix: suffix}
	if len(node.children) == index {
		// вставка в конец слайса (расширение массива)
		node.children = append(node.children, n)
		return
	}

	// вставка в середину слайса со смещением элементов > index вправо
	node.children = append(node.children[:index+1], node.children[index:]...)
	node.children[index] = n
}

// splitBranch - разделяет текущую цепочку на основе суффикса и части ключа
func (node *arrayNode[V]) splitBranch(key []byte, offset int) (*arrayNode[V], int) {
	currentNode := node

	nodeValue := currentNode.value // переносимое значение
	suffix := currentNode.suffix

	// в разделяемом узле удаляем значение и суффикс
	currentNode.present = false
	currentNode.suffix = nil

	// создаем пересекающуюся по префиксу ветвь дерева
	j := 0
	for ; j < len(suffix) && offset+j < len(key) && suffix[j] == key[offset+j]; j++ {
		currentNode = currentNode.insert(key[offset+j], nil)
	}

	if j >= len(suffix) {
		// переносим значение в конец ветви (если был последний узел)
		currentNode.present = true
		currentNode.value = nodeValue
	} else {
		// вставляем конечный узел с оставшейся частью суффикса в старую ветвь
		n := currentNode.insert(suffix[j], suffix[j+1:])
		n.k = suffix[j]
		n.present = true
		n.value = nodeValue
	}

	return currentNode, j
}

// forkSuffix - отделяет суффикс в отдельную ветку по первому символу
func (node *arrayNode[V]) forkSuffix() {
	// переносимое значение
	nodeValue := node.value
	suffix := node.suffix

	// в разделяемом узле удаляем значение и суффикс
	node.present = false
	node.suffix = nil

	n := node.insert(suffix[0], suffix[1:])
	n.k = suffix[0]
	n.present = true
	n.value = nodeValue
}

func (node *arrayNode[V]) walk(key []byte, f func(key []byte, value V) error) error {
	for _, child := range node.children {
		k := append(key, child.k)
		if child.present {
			if err := f(append(k, child.suffix...), child.value); err != nil {
				return err
			}
		}
		if err := child.walk(k, f); err != nil {
			return err
		}
	}

	return nil
}
