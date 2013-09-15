package app

import (
	photoJobs "github.com/paytonrules/photolibrary/website/app/jobs"
	"github.com/robfig/revel"
	"github.com/robfig/revel/modules/jobs/app/jobs"
)

func init() {
	revel.OnAppStart(func() {
		// Fix don't run on weekends
		thumbnailServerUrl, _ := revel.Config.String("thumbnail_server")
		directory, _ := revel.Config.String("root_dir")
		jobs.Schedule("@midnight", photoJobs.GenerateThumbnails{Server: thumbnailServerUrl,
			Directory: directory,
			Duration:  10800})
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
