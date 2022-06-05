package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"unicode"

	"github.com/gin-gonic/gin"
	"go.mongodb.org/mongo-driver/bson"
)

//UPing ...
func UPing(c *gin.Context) {
	c.JSON(200, gin.H{
		"msg": "Server up and running! Fine :)",
	})
}

//UCreateAccount ...
func UCreateAccount(c *gin.Context) {
	accountType := c.Query("is_teacher")
	email := c.Query("email")

	usersCollection := mpDatabaseRS.Collection("users")

	//Validation starts here
	var canCreateAccount bool = true

	//Checking email
	filter, err := usersCollection.Find(ctx, bson.M{"email": email})
	if err != nil {
		log.Fatal(err)
	}

	var playersFiltered []bson.M
	if err = filter.All(ctx, &playersFiltered); err != nil {
		log.Fatal(err)
	}
	if playersFiltered != nil {
		canCreateAccount = false
	}

	if !canCreateAccount {
		fmt.Println("[SERVER] Can't create account, email already exist.", email)
		c.JSON(200, gin.H{
			"msg": "An account with this email is already in system.",
		})
	}

	//Inserting...
	//Creating document
	if canCreateAccount {
		MAIuser := MAIUSER{
			Name:        "",
			Email:       email,
			Token:       "",
			Initials:    "",
			Discipline:  "",
			Disciplines: "",
			Marks:       "",
			Isteacher:   false,
		}
		insertResult, err := usersCollection.InsertOne(ctx, MAIuser)
		if accountType == "no" { //no - student, yes - teacher
			name := c.Query("name")
			disciplines := c.Query("disciplines")
			_, err = usersCollection.UpdateOne(ctx,
				bson.M{"email": email},
				bson.D{
					{"$set", bson.D{{"name", name}}},
					{"$set", bson.D{{"disciplines", disciplines}}},
				},
			)

		} else {
			initials := c.Query("initials")
			discipline := c.Query("discipline")
			_, err = usersCollection.UpdateOne(ctx,
				bson.M{"email": email},
				bson.D{
					{"$set", bson.D{{"initials", initials}}},
					{"$set", bson.D{{"discipline", discipline}}},
					{"$set", bson.D{{"isteacher", true}}},
				},
			)
		}

		if err != nil {
			log.Fatal(err)
		}

		fmt.Println("[SERVER] Created account with ID", insertResult.InsertedID)

		//Generating auth token for user
		uToken, err := GenerateJWT(email, accountType)

		_, _ = usersCollection.UpdateOne(ctx,
			bson.M{"email": email},
			bson.D{
				{"$set", bson.D{{"Token", uToken}}},
			},
		)

		c.JSON(200, gin.H{
			"msg":   "Success: You created account! Your token is written above. Please save it, otherwise you will lose access to your account.",
			"Token": uToken,
		})
	}
}

