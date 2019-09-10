package main

// All material is licensed under the Apache License Version 2.0, January 2004
// https://play.golang.org/p/uB4c33sbfj
// Sample program demonstrating decoupling with interface composition.

import "fmt"

// =============================================================================

// Board represents a surface we can work on.
type Board struct {
	NailsNeeded int
	NailsDriven int
}

// =============================================================================

// NailDriver represents behavior to drive nails into a board.
type NailDriver interface {
	DriveNail(nailSupply *int, b *Board)
}

// NailPuller represents behavior to remove nails into a board.
type NailPuller interface {
	PullNail(nailSupply *int, b *Board)
}

/*
This interface is composed from both the NailDriver and NailPuller interfaces. This is a very common pattern,
taking existing interfaces and grouping them into composed behaviors. You will see how this plays into the code later on.
For now, any concrete type value that implements both the driver and puller behaviors will also implement the NailDrivePuller interface.
*/
type NailDrivePuller interface {
	NailDriver
	NailPuller
}

// =============================================================================

// Mallet is a tool that pounds in nails.
type Mallet struct{}

// DriveNail pounds a nail into the specified board.
func (Mallet) DriveNail(nailSupply *int, b *Board) {

	// Take a nail out of the supply.
	*nailSupply--

	// Pound a nail into the board.
	b.NailsDriven++

	fmt.Println("Mallet: pounded nail into the board.")
}

// Crowbar is a tool that removes nails.
type Crowbar struct{}

// PullNail yanks a nail out of the specified board.
func (Crowbar) PullNail(nailSupply *int, b *Board) {

	// Yank a nail out of the board.
	b.NailsDriven--

	// Put that nail back into the supply.
	*nailSupply++

	fmt.Println("Crowbar: yanked nail out of the board.")
}

// =============================================================================

// Toolbox can contains any type of Driver and Puller.
type Toolbox struct {
	NailDriver
	NailPuller

	nails int
}

/*
	We have not embedded a struct type into our Toolbox but two interface types. T
	his means any concrete type value that implements the NailDriver
	interface can be assigned as the inner type value for the NailDriver embedded interface type. The same holds true for the embedded NailPuller interface type.
	Once a concrete type is assigned, the Toolbox is then guaranteed to implement this behavior.
	Even more, since the toolbox embeds both a NailDriver and NailPuller interface type,
	this means a Toolbox also implements the NailDrivePuller interface as well
*/
// =============================================================================

// Contractor carries out the task of securing boards.
type Contractor struct{}

// Fasten will drive nails into a board.
func (Contractor) Fasten(d NailDriver, nailSupply *int, b *Board) {
	for b.NailsDriven < b.NailsNeeded {
		d.DriveNail(nailSupply, b)
	}
}

/*
 The method Fasten is declared to provide a contractor the behavior to drive the number of nails that are needed into a specified board.
 The method requires the user to pass as the first parameter a value that implements the NailDriver interface.
 This value represents the tool the contractor will use to execute this behavior.
 Using an interface type for the this parameter allows the user of the API to later create and use different tools without the need for the API to change.
 The user is providing the behavior of the tooling and the Fasten method is providing the workflow for when and how the tool is used.
*/

func (Contractor) Unfasten(p NailPuller, nailSupply *int, b *Board) {
	for b.NailsDriven > b.NailsNeeded {
		p.PullNail(nailSupply, b)
	}
}

/*
Notice the Fasten method requires a value of interface type NailDriver and we are passing a value of interface type NailDrivePuller.
This is possible because the compiler knows that any concrete type value that can be stored inside a NailDrivePuller interface value
 must also implement the NailDriver interface.
Therefore, the compiler accepts the method call and the assignment between these two interface type values

*/

func (c Contractor) ProcessBoards(dp NailDrivePuller, nailSupply *int, boards []Board) {
	for i := range boards {
		b := &boards[i]

		fmt.Printf("Contractor: examining board #%d: %+v\n", i+1, b)

		switch {
		case b.NailsDriven < b.NailsNeeded:
			c.Fasten(dp, nailSupply, b)

		case b.NailsDriven > b.NailsNeeded:
			c.Unfasten(dp, nailSupply, b)
		}
	}
}

// =============================================================================

// main is the entry point for the application.
func main() {

	// Inventory of old boards to remove, and the new boards
	// that will replace them.
	boards := []Board{

		// Rotted boards to be removed.
		{NailsDriven: 3},
		{NailsDriven: 1},
		{NailsDriven: 6},

		// Fresh boards to be fastened.
		{NailsNeeded: 6},
		{NailsNeeded: 9},
		{NailsNeeded: 4},
	}

	// Fill a toolbox.
	tb := Toolbox{
		NailDriver: Mallet{},
		NailPuller: Crowbar{},
		nails:      10,
	}

	// Hire a Contractor and put our Contractor to work.
	var c Contractor
	c.ProcessBoards(&tb, &tb.nails, boards)
}
