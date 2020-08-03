package wavefront

import (
	"fmt"
	"strings"
	"testing"

	"github.com/WavefrontHQ/go-wavefront-management-api"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/terraform"
)

func TestAccWavefrontTarget_BasicWebhook(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "description", "Test target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "method", "WEBHOOK"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "recipient", "https://hooks.slack.com/services/test/me"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "content_type", "application/json"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "template", "{}"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.0", "ALERT_OPENED"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.1", "ALERT_RESOLVED"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "custom_headers.%", "1"),
					resource.TestCheckResourceAttrSet("wavefront_alert_target.test_target", "target_id"),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_Updated(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_basic(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Target"),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_new_value(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributesUpdated(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Updated"),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_BasicEmail(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontTarget_email(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "description", "Test target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "method", "EMAIL"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "email_subject", "This is a test"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "recipient", "test@example.com"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "template", "{}"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "is_html_content", "true"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.0", "ALERT_OPENED"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.1", "ALERT_RESOLVED"),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_BasicPagerduty(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontTarget_pagerduty(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "description", "Test target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "method", "PAGERDUTY"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "recipient", "12345678910111213141516171819202"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "template", "{}"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.0", "ALERT_OPENED"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.1", "ALERT_RESOLVED"),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_AlertTargetId(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccCheckWavefrontTarget_alertTargetId(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					// Check against state that the attributes are as we expect
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "name", "Terraform Test Target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "description", "Test target"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "method", "PAGERDUTY"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "recipient", "12345678910111213141516171819202"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "template", "{}"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.#", "2"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.0", "ALERT_OPENED"),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "triggers.1", "ALERT_RESOLVED"),
					resource.TestCheckResourceAttrSet(
						"wavefront_alert_target.test_target", "target_id"),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_Routes(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_routes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "route.#", "1"),
					testAccCheckWavefrontTargetRouteAttributes(&record, []string{"prod"}),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_MultipleRoutes(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_routes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "route.#", "1"),
					testAccCheckWavefrontTargetRouteAttributes(&record, []string{"prod"}),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_addRoutes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "route.#", "2"),
					testAccCheckWavefrontTargetRouteAttributes(&record, []string{"prod", "dev"}),
				),
			},
		},
	})
}

func TestAccWavefrontTarget_UpdateRoutes(t *testing.T) {
	var record wavefront.Target

	resource.Test(t, resource.TestCase{
		PreCheck:     func() { testAccPreCheck(t) },
		Providers:    testAccProviders,
		CheckDestroy: testAccCheckWavefrontTargetDestroy,
		Steps: []resource.TestStep{
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_routes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),

					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "route.#", "1"),
					testAccCheckWavefrontTargetRouteAttributes(&record, []string{"prod"}),
				),
			},
			{
				ExpectNonEmptyPlan: true,
				Config:             testAccCheckWavefrontTarget_changeRoutes(),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckWavefrontTargetExists("wavefront_alert_target.test_target", &record),
					testAccCheckWavefrontTargetAttributes(&record),
					resource.TestCheckResourceAttr(
						"wavefront_alert_target.test_target", "route.#", "1"),
					testAccCheckWavefrontTargetRouteAttributes(&record, []string{"prod2"}),
				),
			},
		},
	})
}

func testAccCheckWavefrontTargetDestroy(s *terraform.State) error {

	targets := testAccProvider.Meta().(*wavefrontClient).client.Targets()

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "wavefront_alert_target" {
			continue
		}

		results, err := targets.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront Target. %s", err)
		}
		if len(results) > 0 {
			return fmt.Errorf("target still exists")
		}
	}

	return nil
}

func testAccCheckWavefrontTargetAttributes(target *wavefront.Target) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if target.Description != "Test target" {
			return fmt.Errorf("bad value: %s", target.Description)
		}

		return nil
	}
}

func testAccCheckWavefrontTargetRouteAttributes(target *wavefront.Target, routeTarget []string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		for _, route := range target.Routes {
			if route.Method != "WEBHOOK" {
				return fmt.Errorf("bad value: %s", route.Method)
			}

			filter := strings.Split(route.Filter, " ")
			if len(filter) != 2 {
				return fmt.Errorf("bad value: %s", route.Filter)
			}

			matchedRoute := false
			for _, target := range routeTarget {
				if strings.Contains(route.Target, target) {
					matchedRoute = true
					break
				}
			}

			if !matchedRoute {
				return fmt.Errorf("bad value: %s", route.Target)
			}
		}
		return nil
	}
}

func testAccCheckWavefrontTargetAttributesUpdated(target *wavefront.Target) resource.TestCheckFunc {
	return func(s *terraform.State) error {

		if target.Title != "Terraform Test Updated" {
			return fmt.Errorf("bad value: %s", target.Title)
		}

		return nil
	}
}

