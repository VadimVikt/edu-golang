package hw04lrucache

import (
	"log"
)

type List interface {
	Len() int
	Front() *ListItem
	Back() *ListItem
	PushFront(v interface{}, key Key) *ListItem
	PushBack(v interface{}, key Key) *ListItem
	Remove(i *ListItem)
	MoveToFront(i *ListItem)
}

type ListItem struct {
	Value interface{}
	Key   Key
	Next  *ListItem
	Prev  *ListItem
}

type list struct {
	head *ListItem
	tail *ListItem
	size int
}

func (l *list) Len() int {
	return l.size
}

func (l *list) Front() *ListItem {
	return l.head
}

func (l *list) Back() *ListItem {
	return l.tail
}

func (l *list) PushFront(v interface{}, key Key) *ListItem {
	newListItem := &ListItem{Value: v, Key: key, Next: l.head, Prev: nil}
	if l.head != nil {
		l.head.Prev = newListItem
	}
	l.head = newListItem
	if l.tail == nil {
		l.tail = newListItem
	}
	l.size++
	return newListItem
}

func (l *list) PushBack(v interface{}, key Key) *ListItem {
	newListItem := &ListItem{Value: v, Next: nil, Key: key, Prev: l.tail}
	if l.tail != nil {
		l.tail.Next = newListItem
	}
	l.tail = newListItem
	if l.head == nil {
		l.head = newListItem
	}
	l.size++
	return newListItem
}

func (l *list) Remove(i *ListItem) {
	if i == nil || l.size == 0 {
		log.Printf("Элемент или список пуст, удаление невозможно.")
		return
	}

	// Удаление первого элемента
	if i.Prev == nil {
		l.head = i.Next // Обновляем голову списка
		if l.head != nil {
			l.head.Prev = nil
		}
		l.size--
		return
	}

	// Удаление последнего элемента
	if i.Next == nil {
		i.Prev.Next = nil
		l.tail = i.Prev
		l.size--
		return
	}

	// Удаление промежуточного элемента
	prevItem := i.Prev
	nextItem := i.Next
	prevItem.Next = nextItem
	nextItem.Prev = prevItem
	l.size--
}

func (l *list) MoveToFront(i *ListItem) {
	// 1. Проверки на граничные случаи
	if i == nil || l.head == nil || i.Prev == nil {
		// Элемент уже в начале, список пуст или элемент не существует
		return
	}

	// 2. Сохраняем ссылку на старый первый элемент ДО изменения связей
	oldFirstItem := l.head

	// 3. Вырезаем элемент из текущего места (работает и для последнего, и для среднего)
	// Связываем соседей между собой, "пропуская" элемент i
	i.Prev.Next = i.Next

	// Если элемент НЕ последний, исправляем ссылку Prev у следующего элемента
	if i.Next != nil {
		i.Next.Prev = i.Prev
	}
	// Если элемент последний, исправляем ссылку Next у предыдущего элемента
	if i.Next == nil {
		i.Prev.Next = i.Next
		l.tail = i.Prev
	}
	// 4. Перемещаем элемент в начало (Head)
	l.head = i
	i.Prev = nil
	i.Next = oldFirstItem

	// 5. Исправляем ссылку Prev у бывшего первого элемента
	oldFirstItem.Prev = i
}

func NewList() List {
	return new(list)
}
