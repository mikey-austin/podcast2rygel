package media

import (
	"github.com/godbus/dbus/v5"
	"github.com/godbus/dbus/v5/prop"
)

// org.gnome.MediaObject2 dbus interface
type MediaObject2 interface {
	Parent() dbus.ObjectPath // The container containing this object. If this is the root container it must point to itself.
	Type() string            // 'container', 'video', 'video.movie', 'audio', 'music', 'image' or 'image.photo'
	Path() dbus.ObjectPath   // D-bus path of the object
	DisplayName() string     // The readable name of this object
}

// org.gnome.UPnP.MediaContainer2 dbus interface
type MediaContainer2 interface {
	MediaObject2 // Inherits the MediaObject2 interface

	ChildCount() int     // u org.gnome.UPnP.MediaContainer2.ChildCount
	ItemCount() int      // u org.gnome.UPnP.MediaContainer2.ItemCount
	ContainerCount() int // u org.gnome.UPnP.MediaContainer2.ContainerCount
	Searchable() bool    // b org.gnome.UPnP.MediaContainer2.Searchable

	// aa{sv} org.gnome.UPnP.MediaContainer2.ListChildren (IN u offset, IN u max, IN as filter)
	ListChildren(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error)

	// aa{sv} org.gnome.UPnP.MediaContainer2.ListContainers (IN u offset, IN u max, IN as filter)
	ListContainers(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error)

	// aa{sv} org.gnome.UPnP.MediaContainer2.ListItems (IN u offset, IN u max, IN as filter)
	ListItems(offset uint, max uint, filter []string) ([]map[string]dbus.Variant, *dbus.Error)
}

// org.gnome.MediaItem2 dbus interface
type MediaItem2 interface {
	MediaObject2 // Inherits the MediaObject2 interface

	Urls() []string   // as org.gnome.UPnP.MediaItem2.URLs
	MimeType() string // s org.gnome.UPnP.MediaItem2.MIMEType
}

func GetMediaContainerMethods(obj MediaContainer2) map[string]interface{} {
	return map[string]interface{}{
		"ListChildren":   obj.ListChildren,
		"ListContainers": obj.ListContainers,
		"ListItems":      obj.ListItems}
}

func GetMediaContainerProps(obj MediaContainer2) prop.Map {
	props := GetProps(obj)
	props["org.gnome.UPnP.MediaContainer2"] = map[string]*prop.Prop{
		"ChildCount": {
			Value:    obj.ChildCount(),
			Writable: false,
			Emit:     prop.EmitFalse},
		"ContainerCount": {
			Value:    obj.ContainerCount(),
			Writable: false,
			Emit:     prop.EmitFalse},
		"ItemCount": {
			Value:    obj.ItemCount(),
			Writable: false,
			Emit:     prop.EmitFalse},
		"Searchable": {
			Value:    obj.Searchable(),
			Writable: false,
			Emit:     prop.EmitFalse}}
	return props
}

// Given a media object, return an exportable map
func GetProps(obj MediaObject2) prop.Map {
	return prop.Map{
		"org.gnome.UPnP.MediaObject2": map[string]*prop.Prop{
			"Parent": {
				Value:    obj.Parent(),
				Writable: false,
				Emit:     prop.EmitFalse},
			"Type": {
				Value:    obj.Type(),
				Writable: false,
				Emit:     prop.EmitFalse},
			"Path": {
				Value:    obj.Path(),
				Writable: false,
				Emit:     prop.EmitFalse},
			"DisplayName": {
				Value:    obj.DisplayName(),
				Writable: false,
				Emit:     prop.EmitFalse},
		},
	}
}
