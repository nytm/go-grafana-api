package gapi

import (
	"testing"

	. "github.com/smartystreets/goconvey/convey"
)

func TestDashboardTags(t *testing.T) {
	Convey("given a dashboard", t, func() {
		d := NewDashboard()

		Convey("should be able to add some tags", func() {
			d.AddTags("home", "ok")
			So(d.Tags(), ShouldContain, "home")
			So(d.Tags(), ShouldContain, "ok")
		})

		Convey("should be able to set all tags", func() {
			d.AddTags("home", "ok")
			d.SetTags("set", "this")
			So(d.Tags(), ShouldNotContain, "home")
			So(d.Tags(), ShouldNotContain, "ok")
			So(d.Tags(), ShouldContain, "this")
			So(d.Tags(), ShouldContain, "set")
		})

		Convey("should be able to remove a tags", func() {
			d.AddTags("home", "ok", "horse")
			d.RemoveTags("home", "ok")
			So(d.Tags(), ShouldNotContain, "home")
			So(d.Tags(), ShouldNotContain, "ok")
			So(d.Tags(), ShouldContain, "horse")
		})
	})
}

func TestDashboardTitle(t *testing.T) {
	Convey("given a dashboard", t, func() {
		d := NewDashboard()
		d.Model["title"] = "hihi"

		Convey("it should give the title", func() {
			t, ok := d.Title()
			So(ok, ShouldBeTrue)
			So(t, ShouldEqual, "hihi")
		})
	})
}
