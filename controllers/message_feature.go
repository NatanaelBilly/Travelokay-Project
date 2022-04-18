package controllers

import (
	"log"
	"strconv"
	"time"

	"github.com/Travelokay-Project/models"
	"github.com/go-co-op/gocron"
	"gopkg.in/gomail.v2"
)

func SendReceipt(emailReceiver string, newOrder models.Order, price int) {

	m := gomail.NewMessage()

	// Get value from env
	emailSender := LoadEnv("EMAIL_SENDER")
	emailPassword := LoadEnv("EMAIL_PASS")

	// Set email content
	m.SetHeader("From", emailSender)
	m.SetHeader("To", emailReceiver)
	m.SetHeader("Subject", "Travelokay Order Receipt")

	text := `<h1>Your Purchase Receipt</h1></br>
		<p>You have made a purchase via Traveloka app with the following details:</p>
		<table>
		<tr>
			<td><b>Order ID</b></td>
			<td>: ` + strconv.Itoa(newOrder.ID) + `</td>
		</tr>
		<tr>
			<td><b>Order date</b></td>
			<td>: ` + newOrder.OrderDate + `</td>
		</tr>
		<tr>
			<td><b>Order status</b></td>
			<td>: ` + newOrder.OrderStatus + `</td>
		</tr>
		<tr>
			<td><b>Transaction type</b></td>
			<td>: ` + newOrder.TransactionType + `</td>
		</tr>
		<tr>
			<td><b>Price</b></td>
			<td>: ` + strconv.Itoa(price) + `</td>
		</tr>
		</table>`

	m.SetBody("text/html", text)

	d := gomail.NewDialer("smtp.gmail.com", 465, emailSender, emailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}
func OfferMail(hari int) {
	// Connect to database
	db := Connect()
	defer db.Close()

	m := gomail.NewMessage()

	// Get value from env
	emailSender := LoadEnv("EMAIL_SENDER")
	emailPassword := LoadEnv("EMAIL_PASS")
	rows, errQuery := db.Query("SELECT email FROM users")
	if errQuery != nil {
		log.Fatal(errQuery)
		return
	}
	var user models.User
	var allEmail string
	defer rows.Close()
	for rows.Next() {
		if err := rows.Scan(&user.Email); err != nil {
			log.Fatal(errQuery)
			return
		}
		allEmail += user.Email + ","
	}
	emailReceiver := allEmail[:len(allEmail)-1]
	// Set email content
	m.SetHeader("From", emailSender)
	m.SetHeader("To", emailReceiver)
	if hari == 0 {
		m.SetHeader("Subject", "test aja")
	}
	if hari == 1 {
		m.SetHeader("Subject", "Idul Fitri Promotion Offer")
	} else if hari == 2 {
		m.SetHeader("Subject", "Christmast Promotion Offer")
	} else if hari == 3 {
		m.SetHeader("Subject", "New Year Promotion Offer")
	}

	text := "<h1>Here Is Your Best Deal Offer</h1></br>" +
		"<p><a href='#'>click here</a> to see your deal</p>"
	m.SetBody("text/html", text)

	d := gomail.NewDialer("smtp.gmail.com", 465, emailSender, emailPassword)

	// Send the email
	if err := d.DialAndSend(m); err != nil {
		log.Println(err)
	}
}

func GocronEvent() {
	s := gocron.NewScheduler(time.UTC)
	s.Cron("* /1 * * * *").Do(OfferMail, 0)    // every minutes
	s.Cron("* * * /2 /5 *").Do(OfferMail, 1)   // every idul fitri
	s.Cron("* * * /25 /12 *").Do(OfferMail, 2) // every christmast
	s.Cron("* * * /1 /1 *").Do(OfferMail, 3)   // every new year

	// starts the scheduler asynchronously
	s.StartAsync()
	// starts the scheduler and blocks current execution path
	s.StartBlocking()
}
