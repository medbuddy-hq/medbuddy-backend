package jobs

import (
	"context"
	"fmt"
	"github.com/go-co-op/gocron"
	"medbuddy-backend/internal/config"
	"medbuddy-backend/internal/constant"
	"medbuddy-backend/internal/model"
	"medbuddy-backend/pkg/repository/mongo"
	"medbuddy-backend/utility"
	"time"
)

var (
	emailSubject = "Reminder to take %s"

	CronScheduler *gocron.Scheduler
	logger        = utility.NewLogger()
)

type Cron struct {
	scheduler *gocron.Scheduler
}

func NewCronJob() *Cron {
	s := gocron.NewScheduler(time.UTC)

	CronScheduler = s
	return &Cron{s}
}

func (c *Cron) StartJobs() {
	// 4
	c.scheduler.Every(constant.TimeLapseInMinutes).Minute().Do(fetchTasks)

	// 5
	c.scheduler.StartAsync()
	logger.Info("Background cron jobs started...")
}

func StopJobs() {
	CronScheduler.Stop()
	logger.Info("SUCCESSFULLY STOPPED CRON JOBS")
}

func fetchTasks() {
	ctx := context.Background()

	dbRepo := mongo.GetDB()
	tasks, err := dbRepo.GetLatestTasks(ctx, time.Now())
	if err != nil {
		logger.Error("Could not fetch latest tasks, got error: ", err.Error())
		return
	}

	logger.Infof("Successfully fetched %v tasks to be executed \n", len(tasks))

	var timeIntervals []time.Duration
	for idx, task := range tasks {
		var interval time.Duration
		if idx == 0 {
			interval = task.Time.Sub(time.Now())
		}
		if idx > 0 {
			interval = task.Time.Sub(tasks[idx-1].Time)

		}

		timeIntervals = append(timeIntervals, interval)
	}

	runTasks(timeIntervals, tasks)
}

func runTasks(intervals []time.Duration, tasks []model.LatestTaskResponse) {
	for i := 0; i < len(tasks); i++ {
		time.Sleep(intervals[i])
		go func(task model.LatestTaskResponse) {
			config := config.GetConfig()
			emailEntity := utility.NewEmail(config.EmailDomain, fmt.Sprintf(emailSubject, task.Medication.Medicine.Name),
				task.Medication.Patient.Email, config.MailgunEmailKey)

			err := emailEntity.SendReminderEmail(logger, &task.Medication)
			if err != nil {
				logger.Errorf("Got error while sending email to '%s', error: %s", task.Medication.Patient.Email, err.Error())
				return
			}

			logger.Infof("Successfully sent reminder email to '%s'", task.Medication.Patient.Email)

			db := mongo.GetDB()
			if err := db.UpdateTask(context.Background(), task.ID, constant.TaskDone); err != nil {
				logger.Error("Error updating task to `done`, error: ", err.Error())
				return
			}
			logger.Infof("Successfully updated task to send reminder email to '%s'\n", task.Medication.Patient.Email)
		}(tasks[i])
	}
}
