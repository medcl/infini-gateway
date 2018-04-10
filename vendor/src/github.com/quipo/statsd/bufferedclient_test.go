package statsd

import (
	"fmt"
	"os"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"testing"
	"time"
)

func TestBufferedTotal(t *testing.T) {
	ln, udpAddr := newLocalListenerUDP(t)
	defer ln.Close()

	prefix := "myproject."

	client := NewStatsdClient(udpAddr.String(), prefix)
	buffered := NewStatsdBuffer(time.Millisecond*20, client)

	ch := make(chan string, 0)

	s := map[string]int64{
		"a:b:c": 5,
		"d:e:f": 2,
		"x:b:c": 5,
		"g.h.i": 1,
	}

	expected := make(map[string]int64)
	for k, v := range s {
		expected[k] = v
	}

	// also test %HOST% replacement
	s["zz.%HOST%"] = 1
	hostname, err := os.Hostname()
	expected["zz."+hostname] = 1

	go doListenUDP(t, ln, ch, len(s))

	err = buffered.CreateSocket()
	if nil != err {
		t.Fatal(err)
	}
	defer buffered.Close()

	for k, v := range s {
		buffered.Total(k, v)
	}

	actual := make(map[string]int64)

	re := regexp.MustCompile(`^(.*)\:(\d+)\|(\w).*$`)

	received := 0

	for received < len(s) {
		batch := <-ch
		for _, x := range strings.Split(batch, "\n") {
			x = strings.TrimSpace(x)
			//fmt.Println(x)
			if !strings.HasPrefix(x, prefix) {
				t.Errorf("Metric without expected prefix: expected '%s', actual '%s'", prefix, x)
				return
			}
			received++
			vv := re.FindStringSubmatch(x)
			fmt.Println(vv, x)
			if vv[3] != "t" {
				t.Errorf("Metric without expected suffix: expected 't', actual '%s'", vv[3])
			}
			v, err := strconv.ParseInt(vv[2], 10, 64)
			if err != nil {
				t.Error(err)
			}
			actual[vv[1][len(prefix):]] = v
		}
	}

	if !reflect.DeepEqual(expected, actual) {
		t.Errorf("did not receive all metrics: Expected: %T %v, Actual: %T %v ", expected, expected, actual, actual)
	}
}
