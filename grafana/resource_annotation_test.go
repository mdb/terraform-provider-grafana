package grafana

import (
	"fmt"
	"net/url"
	"regexp"
	"testing"

	gapi "github.com/nytm/go-grafana-api"

	"github.com/hashicorp/terraform/helper/resource"
	"github.com/hashicorp/terraform/terraform"
)

func TestAccAnnotation_basic(t *testing.T) {
	var annotation gapi.Annotation

	resource.Test(t, resource.TestCase{
		PreCheck:  func() { testAccPreCheck(t) },
		Providers: testAccProviders,
		//CheckDestroy: testAccAnnotationCheckDestroy(&annotation),
		Steps: []resource.TestStep{
			// first step creates the resource
			{
				Config: testAccAnnotationConfig_basic,
				Check: resource.ComposeTestCheckFunc(
					testAccAnnotationCheckExists("grafana_annotation.test", &annotation),
					resource.TestMatchResourceAttr(
						"grafana_annotation.test", "text", regexp.MustCompile(".*Terraform Acceptance Test.*"),
					),
				),
			},
			// second step updates it with a new text
			{
				Config: testAccAnnotationConfig_update,
				Check: resource.ComposeTestCheckFunc(
					testAccAnnotationCheckExists("grafana_annotation.test", &annotation),
					resource.TestMatchResourceAttr(
						"grafana_annotation.test", "text", regexp.MustCompile(".*Updated Text.*"),
					),
				),
			},
			// final step checks importing the current state we reached in the step above
			{
				ResourceName:      "grafana_annotation.test",
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func testAccAnnotationCheckExists(rn string, annotation *gapi.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[rn]
		if !ok {
			return fmt.Errorf("resource not found: %s", rn)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("resource id not set")
		}

		client := testAccProvider.Meta().(*gapi.Client)
		as, err := client.Annotations(url.Values{})
		if err == nil {
			return err
		}

		for _, a := range as {
			if a.ID == annotation.ID {
				return nil
			}
		}

		return fmt.Errorf("failed to find annotation ID: %v", annotation.ID)
	}
}

func testAccAnnotationCheckDestroy(annotation *gapi.Annotation) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		client := testAccProvider.Meta().(*gapi.Client)
		as, err := client.Annotations(url.Values{})
		if err == nil {
			return err
		}

		for _, a := range as {
			if a.ID == annotation.ID {
				return fmt.Errorf("failed to destroy annotation ID: %v", annotation.ID)
			}
		}

		return nil
	}
}

const testAccAnnotationConfig_basic = `
resource "grafana_annotation" "test" {
	text = "Terraform Acceptance Test"
}
`

const testAccAnnotationConfig_update = `
resource "grafana_annotation" "test" {
	text = "Updated Text"
}
`
