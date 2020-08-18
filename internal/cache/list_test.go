package cache //nolint:golint,stylecheck

import (
	"testing"

	"github.com/stretchr/testify/require"
)

func TestList(t *testing.T) { //nolint:go-lint
	t.Run("empty list", func(t *testing.T) {
		l := NewList()

		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	t.Run("complex", func(t *testing.T) {
		l := NewList()

		l.PushFront(10) // [10]
		l.PushBack(20)  // [10, 20]
		l.PushBack(30)  // [10, 20, 30]
		require.Equal(t, l.Len(), 3)

		middle := l.Back().Next // 20
		l.Remove(middle)        // [10, 30]
		require.Equal(t, l.Len(), 2)

		for i, v := range [...]int{40, 50, 60, 70, 80} {
			if i%2 == 0 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // [80, 60, 40, 10, 30, 50, 70]

		require.Equal(t, l.Len(), 7)
		require.Equal(t, 80, l.Front().Value)
		require.Equal(t, 70, l.Back().Value)

		l.MoveToFront(l.Front()) // [80, 60, 40, 10, 30, 50, 70]
		l.MoveToFront(l.Back())  // [70, 80, 60, 40, 10, 30, 50]

		elems := make([]int, 0, l.Len())
		for i := l.Back(); i != nil; i = i.Next {
			elems = append(elems, i.Value.(int))
		}
		require.Equal(t, []int{50, 30, 10, 40, 60, 80, 70}, elems)
	})

	// new case
	t.Run("remove all fronts", func(t *testing.T) { //nolint:go-lint
		elems := []string{"zero", "one", "two", "three", "four", "five", "six", "seven", "eight", "nine"}
		l := NewList()
		for i, v := range elems {
			if i < len(elems)/2 {
				l.PushFront(v)
			} else {
				l.PushBack(v)
			}
		} // Front()<["four", "three", "two", "one", "zero", "five", "six", "seven", "eight", "nine"]<Back()
		require.Equal(t, l.Len(), 10)
		require.Equal(t, "four", l.Front().Value)
		require.Equal(t, "nine", l.Back().Value)

		n := l.Len()
		for i := 0; i < n; i++ { // remove all by "Front()"
			l.Remove(l.Front())
		}
		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})

	// new case
	t.Run("remove all backs", func(t *testing.T) { //nolint:go-lint
		elems := []float64{0.0, 0.1, 0.2, 0.3, 0.4, 0.5, 0.6, 0.7, 0.8, 0.9}
		l := NewList()
		for i, v := range elems {
			if i < len(elems)/2 {
				l.PushBack(v)
			} else {
				l.PushFront(v)
			}
		} // Front()<[0.9, 0.8, 0.7, 0.6, 0.5, 0.0, 0.1, 0.2, 0.3, 0.4]<Back()
		require.Equal(t, l.Len(), 10)
		require.Equal(t, 0.9, l.Front().Value)
		require.Equal(t, 0.4, l.Back().Value)

		n := l.Len()
		for i := 0; i < n; i++ { // remove all by "Back()"
			l.Remove(l.Back())
		}
		require.Equal(t, l.Len(), 0)
		require.Nil(t, l.Front())
		require.Nil(t, l.Back())
	})
}
