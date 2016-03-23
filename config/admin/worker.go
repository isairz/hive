package admin

import (
	"fmt"
	"time"

	"github.com/qor/media_library"
	"github.com/qor/worker"
)

func getWorker() *worker.Worker {
	Worker := worker.New()

	type sendNewsletterArgument struct {
		Subject      string
		Content      string `sql:"size:65532"`
		SendPassword string
		worker.Schedule
	}

	Worker.RegisterJob(&worker.Job{
		Name: "Send Newsletter",
		Handler: func(argument interface{}, qorJob worker.QorJobInterface) error {
			qorJob.AddLog("Started sending newsletters...")
			qorJob.AddLog(fmt.Sprintf("Argument: %+v", argument.(*sendNewsletterArgument)))
			for i := 1; i <= 100; i++ {
				time.Sleep(100 * time.Millisecond)
				qorJob.AddLog(fmt.Sprintf("Sending newsletter %v...", i))
				qorJob.SetProgress(uint(i))
			}
			qorJob.AddLog("Finished send newsletters")
			return nil
		},
		Resource: Admin.NewResource(&sendNewsletterArgument{}),
	})

	type importMangaArgument struct {
		File media_library.FileSystem
	}

	return Worker
}
