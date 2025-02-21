package seeding

import (
	"github.com/kayprogrammer/ednet-fiber-api/modules/admin"
	"github.com/kayprogrammer/ednet-fiber-api/modules/courses"
	"github.com/kayprogrammer/ednet-fiber-api/modules/instructors"
)

var (
	instructorManager = instructors.InstructorManager{}
	courseManager = courses.CourseManager{}
	adminManager = admin.AdminManager{}
)
