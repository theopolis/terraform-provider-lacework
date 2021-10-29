package integration

import (
	"testing"

	"github.com/gruntwork-io/terratest/modules/terraform"
	"github.com/stretchr/testify/assert"
)

// TestServiceNowRestAlertChannelCreate applies integration terraform '../examples/resource_lacework_alert_channel_service_now'
// Uses the go-sdk to verify the created integration
// Applies an update with new channel name and Terraform destroy
func TestServiceNowRestAlertChannelCreate(t *testing.T) {
	terraformOptions := terraform.WithDefaultRetryableErrors(t, &terraform.Options{
		TerraformDir: "../examples/resource_lacework_alert_channel_service_now",
		Vars: map[string]interface{}{
			"channel_name":         "Service Now Alert Channel Example",
			"instance_url":         "https://dev123.service-now.com",
			"username":             "snow-user",
			"password":             "snow-pass",
			"custom_template_file": customTemplate,
		},
	})
	defer terraform.Destroy(t, terraformOptions)

	// Create new ServiceNowRest Alert Channel
	create := terraform.InitAndApplyAndIdempotent(t, terraformOptions)
	created := GetAlertChannelProps(create)
	if data, ok := created.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)

		actualName := terraform.Output(t, terraformOptions, "channel_name")
		actualUrl := terraform.Output(t, terraformOptions, "instance_url")
		actualUsername := terraform.Output(t, terraformOptions, "username")
		actualTemplate := terraform.Output(t, terraformOptions, "custom_template_file")

		assert.Equal(t, "Service Now Alert Channel Example", created.Data.Name)
		assert.Equal(t, "https://dev123.service-now.com", data["instanceUrl"])
		assert.Equal(t, templateEncoded,
			data["customTemplateFile"])
		assert.Equal(t, "snow-user", data["userName"])

		assert.Equal(t, "Service Now Alert Channel Example", actualName)
		assert.Equal(t, "https://dev123.service-now.com", actualUrl)
		assert.Equal(t, "snow-user", actualUsername)
		assert.Equal(t, customTemplate, actualTemplate)
	}
	// Update ServiceNowRest Alert Channel
	terraformOptions.Vars = map[string]interface{}{
		"channel_name": "Service Now Alert Channel Updated",
		"instance_url": "https://dev321.service-now.com",
		"username":     "snow-user-updated",
	}

	update := terraform.Apply(t, terraformOptions)
	updated := GetAlertChannelProps(update)
	if data, ok := updated.Data.Data.(map[string]interface{}); ok {
		assert.True(t, ok)

		actualName := terraform.Output(t, terraformOptions, "channel_name")
		actualUrl := terraform.Output(t, terraformOptions, "instance_url")
		actualUsername := terraform.Output(t, terraformOptions, "username")
		actualTemplate := terraform.Output(t, terraformOptions, "custom_template_file")

		assert.Equal(t, "Service Now Alert Channel Updated", updated.Data.Name)
		assert.Equal(t, "https://dev321.service-now.com", data["instanceUrl"])
		assert.Equal(t, templateEncoded, data["customTemplateFile"])
		assert.Equal(t, "snow-user-updated", data["userName"])

		assert.Equal(t, "Service Now Alert Channel Updated", actualName)
		assert.Equal(t, "https://dev321.service-now.com", actualUrl)
		assert.Equal(t, "snow-user-updated", actualUsername)
		assert.Equal(t, customTemplate, actualTemplate)
	}
}

var customTemplate = "  {\n    \"description\" : \"Generated by Lacework:\",\n    \"approval\" : \"Approved\"\n  }\n"
var templateEncoded = "data:application/json;name=i.json;base64,ICB7CiAgICAiZGVzY3JpcHRpb24iIDogIkdlbmVyYXRlZCBieSBMYWNld29yazoiLAogICAgImFwcHJvdmFsIiA6ICJBcHByb3ZlZCIKICB9Cg=="