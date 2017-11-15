package main

import (
	"log"
	"skeletor/utils"
)

func saveUserProfile(profile *Profile) {
	profile.Password = utils.HashPassword(profile.Password)
	err := session.QueryRow(`INSERT INTO profile ( 
		firstname, 
		lastname, 
		username, 
		title, 
		password, 
		mobilenumber
	) VALUES($1, $2, $3, $4, $5, $6) RETURNING id`,
		profile.Firstname,
		profile.Lastname,
		profile.Username,
		profile.Title,
		profile.Password,
		profile.MobileNumber).Scan(&profile.Id)
	if err != nil {
		log.Print(err)
	}
	profile.Password = ""
}