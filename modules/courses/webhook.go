package courses

import (
	"encoding/json"
	"log"

	"github.com/gofiber/fiber/v2"
	"github.com/google/uuid"
	"github.com/kayprogrammer/ednet-fiber-api/config"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/ent/enrollment"
	"github.com/stripe/stripe-go/v82"
	"github.com/stripe/stripe-go/v82/webhook"
)

func StripeWebhook(db *ent.Client, cfg config.Config) fiber.Handler {
	return func(c *fiber.Ctx) error {
		body := c.Body()
		ctx := c.Context()

		if len(body) == 0 {
			return config.APIError(c, fiber.StatusBadRequest, config.ServerErr("Empty request body"))
		}
		event, err := webhook.ConstructEvent(body, c.Get("Stripe-Signature"), cfg.StripeWebhookSecret)
		if err != nil {
			return config.APIError(c, fiber.StatusBadRequest, config.ServerErr("Webhook signature verification failed"))
		}
		var session stripe.CheckoutSession
		if err := json.Unmarshal(event.Data.Raw, &session); err != nil {
			return config.APIError(c, fiber.StatusBadRequest, config.ServerErr("Failed to parse webhook JSON"))
		}
		enrollmentID, err := uuid.Parse(session.ClientReferenceID)
		if err != nil {
			return config.APIError(c, fiber.StatusBadRequest, config.ServerErr("Invalid enrollment ID"))
		}
		log.Println("Parsed Enrollment ID:", enrollmentID)
		switch event.Type {
		case "checkout.session.completed", "checkout.session.async_payment_succeeded":
			courseManager.UpdateEnrollment(db, ctx, enrollmentID, enrollment.PaymentStatusSuccessful)
		case "checkout.session.expired":
			courseManager.UpdateEnrollment(db, ctx, enrollmentID, enrollment.PaymentStatusCancelled)
		case "checkout.session.async_payment_failed":
			courseManager.UpdateEnrollment(db, ctx, enrollmentID, enrollment.PaymentStatusFailed)
		}

		return c.SendStatus(fiber.StatusOK)
	}
}
