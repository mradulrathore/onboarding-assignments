package user

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/mradulrathore/user-management/service/course/enum"
	"github.com/mradulrathore/user-management/validationutil"

	validation "github.com/go-ozzo/ozzo-validation"
)

type User struct {
	Name    string        `json:"fullname"`
	Age     int           `json:"age"`
	Address string        `json:"address"`
	RollNo  int           `json:"rollno"`
	Courses []enum.Course `json:"courses"`
}

func New(name string, age int, address string, rollNo int, courseEnrol []string) (User, error) {
	var user User
	var err error

	user.Name = name
	user.Age = age
	user.Address = address
	user.RollNo = rollNo
	user.Courses, err = getCourse(courseEnrol)
	if err != nil {
		return User{}, err
	}

	if err = user.validate(); err != nil {
		log.Println(err)
		return User{}, err
	}

	return user, nil
}

func getCourse(courseEnrol []string) ([]enum.Course, error) {
	var courseEnum []enum.Course
	var err error
	var course enum.Course
	for _, c := range courseEnrol {
		course, err = enum.CourseString(c)
		if err != nil {
			log.Println(err)
			return []enum.Course{}, err
		}
		courseEnum = append(courseEnum, course)
	}

	return courseEnum, nil
}

func (user User) validate() error {
	return validation.ValidateStruct(&user,
		validation.Field(&user.Name, validation.Required),
		validation.Field(&user.Age, validation.Required, validation.By(validationutil.CheckPositive)),
		validation.Field(&user.Address, validation.Required),
		validation.Field(&user.RollNo, validation.Required, validation.By(validationutil.CheckPositive)),
		validation.Field(&user.Courses, validation.Required),
	)
}

func (user User) String() string {
	return fmt.Sprintf("	%s	|	%d	|	%s	|	%d	|	%s|\n", user.Name, user.Age, user.Address, user.RollNo, courseString(user.Courses))
}

func courseString(course []enum.Course) []string {
	var courses []string
	for _, c := range course {
		courses = append(courses, c.String())
	}

	return courses
}

func EncodeUser(users []User) ([]byte, error) {
	userB, err := json.Marshal(users)
	if err != nil {
		log.Println(err)
		return []byte{}, err
	}

	return userB, nil
}

func DecodeUser(userB []byte) ([]User, error) {
	var users []User
	if err := json.Unmarshal(userB, &users); err != nil {
		log.Println(err)
		return []User{}, err
	}

	return users, nil

}