package main

type config struct {
	Global backupSetting
	Jobs   []backupJob
}

type backupJob struct {
	Name string
	Dir  string
	backupSetting
}

type backupSetting struct {
	WebDAV     string `yaml:"webdav"`
	TimeFormat string `yaml:"timeFmt"`
	User       string
	Password   string
	Rotate     string
}

func (c *config) ListJobs() []backupJob {
	result := make([]backupJob, len(c.Jobs))
	for i, item := range c.Jobs {
		result[i].Name = item.Name
		result[i].Dir = item.Dir

		result[i].TimeFormat = item.TimeFormat
		if result[i].TimeFormat == "" {
			result[i].TimeFormat = c.Global.TimeFormat
		}

		result[i].WebDAV = item.WebDAV
		if result[i].WebDAV == "" {
			result[i].WebDAV = c.Global.WebDAV
		}

		result[i].User = item.User
		if result[i].User == "" {
			result[i].User = c.Global.User
		}

		result[i].Password = item.Password
		if result[i].Password == "" {
			result[i].Password = c.Global.Password
		}

		result[i].Rotate = item.Rotate
		if result[i].Rotate == "" {
			result[i].Rotate = c.Global.Rotate
		}
	}

	return result
}
