package aws

import (
	"fmt"
	"regexp"
	"testing"

	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/service/route53recoveryreadiness"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/acctest"
	"github.com/hashicorp/terraform-plugin-sdk/v2/helper/resource"
	"github.com/hashicorp/terraform-plugin-sdk/v2/terraform"
)

func TestAccAwsRoute53RecoveryReadinessReadinessCheck_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	rSetName := acctest.RandomWithPrefix("tf-acc-test-set")
	resourceName := "aws_route53recoveryreadiness_readiness_check.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessReadinessCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessReadinessCheckConfig(rName, rSetName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessReadinessCheckExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`readiness-check/.+`)),
					resource.TestCheckResourceAttr(resourceName, "resource_set_name", rSetName),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
		},
	})
}

func TestAccAwsRoute53RecoveryReadinessReadinessCheck_tags(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_readiness_check.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessReadinessCheckDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessReadinessCheckConfig_Tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessReadinessCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1"),
				),
			},
			{
				ResourceName:      resourceName,
				ImportState:       true,
				ImportStateVerify: true,
			},
			{
				Config: testAccAwsRoute53RecoveryReadinessReadinessCheckConfig_Tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessReadinessCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAwsRoute53RecoveryReadinessReadinessCheckConfig_Tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessReadinessCheckExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func testAccCheckAwsRoute53RecoveryReadinessReadinessCheckDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53recoveryreadiness_readiness_check" {
			continue
		}

		input := &route53recoveryreadiness.GetReadinessCheckInput{
			ReadinessCheckName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetReadinessCheck(input)
		if err == nil {
			return fmt.Errorf("Route53RecoveryReadiness Readiness Check (%s) not deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAwsRoute53RecoveryReadinessReadinessCheckExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

		input := &route53recoveryreadiness.GetReadinessCheckInput{
			ReadinessCheckName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetReadinessCheck(input)

		return err
	}
}

func testAccPreCheckAwsRoute53RecoveryReadinessReadinessCheck(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	input := &route53recoveryreadiness.ListReadinessChecksInput{}

	_, err := conn.ListReadinessChecks(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAwsCloudWatchAlarmResourceSetConfig(rSetName string) string {
	return fmt.Sprintf(`
resource "aws_cloudwatch_metric_alarm" "test" {
  alarm_name                = "ResourceSetTestAlarm"
  comparison_operator       = "GreaterThanOrEqualToThreshold"
  evaluation_periods        = "2"
  metric_name               = "CPUUtilization"
  namespace                 = "AWS/EC2"
  period                    = "120"
  statistic                 = "Average"
  threshold                 = "80"
  alarm_description         = "This metric monitors ec2 cpu utilization"
  insufficient_data_actions = []

  dimensions = {
    InstanceId = "i-abc123"
  }
}

resource "aws_route53recoveryreadiness_resource_set" "test" {
	resource_set_name = %q
	resource_set_type = "AWS::CloudWatch::Alarm"
	
	resources {
	  resource_arn = aws_cloudwatch_metric_alarm.test.arn
	}
  }
`, rSetName)
}

func testAccAwsRoute53RecoveryReadinessReadinessCheckConfig(rName, rSetName string) string {
	return composeConfig(testAccAwsCloudWatchAlarmResourceSetConfig(rSetName), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_readiness_check" "test" {
  readiness_check_name = %q
  resource_set_name = aws_route53recoveryreadiness_resource_set.test.resource_set_name
}
`, rName))
}

func testAccAwsRoute53RecoveryReadinessReadinessCheckConfig_Tags1(rName, tagKey1, tagValue1 string) string {
	return composeConfig(testAccAwsCloudWatchAlarmResourceSetConfig("resource-set-for-testing"), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_readiness_check" "test" {
  readiness_check_name = %[1]q
  resource_set_name = aws_route53recoveryreadiness_resource_set.test.resource_set_name
  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccAwsRoute53RecoveryReadinessReadinessCheckConfig_Tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return composeConfig(testAccAwsCloudWatchAlarmResourceSetConfig("resource-set-for-testing"), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_readiness_check" "test" {
  readiness_check_name = %[1]q
  resource_set_name = aws_route53recoveryreadiness_resource_set.test.resource_set_name
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}