//ActionTeacher ...
func ActionTeacher(c *gin.Context) {
	userToken := string(c.GetHeader("Token"))
	var teacher MAIUSER
	var allok bool = true

	usersCollection := mpDatabaseRS.Collection("users")

	//Getting info
	err := usersCollection.FindOne(ctx, bson.M{"Token": userToken}).Decode(&teacher)
	if err != nil {
		log.Print(err)
		c.JSON(200, gin.H{
			"msg": "Your token is invalid.",
		})
		allok = false
	}
	action := c.Query("action")
	if action == "add_mark" {
		var student MAIUSER
		student_name := c.Query("student_name")
		err := usersCollection.FindOne(ctx, bson.M{"name": student_name}).Decode(&student)
		if err != nil {
			fmt.Printf("[SERVER] Teacher entered invalid student name")
			c.JSON(200, gin.H{
				"msg": "Something went wrong... Please check if student name is valid",
			})
			allok = false
		}

		//Checking if teacher discipline is in student's discipline list; inserting mark
		if strings.Contains(student.Disciplines, teacher.Discipline) && allok {
			mark_type := c.Query("mark_type")
			mark := c.Query("mark")
			coef := c.Query("coef")
			cf, _ := strconv.Atoi(coef)
			if coef == "" || cf < 1 {
				cf = 1
			}
			existing_marks := student.Marks
			new_marks := ""
			var u int = 0
			if existing_marks == "" {
				new_marks = existing_marks + teacher.Discipline + "-" + mark_type + "-" + mark
			} else {
				new_marks = existing_marks + " " + teacher.Discipline + "-" + mark_type + "-" + mark
			}
			for u < cf-1 {
				new_marks = new_marks + " " + teacher.Discipline + "-" + mark_type + "-" + mark
				u++
			}
			_, err = usersCollection.UpdateOne(ctx,
				bson.M{"name": student_name},
				bson.D{
					{"$set", bson.D{{"marks", new_marks}}},
				},
			)

			fmt.Println("[SERVER] Teacher " + teacher.Initials + " added a new mark " + mark + " to " + student_name)
			if cf == 1 {
				c.JSON(200, gin.H{
					"msg": "Success. You've added " + mark + " to " + student_name + "'s marks list.",
				})
			} else {
				c.JSON(200, gin.H{
					"msg": "Success. You've added " + mark + " (x" + coef + ")" + " to " + student_name + "'s marks list.",
				})
			}

		} else {
			if allok {
				fmt.Println("[SERVER] Teacher tried to access a student he/she does not teach.")
				c.JSON(200, gin.H{
					"msg": "You are not teaching this student. Please check if you have entered valid student name. Access denied.",
				})
			}
		}
	} else {
		c.JSON(200, gin.H{
			"msg": "No such action. Please ensure you've entered valid action name",
		})
	}
}

//ArContains: if x string is substring of any ar element.
func ArContains(ar []string, x string) bool {
	for i := 0; i < len(ar); i++ {
		if strings.Contains(ar[i], x) {
			return true
		}
	}
	return false
}

//ArContainsWithIndex: if x string is substring of any ar element - return element index
func ArContainsWithIndex(ar []string, x string) int {
	for i := 0; i < len(ar); i++ {
		if strings.Contains(ar[i], x) {
			return i
		}
	}
	return -1
}

//GetArElementBySubstring: if x string is substring of any ar element - return all such elements
func GetArElementsBySubstring(ar []string, x string) []string {
	var res []string
	for i := 0; i < len(ar); i++ {
		if strings.Contains(ar[i], x) {
			res = append(res, ar[i])
		}
	}
	return res
}

//RepalceArElementBySubstring: if substring string is substring of any ar element - replace this element with new string
func ReplaceArElementBySubstring(ar []string, substring string, new string) []string {
	for i := 0; i < len(ar); i++ {
		if strings.Contains(ar[i], substring) {
			ar[i] = new
			break
		}
	}
	return ar
}

//GetDisciplinesListFromAr: returns all unique disciplines from marks array
func GetDisciplinesListFromAr(ar []string) []string {
	var d []string
	for i := 0; i < len(ar); i++ {
		if !ArContains(d, strings.Split(ar[i], " ")[0]) {
			d = append(d, strings.Split(ar[i], " ")[0])
		}
	}
	return d
}

