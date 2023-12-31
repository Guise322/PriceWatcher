package message

import (
	"testing"

	"github.com/sirupsen/logrus"
)

var val = MessageHourVal{}

func TestValidate(t *testing.T) {
	logrus.Info("Start to test the func Validate with true values")
	withTrueValues(t)
	logrus.Info("Start to test the func Validate with false values")
	withFalseValues(t)
}

func withTrueValues(t *testing.T) {
	tHour := 12
	want := true
	if got := val.Validate(tHour, []int{12, 17}); got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}

	tHour = 17
	want = true
	if got := val.Validate(tHour, []int{12, 17}); got != want {
		t.Errorf("Got %v, wanted %v", got, want)
	}
}

func withFalseValues(t *testing.T) {
	tHours := [4]int{11, 13, 16, 18}
	want := false
	for i := 0; i < len(tHours); i++ {
		if got := val.Validate(tHours[i], []int{12, 17}); got != want {
			t.Errorf("Got %v, wanted %v", got, want)
		}
	}
}
