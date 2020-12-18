package uwp

import (
	"log"
	"reflect"
	"testing"
	"time"

	"github.com/pkg/errors"
)

func TestCasePool(t *testing.T) {
	p := New(2)
	defer func() {
		log.Println("error:", p.Wait().Close())
	}()

	p.Run().
		Add(func() error {
			log.Println("Start 1 task")
			time.Sleep(time.Second)
			log.Println("Finish 1 task")
			return errors.New("first error")
		}).
		Add(func() error {
			log.Println("Start 2 task")
			time.Sleep(time.Second * 2)
			log.Println("Finish 2 task")
			return nil
		}).
		Add(func() error {
			log.Println("Start 3 task")
			time.Sleep(time.Second * 2)
			log.Println("Finish 3 task")
			return nil
		}).
		Add(func() error {
			log.Println("Start 4 task")
			time.Sleep(time.Second * 2)
			log.Println("Finish 4 task")
			return errors.New("four error")
		})

}

func TestNewPool(t *testing.T) {
	type args struct {
		concurency int
	}
	tests := []struct {
		name string
		args args
		want *Pool
	}{
		// TODO: Add test cases.
	}
	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			if got := New(tt.args.concurency); !reflect.DeepEqual(got, tt.want) {
				t.Errorf("NewPool() = %v, want %v", got, tt.want)
			}
		})
	}
}