//ActionStudent ...
func ActionStudent(c *gin.Context) {
	userToken := string(c.GetHeader("Token"))
	var student MAIUSER

	usersCollection := mpDatabaseRS.Collection("users")

	err := usersCollection.FindOne(ctx, bson.M{"Token": userToken}).Decode(&student)
	if err != nil {
		log.Print(err)
		c.JSON(200, gin.H{
			"msg": "Your token is invalid.",
		})
	}

	action := c.Query("action")
	if action == "get_marks" {
		var anyMark bool = true
		fmt.Println("[SERVER] Sending marks to " + student.Name)
		subject := c.Query("subject")
		marks := student.Marks
		if marks == "" {
			c.JSON(200, gin.H{
				"msg": "You don't have any marks yet.",
			})
			anyMark = false
		}
		if anyMark {
			marks_arr := strings.Split(marks, " ") //Math-Exam-5 Math-Test-4...
			var marks_out []string

			//Firstly adding all subjects to marks_out...
			for i := 0; i < len(marks_arr); i++ {
				marks_parts := strings.Split(marks_arr[i], "-")
				if !ArContains(marks_out, marks_parts[0]) {
					marks_out = append(marks_out, marks_parts[0])
					if ArContains(marks_arr, marks_parts[0]+"-Exam") && ArContains(marks_arr, marks_parts[0]+"-Test") {
						marks_out = append(marks_out, marks_parts[0])
					}
					var hasExam bool = false
					if ArContains(marks_arr, marks_parts[0]+"-Exam") {
						ind := ArContainsWithIndex(marks_out, marks_parts[0])
						marks_out[ind] = marks_out[ind] + " Exam"
						hasExam = true
					}
					if ArContains(marks_arr, marks_parts[0]+"-Test") {
						ind := ArContainsWithIndex(marks_out, marks_parts[0])
						if hasExam {
							marks_out[ind+1] = marks_out[ind+1] + " Test"
						} else {
							marks_out[ind] = marks_out[ind] + " Test"
						}
					}
				}
			}

			//After that marks_out contains the following: Math Exam Math Test Russian Test...
			for i := 0; i < len(marks_arr); i++ {
				marks_parts := strings.Split(marks_arr[i], "-") //"Math", "Test", "5"
				if ArContains(marks_out, marks_parts[0]) {      //marks_out: Math Exam 5 4 5 5;;  Math Test 4 5 5;; Russian Test 5 5 5;;
					part := GetArElementsBySubstring(marks_out, marks_parts[0]) //like above, gets elements only with "Math" or "English" or etc.
					mark_type := marks_parts[1]                                 //"Test" or "Exam"
					var sw bool = false
					if strings.Contains(part[0], mark_type) {
						part[0] += " " + marks_parts[2] //added mark
						sw = false
					} else if strings.Contains(part[1], mark_type) {
						part[1] += " " + marks_parts[2] //added mark
						sw = true
					}

					//Adding part[0] or part[1] to marks_out
					if sw {
						marks_out = ReplaceArElementBySubstring(marks_out, marks_parts[0]+" "+marks_parts[1], part[1])
					} else {
						marks_out = ReplaceArElementBySubstring(marks_out, marks_parts[0]+" "+marks_parts[1], part[0])
					}
				}
			}

			//Marks ready
			//Calculating mean
			var marks_mean []string
			var sum int = 0
			var counter int = 0
			d := GetDisciplinesListFromAr(marks_out)
			for t := 0; t < len(d); t++ {
				sum = 0
				counter = 0
				marks_by_subj := GetArElementsBySubstring(marks_out, d[t]) //Getting marks strings by subject
				for k := 0; k < len(marks_by_subj); k++ {                  //math math math
					mar := strings.Split(marks_by_subj[k], " ")
					for g := 0; g < len(mar); g++ {
						if unicode.IsDigit([]rune(mar[g])[0]) { //if mar[g] is digit, add it to sum
							bufer, _ := strconv.Atoi(mar[g])
							sum += bufer
							counter++
						}
					}
				}
				var mean float64 = float64(sum) / float64(counter) //Calculating mean
				marks_mean = append(marks_mean, d[t]+": "+strconv.FormatFloat(mean, 'f', 1, 32))
			}

			//Sending marks to user
			if subject == "" { //if filter subject wasn't specified in request
				c.JSON(200, gin.H{
					"Your marks": marks_out,
					"Marks mean": marks_mean,
				})
				fmt.Println("[SERVER] Successfully processed marks for " + student.Name)
			} else { //if user specified subject, which marks he/she wants to see
				response := GetArElementsBySubstring(marks_out, subject)
				responseName := "Your " + subject + " marks"
				responseMean := GetArElementsBySubstring(marks_mean, subject)
				c.JSON(200, gin.H{
					responseName:      response,
					subject + " mean": responseMean,
				})
				fmt.Println("[SERVER] Successfully processed marks for " + student.Name)
			}
		}
	} else {
		c.JSON(200, gin.H{
			"msg": "No such action. Please ensure you've entered valid action name",
		})
	}
}
