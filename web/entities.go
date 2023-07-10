package web

type ScheduledTask struct {
	Name    string `json:"name"`
	Cron    string `json:"cron"`
	Retries int    `json:"retries"`
	TaskID  int    `json:"task_id"`
	Enabled bool   `json:"enabled"`
}
