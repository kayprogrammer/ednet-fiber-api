package general

import (
	"github.com/gofiber/fiber/v2"
	"github.com/kayprogrammer/ednet-fiber-api/ent"
	"github.com/kayprogrammer/ednet-fiber-api/modules/base"
)

// @Summary Retrieve site details
// @Description This endpoint retrieves few details of the site/application.
// @Tags General
// @Success 200 {object} schemas.SiteDetailResponseSchema
// @Router /general/site-detail [get]
func GetSiteDetails(db *ent.Client) fiber.Handler {
	return func(c *fiber.Ctx) error {
		sitedetail := SiteDetailService{}.GetOrCreate(db, c.Context())
		responseSiteDetail := SiteDetailResponseSchema{
			ResponseSchema: base.ResponseMessage("Site Details Fetched!"),
			Data:           SiteDetailSchema{}.Init(sitedetail),
		}
		return c.Status(200).JSON(responseSiteDetail)
	}
}