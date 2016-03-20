package admin

import (
	"fmt"
	"path"
	"time"

	"github.com/qor/exchange"
	"github.com/qor/exchange/backends/csv"
	"github.com/qor/media_library"
	"github.com/qor/qor"
	"github.com/isairz/hive/app/models"
	"github.com/isairz/hive/db"
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

	Worker.RegisterJob(&worker.Job{
		Name:  "Import Mangas",
		Group: "Mangas Management",
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
			argument := arg.(*importMangaArgument)

			context := &qor.Context{DB: db.DB}

			var errorCount uint

			if err := MangaExchange.Import(
				csv.New(path.Join("public", argument.File.URL())),
				context,
				func(progress exchange.Progress) error {
					var cells = []worker.TableCell{
						{Value: fmt.Sprint(progress.Current)},
					}

					var hasError bool
					for _, cell := range progress.Cells {
						var tableCell = worker.TableCell{
							Value: fmt.Sprint(cell.Value),
						}

						if cell.Error != nil {
							hasError = true
							errorCount++
							tableCell.Error = cell.Error.Error()
						}

						cells = append(cells, tableCell)
					}

					if hasError {
						if errorCount == 1 {
							var headerCells = []worker.TableCell{
								{Value: "Line No."},
							}
							for _, cell := range progress.Cells {
								headerCells = append(headerCells, worker.TableCell{
									Value: cell.Header,
								})
							}
							qorJob.AddResultsRow(headerCells...)
						}

						qorJob.AddResultsRow(cells...)
					}

					qorJob.SetProgress(uint(float32(progress.Current) / float32(progress.Total) * 100))
					qorJob.AddLog(fmt.Sprintf("%d/%d Importing manga %v", progress.Current, progress.Total, progress.Value.(*models.Manga).ID))
					return nil
				},
			); err != nil {
				qorJob.AddLog(err.Error())
			}

			return nil
		},
		Resource: Admin.NewResource(&importMangaArgument{}),
	})

	Worker.RegisterJob(&worker.Job{
		Name:  "Export Mangas",
		Group: "Mangas Management",
		Handler: func(arg interface{}, qorJob worker.QorJobInterface) error {
			qorJob.AddLog("Exporting mangas...")

			context := &qor.Context{DB: db.DB}
			fileName := fmt.Sprintf("/downloads/mangas.%v.csv", time.Now().UnixNano())
			if err := MangaExchange.Export(
				csv.New(path.Join("public", fileName)),
				context,
				func(progress exchange.Progress) error {
					qorJob.AddLog(fmt.Sprintf("%v/%v Exporting manga %v", progress.Current, progress.Total, progress.Value.(*models.Manga).ID))
					return nil
				},
			); err != nil {
				qorJob.AddLog(err.Error())
			}

			qorJob.SetProgressText(fmt.Sprintf("<a href='%v'>Download exported mangas</a>", fileName))
			return nil
		},
	})
	return Worker
}
