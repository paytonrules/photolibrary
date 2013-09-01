package app

import (
  "github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
	photoJobs "github.com/paytonrules/photolibrary/app/jobs"
)

func init() {
  revel.OnAppStart(func() {
    // Fix don't run on weekends
    jobs.Schedule("0 0 9 * * 1-5", photoJobs.GenerateThumbnails{Duration: 28800})
    jobs.Schedule("0 0 0 * * 0,6", photoJobs.GenerateThumbnails{Duration: 21600})
  })

	// Filters is the default set of global filters.
	revel.Filters = []revel.Filter{
		revel.PanicFilter,             // Recover from panics and display an error page instead.
		revel.RouterFilter,            // Use the routing table to select the right Action
		revel.FilterConfiguringFilter, // A hook for adding or removing per-Action filters.
		revel.ParamsFilter,            // Parse parameters into Controller.Params.
		revel.SessionFilter,           // Restore and write the session cookie.
		revel.FlashFilter,             // Restore and write the flash cookie.
		revel.ValidationFilter,        // Restore kept validation errors and save new ones from cookie.
		revel.I18nFilter,              // Resolve the requested language
		revel.InterceptorFilter,       // Run interceptors around the action.
		revel.ActionInvoker,           // Invoke the action.
	}
}