func testAccCheckWavefrontTargetExists(n string, target *wavefront.Target) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[n]

		if !ok {
			return fmt.Errorf("not found: %s", n)
		}

		if rs.Primary.ID == "" {
			return fmt.Errorf("no Record ID is set")
		}

		targets := testAccProvider.Meta().(*wavefrontClient).client.Targets()

		results, err := targets.Find(
			[]*wavefront.SearchCondition{
				{
					Key:            "id",
					Value:          rs.Primary.ID,
					MatchingMethod: "EXACT",
				},
			})
		if err != nil {
			return fmt.Errorf("error finding Wavefront Target %s", err)
		}
		// resource has been deleted out of band. So unset ID
		if len(results) != 1 {
			return fmt.Errorf("no Targets Found")
		}
		if *results[0].ID != rs.Primary.ID {
			return fmt.Errorf("target not found")
		}

		*target = *results[0]

		return nil
	}
}

func testAccCheckWavefrontTarget_basic() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Target"
	description = "Test target"
	method = "WEBHOOK"
  recipient = "https://hooks.slack.com/services/test/me"
	content_type = "application/json"
	custom_headers = {
		"Testing" = "true"
	}
	template = "{}"
	triggers = [
		"ALERT_OPENED",
		"ALERT_RESOLVED"
	]
}
`)
}

func testAccCheckWavefrontTarget_new_value() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Updated"
	description = "Test target"
	method = "WEBHOOK"
  recipient = "https://hooks.slack.com/services/test/me"
	content_type = "application/json"
	custom_headers = {
		"Testing" = "true"
	}
	template = "{}"
	triggers = [
		"ALERT_OPENED",
		"ALERT_RESOLVED"
	]
}
`)
}

func testAccCheckWavefrontTarget_routes() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
	name 		   = "Terraform Test Target"
	description    = "Test target"
	method 		   = "WEBHOOK"
	recipient	   = "https://hooks.slack.com/services/test/me"
	content_type   = "application/json"
	custom_headers = {
		"Testing" = "true"
	}
    template       = "{}"
    triggers 	   = [
		"ALERT_OPENED",
		"ALERT_RESOLVED",
	]
  	route {
		method = "WEBHOOK"
		target = "https://hooks.slack.com/services/test/me/prod"
    	filter = {
      		key   = "env"
      		value = "prod.*"
		}
  	}
}`)
}

func testAccCheckWavefrontTarget_addRoutes() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
	name 		   = "Terraform Test Target"
	description    = "Test target"
	method 		   = "WEBHOOK"
	recipient	   = "https://hooks.slack.com/services/test/me"
	content_type   = "application/json"
	custom_headers = {
		"Testing" = "true"
	}
    template       = "{}"
    triggers 	   = [
		"ALERT_OPENED",
		"ALERT_RESOLVED",
	]
  	route {
		method = "WEBHOOK"
		target = "https://hooks.slack.com/services/test/me/prod"
    	filter = {
      		key   = "env"
      		value = "prod.*"
		}
  	}
  	route {
		method = "WEBHOOK"
		target = "https://hooks.slack.com/services/test/me/dev"
    	filter = {
      		key   = "env"
      		value = "dev.*"
		}
  	}
}`)
}

func testAccCheckWavefrontTarget_changeRoutes() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
	name 		   = "Terraform Test Target"
	description    = "Test target"
	method 		   = "WEBHOOK"
	recipient	   = "https://hooks.slack.com/services/test/me"
	content_type   = "application/json"
	custom_headers = {
		"Testing" = "true"
	}
    template       = "{}"
    triggers 	   = [
		"ALERT_OPENED",
		"ALERT_RESOLVED",
	]
  	route {
		method = "WEBHOOK"
		target = "https://hooks.slack.com/services/test/me/prod2"
    	filter = {
      		key   = "env"
      		value = "prod2.*"
		}
  	}
}`)
}

func testAccCheckWavefrontTarget_email() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Target"
  description = "Test target"
  method = "EMAIL"
  recipient = "test@example.com"
  email_subject = "This is a test"
  is_html_content = true
  template = "{}"
  triggers = [
    "ALERT_OPENED",
    "ALERT_RESOLVED"
  ]
}
`)
}

func testAccCheckWavefrontTarget_pagerduty() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name = "Terraform Test Target"
	description = "Test target"
	method = "PAGERDUTY"
  recipient = "12345678910111213141516171819202"
	template = "{}"
	triggers = [
		"ALERT_OPENED",
		"ALERT_RESOLVED"
	]
}
`)
}

func testAccCheckWavefrontTarget_alertTargetId() string {
	return fmt.Sprintf(`
resource "wavefront_alert_target" "test_target" {
  name        = "Terraform Test Target"
  description = "Test target"
  method      = "PAGERDUTY"
  recipient   = "12345678910111213141516171819202"
  template    = "{}"
  triggers    = [
    "ALERT_OPENED",
	"ALERT_RESOLVED"
  ]
}

resource "wavefront_alert" "test_alert" {
  name                   = "Terraform Test Alert"
  target                 = wavefront_alert_target.test_target.target_id
  condition              = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total ) > 80"
  additional_information = "This is a Terraform Test Alert"
  display_expression     = "100-ts(\"cpu.usage_idle\", environment=preprod and cpu=cpu-total )"
  minutes                = 5
  resolve_after_minutes  = 5
  severity               = "WARN"
  tags = [
    "terraform",
    "test"
  ]
}
`)
}
