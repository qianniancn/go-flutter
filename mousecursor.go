package flutter

import (
	"errors"
	"fmt"
	"github.com/go-flutter-desktop/go-flutter/plugin"
	"github.com/go-gl/glfw/v3.3/glfw"
)

const mousecursorChannelName = "flutter/mousecursor"

// mousecursorPlugin implements flutter.Plugin and handles method calls to the
// flutter/mousecursor channel.
type mousecursorPlugin struct {
	channel   *plugin.MethodChannel
	messenger plugin.BinaryMessenger
	window    *glfw.Window
}

var defaultMousecursorPlugin = &mousecursorPlugin{}

var _ PluginGLFW = &mousecursorPlugin{} // compile-time type check

func (p *mousecursorPlugin) InitPlugin(messenger plugin.BinaryMessenger) error {
	p.messenger = messenger

	return nil
}
func (p *mousecursorPlugin) InitPluginGLFW(window *glfw.Window) error {
	p.window = window
	p.channel = plugin.NewMethodChannel(p.messenger, mousecursorChannelName, plugin.StandardMethodCodec{})
	p.channel.HandleFuncSync("activateSystemCursor", p.handleActivateSystemCursor)
	return nil
}

func (p *mousecursorPlugin) handleActivateSystemCursor(arguments interface{}) (reply interface{}, err error) {
	args := arguments.(map[interface{}]interface{})
	switch kind := args["kind"]; {
	case kind == "none":
		p.window.SetInputMode(glfw.CursorMode, glfw.CursorHidden)
	case kind == "basic" || kind == "forbidden" || kind == "grab" || kind == "grabbing":
		// GLFW has no cursors for "forbidden", "grab" and "grabbing"
		p.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		p.window.SetCursor(glfw.CreateStandardCursor(glfw.ArrowCursor))
	case kind == "click":
		p.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		p.window.SetCursor(glfw.CreateStandardCursor(glfw.HandCursor))
	case kind == "text":
		p.window.SetInputMode(glfw.CursorMode, glfw.CursorNormal)
		p.window.SetCursor(glfw.CreateStandardCursor(glfw.IBeamCursor))
	default:
		return nil, errors.New(fmt.Sprintf("cursor kind %s not implemented", args["kind"]))
	}
	return nil, nil
}
