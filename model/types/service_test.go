package types_test

import (
	"time"

	"github.com/soprasteria/docktor/model/types"

	. "github.com/onsi/ginkgo"
	. "github.com/onsi/gomega"
)

var _ = Describe("Service", func() {
	Context("Given a list of images", func() {
		images := []types.Image{
			{
				Created: time.Date(2005, time.November, 10, 22, 0, 0, 0, time.UTC),
				Name:    "Image2005",
			},
			{
				Created: time.Date(2016, time.November, 10, 23, 0, 0, 0, time.UTC),
				Name:    "Image2016",
			},
			{
				Created: time.Date(2015, time.November, 10, 23, 0, 0, 0, time.UTC),
				Name:    "Image2015",
			},
			{
				Created: time.Date(2010, time.November, 10, 23, 0, 0, 0, time.UTC),
				Name:    "Image2010",
			},
		}
		service := types.Service{
			Images: images,
		}
		Context("When I try to get the last image", func() {
			image, err := service.GetLatestImage()
			It("Then I should get the latest created", func() {
				Ω(err).ShouldNot(HaveOccurred())
				Ω(image.Name).Should(Equal("Image2016"))
			})
		})
	})
})