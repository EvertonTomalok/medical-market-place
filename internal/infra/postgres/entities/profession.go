package entities

import "fmt"

type Profession string

const (
	CNA Profession = "CNA"
	LVW Profession = "LVW"
	RN  Profession = "RN"
)

/**
* 	Returns the string representation of the enum
**/
func (p Profession) ToString() string {
	return fmt.Sprintf("%v", p)
}

/**
*	Check if the enum value is valid
**/
func (p Profession) IsValid() bool {
	for _, profession := range []string{"CNA", "LVM", "RN"} {
		if p.ToString() == profession {
			return true
		}
	}
	return false
}
