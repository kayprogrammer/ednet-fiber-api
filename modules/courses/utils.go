package courses

import (
	"context"
	"fmt"
	"log"
	"net/url"

	"github.com/gofiber/fiber/v2"
	"github.com/gosimple/slug"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/course"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/checkout/session"
)

// GenerateCourseSlug creates a unique slug with a 7-character random suffix.
func GenerateCourseSlug(db *ent.Client, ctx context.Context, name string) string {
	baseSlug := slug.Make(name)
	courseSlug := baseSlug

	const maxAttempts = 5
	for i := 0; i < maxAttempts; i++ {
		// Check if slug exists
		course, _ := db.Course.Query().Where(course.Slug(courseSlug)).Only(ctx)
		if course == nil {
			break
		}
		// Generate a new slug with random 7-character suffix
		courseSlug = fmt.Sprintf("%s-%s", baseSlug, config.GetRandomString(7))
	}
	return courseSlug
}

func GetCurrentOrigin(c *fiber.Ctx) string {
	origin := c.Get("Origin")
	if origin == "" {
		// fallback to Referer if Origin is not set
		referer := c.Get("Referer")
		if referer != "" {
			// extract origin from Referer if needed
			u, err := url.Parse(referer)
			if err == nil {
				origin = u.Scheme + "://" + u.Host
			}
		}
	}
	return origin
}

func CreateCheckoutSession(cfg config.Config, course *ent.Course, successUrl string, cancelUrl string, enrollmentObj *ent.Enrollment) (*string, *config.ErrorResponse) {
	stripe.Key = cfg.StripeSecretKey

	price := course.DiscountPrice
	if price == 0.0 {
		price = course.Price
	}
	log.Println(price)
	params := &stripe.CheckoutSessionParams{
		ClientReferenceID: stripe.String(enrollmentObj.ID.String()),
		PaymentMethodTypes: stripe.StringSlice([]string{"card"}),
		Mode:               stripe.String("payment"),
		LineItems: []*stripe.CheckoutSessionLineItemParams{
			{
				PriceData: &stripe.CheckoutSessionLineItemPriceDataParams{
					Currency: stripe.String("usd"),
					ProductData: &stripe.CheckoutSessionLineItemPriceDataProductDataParams{
						Name: stripe.String(course.Title),
						Description: stripe.String(course.Desc),
						Images: stripe.StringSlice([]string{course.ThumbnailURL}),
					},
					UnitAmount: stripe.Int64(int64(price * 100)), // e.g., 5000 = $50.00
				},
				Quantity: stripe.Int64(1),
			},
		},
		SuccessURL: stripe.String(successUrl),
		CancelURL:  stripe.String(cancelUrl),
	}

	s, err := session.New(params)
	if err != nil {
		log.Println("Stripe error: ", err)
		err := config.RequestErr(config.ERR_SERVER_ERROR, "Something went wrong")
		return nil, &err
	}
	return &s.URL, nil
}
