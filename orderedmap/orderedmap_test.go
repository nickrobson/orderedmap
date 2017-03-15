package orderedmap

import (
    "testing"
    "fmt"
    "strings"
)

type kv struct {
    key string
    val interface{}
}

func TestNewOrderedMap(t *testing.T) {
    m := NewOrderedMap()

    assertEqualsInt(t, uint32(0), m.size, "map started non-empty!")
    assertNil(t, m.first, "map started with non-nil first")
    assertNil(t, m.last, "map started with non-nil last")
}

func getTestData() []*kv {
    data := make([]*kv, 3)
    data[0] = &kv{"a", 1}
    data[1] = &kv{"b", "hello world"}
    data[2] = &kv{"c", make([]string, 0)}
    return data
}

func TestOrderedMap_Set(t *testing.T) {
    m := NewOrderedMap()

    data := getTestData()
    for _, kv := range data {
        m.Set(kv.key, kv.val)
    }

    assertEqualsInt(t, uint32(len(data)), m.size, "wrong size!")
}

func TestOrderedMap_Get(t *testing.T) {
    m := NewOrderedMap()

    v, b := m.Get("a")
    assertThat(t, !b, "map returned true when not-contained")
    assertThat(t, v == nil, "map returned non-nil value when not-contained")

    m.Set("a", 1)
    v, b = m.Get("a")
    assertThat(t, b, "map says it doesn't contain the item")
    assertThat(t, v == 1, "map returned wrong value on Get")

    n, _ := m.getnodeat(0)
    assertNotNil(t, n, "inserted node is nil")
    assertEquals(t, n, m.first, "first node is not new node")
    assertEquals(t, n, m.last, "last node is not new node")
}

func TestOrderedMap_LinkedList(t *testing.T) {
    m := NewOrderedMap()

    assertNil(t, m.first, "map started with non-nil first")
    assertNil(t, m.last, "map started with non-nil last")

    data := getTestData()

    for i, kv := range data {
        ui := uint32(i)
        assertEqualsInt(t, m.size, ui, "map size is incorrect!")
        m.Set(kv.key, kv.val)
        n, _ := m.getnodeat(ui)
        assertEquals(t, n, m.last, "last is not new node")
        for j := uint32(0); j <= ui; j++ {
            n, _ := m.getnodeat(j)
            next, _ := m.getnodeat(j + 1)
            if j == 0 {
                assertEquals(t, n, m.first, "0th node is not first")
                assertEquals(t, nil, n.prev, "first node has previous")
                assertEquals(t, next, n.next, "first node has wrong next")
            }
            if j == ui {
                assertEquals(t, n, m.last, "i'th node is not last")
                assertNil(t, n.next, "last node has next")
                assertNil(t, next, "next is non-nil")
            }
            if j > 0 {
                prev, _ := m.getnodeat(j - 1)
                assertNotNil(t, prev, "prev is nil")
                assertEquals(t, prev, n.prev, "node's prev is wrong")
                if j == ui {
                    assertEquals(t, prev, m.last.prev, "last node's prev is wrong")
                }
            }
        }
    }
}

func assertThat(t *testing.T, condition bool, message string) {
    if !condition {
        t.Fatal(message)
    }
}

func assertEqualsInt(t *testing.T, expected, actual uint32, message string) {
    if expected != actual {
        msg := make([]string, 2)
        msg[0] = message
        msg[1] = fmt.Sprintf("expected (%#v) got (%#v)", expected, actual)
        t.Fatal(strings.Join(msg, " : "))
    }
}

func assertEquals(t *testing.T, expected, actual *node, message string) {
    if expected != actual {
        msg := make([]string, 2)
        msg[0] = message
        msg[1] = fmt.Sprintf("expected (%#v) got (%#v)", expected, actual)
        t.Fatal(strings.Join(msg, " : "))
    }
}

func assertNotEquals(t *testing.T, expected, actual *node, message string) {
    if expected == actual {
        msg := make([]string, 2)
        msg[0] = message
        msg[1] = fmt.Sprintf("expected (%#v) got (%#v)", expected, actual)
        t.Fatal(strings.Join(msg, " : "))
    }
}

func assertNil(t *testing.T, actual *node, message string) {
    if actual != nil {
        msg := make([]string, 2)
        msg[0] = message
        msg[1] = fmt.Sprintf("expected nil got (%#v)", actual)
        t.Fatal(strings.Join(msg, " : "))
    }
}

func assertNotNil(t *testing.T, actual *node, message string) {
    if actual == nil {
        msg := make([]string, 2)
        msg[0] = message
        msg[1] = fmt.Sprintf("expected non-nil got (%#v)", actual)
        t.Fatal(strings.Join(msg, " : "))
    }
}