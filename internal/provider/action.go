package provider

import (
	"context"
	"fmt"
	"time"

	"github.com/chromedp/cdproto/cdp"
	"github.com/chromedp/cdproto/network"
	"github.com/chromedp/chromedp"
	"github.com/chromedp/chromedp/kb"
	"github.com/hashicorp/terraform-plugin-framework/types"
)

type Action struct {
	dpaction  chromedp.Action
	valueName string
	value     *string
}

func NewAction(action chromedp.Action, valueName string, value *string) *Action {
	return &Action{
		dpaction:  action,
		valueName: valueName,
		value:     value,
	}
}

func (a *Action) Action(values map[string]*string) chromedp.Action {
	if a.valueName != "" && a.value != nil {
		values[a.valueName] = a.value
	}
	return a.dpaction
}

func actionBuilder(actionArgs []types.String) (*Action, error) {
	if len(actionArgs) < 1 {
		return nil, fmt.Errorf("malformed action")
	}
	verb := actionArgs[0]
	args := actionArgs[1:]

	var dpAction chromedp.Action
	var valueName string
	var outputValue *string

	switch verb.ValueString() {
	case "navigate":
		if len(args) != 1 {
			return nil, fmt.Errorf("navigate action expects only 1 argument (URL), got %d: %v", len(args), args)
		}
		url := args[0].ValueString()
		dpAction = chromedp.Navigate(url)
	case "wait_visible":
		if len(args) != 1 {
			return nil, fmt.Errorf("wait_visible action expects only 1 argument (selector), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		dpAction = chromedp.WaitVisible(selector)
	case "click":
		if len(args) < 1 {
			return nil, fmt.Errorf("click action expects at 1 least argument (selector and options), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		var opts []func(*chromedp.Selector)
		for _, optName := range args[1:] {
			switch optName.ValueString() {
			case "visible":
				opts = append(opts, chromedp.NodeVisible)
			}
		}
		dpAction = chromedp.Click(selector, opts...)
	case "value":
		if len(args) != 2 {
			return nil, fmt.Errorf("value action expects 2 arguments (selector and value name), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		valueName = args[1].ValueString()
		outputValue = new(string)
		dpAction = chromedp.Value(selector, outputValue)
	case "focus":
		if len(args) != 1 {
			return nil, fmt.Errorf("focus action expects only 1 argument (selector), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		dpAction = chromedp.Focus(selector)
	case "sleep":
		if len(args) != 1 {
			return nil, fmt.Errorf("sleep action expects only 1 argument (duration), got %d: %v", len(args), args)
		}
		d, err := time.ParseDuration(args[0].ValueString())
		if err != nil {
			return nil, fmt.Errorf("can't parse duration for sleep: %w", err)
		}
		dpAction = chromedp.Sleep(d)
	case "text":
		if len(args) != 2 {
			return nil, fmt.Errorf("text action expects 2 arguments (selector and value name), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		valueName = args[1].ValueString()
		outputValue = new(string)
		dpAction = chromedp.TextContent(selector, outputValue)
	case "cookie":
		if len(args) < 2 {
			return nil, fmt.Errorf("cookie action expects at least 2 arguments (cookie name, value, and optional domain), got %d: %v", len(args), args)
		}
		cookieName := args[0].ValueString()
		cookieValue := args[1].ValueString()
		cookieDomain := ""

		if len(args) == 2 {
			cookieDomain = args[2].ValueString()
		}
		dpAction = chromedp.ActionFunc(func(ctx context.Context) error {
			expr := cdp.TimeSinceEpoch(time.Now().Add(24 * time.Hour))
			setcookieBuilder := network.SetCookie(cookieName, cookieValue).
				WithExpires(&expr)
			if cookieDomain != "" {
				setcookieBuilder = setcookieBuilder.WithDomain(cookieDomain)
			}
			err := setcookieBuilder.Do(ctx)
			if err != nil {
				return err
			}
			return nil
		})
	case "set_value":
		if len(args) != 2 {
			return nil, fmt.Errorf("set_value action expects 2 arguments (selector and value), got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		value := args[1].ValueString()
		dpAction = chromedp.SetValue(selector, value)
	case "press_enter":
		if len(args) != 0 {
			return nil, fmt.Errorf("press_enter action expects 0 arguments, got %d: %v", len(args), args)
		}
		selector := args[0].ValueString()
		dpAction = chromedp.SendKeys(selector, kb.Enter)
	default:
		return nil, fmt.Errorf("unknown action: %s", verb)
	}
	return NewAction(dpAction, valueName, outputValue), nil
}
