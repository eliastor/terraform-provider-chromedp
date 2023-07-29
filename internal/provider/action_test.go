package provider

import (
	"testing"

	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/assert"
)

func TestNewActionValue(t *testing.T) {
	value := new(string)
	act := NewAction(chromedp.Value("body", value), "value", value)
	assert.NotNil(t, act)

	m := map[string]*string{}

	_ = act.Action(m)

	assert.Same(t, value, m["value"])
}

func TestAction_Action(t *testing.T) {

}
