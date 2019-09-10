package main

import (
	"fmt"
)

type user struct {
	name string 
	email string 
}

type admin struct {
	user //Embed Type 
	level string 

}

//Method for user 
func (u *user) notify(){
	fmt.Printf(" Sending User Email to %s<%s>\n",
				 u.name, 
				 u.email)
}

func main() {

	/* 
	 inner type promotion, everything that is declared within the inner type is promoted to the outer type. 
	 This means through a value of the outer type, we can access any field or method associated with the inner type value
	  directly based on the rules of exporting.
	*/
adminUser : = admin{
				user: user {
					name: "Joshua",
					email: "jaguilar20@gmail.com",
				}
				level:"Super"
			}

//long Way 			
adminUser.user.notify()

//Can Access Method Directly
adminUser.notify()

}
