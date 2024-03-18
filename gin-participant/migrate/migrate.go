package main

import (
	"gin-crud/initializers"
	model "gin-crud/models"
)

func init() {
	initializers.LoadEnvVariables()
	initializers.DatabaseInit()
}

func main() {
	c1 := initializers.DB.Migrator().HasTable(&model.ParticipantData{})
	c2 := initializers.DB.Migrator().HasTable("participant_data")

	if c1 == true && c2 == true {
		initializers.DB.Migrator().DropTable(&model.ParticipantData{})
		initializers.DB.Migrator().DropTable("participant_data")
	}

	c3 := initializers.DB.Migrator().HasTable(&model.SystemData{})
	c4 := initializers.DB.Migrator().HasTable("system_data")

	if c3 == true && c4 == true {
		initializers.DB.Migrator().DropTable(&model.SystemData{})
		initializers.DB.Migrator().DropTable("system_data")
	}

	c5 := initializers.DB.Migrator().HasTable(&model.MentorData{})
	c6 := initializers.DB.Migrator().HasTable("mentor_data")

	if c5 == true && c6 == true {
		initializers.DB.Migrator().DropTable(&model.MentorData{})
		initializers.DB.Migrator().DropTable("mentor_data")
	}

	initializers.DB.AutoMigrate(&model.SystemData{})
	initializers.DB.AutoMigrate(&model.ParticipantData{})
	initializers.DB.AutoMigrate(&model.MentorData{})

}
