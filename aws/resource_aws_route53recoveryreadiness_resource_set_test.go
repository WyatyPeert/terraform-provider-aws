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

func TestAccAwsRoute53RecoveryReadinessResourceSet_basic(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetConfig(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`resource-set.+`)),
					resource.TestCheckResourceAttr(resourceName, "resources.#", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "0"),
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

func TestAccAwsRoute53RecoveryReadinessResourceSet_tags(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetConfig_Tags1(rName, "key1", "value1"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
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
				Config: testAccAwsRoute53RecoveryReadinessResourceSetConfig_Tags2(rName, "key1", "value1updated", "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "2"),
					resource.TestCheckResourceAttr(resourceName, "tags.key1", "value1updated"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetConfig_Tags1(rName, "key2", "value2"),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					resource.TestCheckResourceAttr(resourceName, "tags.%", "1"),
					resource.TestCheckResourceAttr(resourceName, "tags.key2", "value2"),
				),
			},
		},
	})
}

func TestAccAwsRoute53RecoveryReadinessResourceSet_readinessScope(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck:          func() { testAccPreCheck(t); testAccPreCheckAwsRoute53RecoveryReadiness(t) },
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetConfig_ReadinessScopes(rName),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`resource-set.+`)),
					resource.TestCheckResourceAttr(resourceName, "resources.0.readiness_scopes.#", "1"),
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

func TestAccAwsRoute53RecoveryReadinessResourceSet_basicDnsTargetResource(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	domainName := "myTestDomain.test"
	hzArn := "arn:aws:route53::01234567890:hostedzone/ZZZZZZZZZZZZZZ"
	recordType := "A"
	recordSetId := "12345"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAwsRoute53RecoveryReadiness(t)
		},
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetBasicDnsTargetResourceConfig(rName, domainName, hzArn, recordType, recordSetId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`resource-set.+`)),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.hosted_zone_arn", hzArn),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.record_type", recordType),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.record_set_id", recordSetId),
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

func TestAccAwsRoute53RecoveryReadinessResourceSet_DnsTargetResourceNlbTarget(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	nlbArn := "arn:aws:elasticloadbalancing:us-east-2:123456789012:loadbalancer/net/my-load-balancer/1234567890123456"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAwsRoute53RecoveryReadiness(t)
		},
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetDnsTargetResourceNlbTargetConfig(rName, nlbArn),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`resource-set.+`)),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.target_resource.0.nlb_resource.0.arn", nlbArn),
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

func TestAccAwsRoute53RecoveryReadinessResourceSet_DnsTargetResourceR53Target(t *testing.T) {
	rName := acctest.RandomWithPrefix("tf-acc-test")
	resourceName := "aws_route53recoveryreadiness_resource_set.test"
	domainName := "my.target.domain"
	recordSetId := "987654321"
	resource.ParallelTest(t, resource.TestCase{
		PreCheck: func() {
			testAccPreCheck(t)
			testAccPreCheckAwsRoute53RecoveryReadiness(t)
		},
		ErrorCheck:        testAccErrorCheck(t, route53recoveryreadiness.EndpointsID),
		ProviderFactories: testAccProviderFactories,
		CheckDestroy:      testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy,
		Steps: []resource.TestStep{
			{
				Config: testAccAwsRoute53RecoveryReadinessResourceSetDnsTargetResourceR53TargetConfig(rName, domainName, recordSetId),
				Check: resource.ComposeTestCheckFunc(
					testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(resourceName),
					testAccMatchResourceAttrGlobalARN(resourceName, "arn", "route53-recovery-readiness", regexp.MustCompile(`resource-set.+`)),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.target_resource.0.r53_resource.0.domain_name", domainName),
					resource.TestCheckResourceAttr(resourceName, "resources.0.dns_target_resource.0.target_resource.0.r53_resource.0.record_set_id", recordSetId),
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

func testAccCheckAwsRoute53RecoveryReadinessResourceSetDestroy(s *terraform.State) error {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	for _, rs := range s.RootModule().Resources {
		if rs.Type != "aws_route53recoveryreadiness_resource_set" {
			continue
		}

		input := &route53recoveryreadiness.GetResourceSetInput{
			ResourceSetName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetResourceSet(input)
		if err == nil {
			return fmt.Errorf("Route53RecoveryReadiness Resource Set (%s) not deleted", rs.Primary.ID)
		}
	}

	return nil
}

func testAccCheckAwsRoute53RecoveryReadinessResourceSetExists(name string) resource.TestCheckFunc {
	return func(s *terraform.State) error {
		rs, ok := s.RootModule().Resources[name]
		if !ok {
			return fmt.Errorf("Not found: %s", name)
		}

		conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

		input := &route53recoveryreadiness.GetResourceSetInput{
			ResourceSetName: aws.String(rs.Primary.ID),
		}

		_, err := conn.GetResourceSet(input)

		return err
	}
}

func testAccPreCheckAwsRoute53RecoveryReadinessResourceSet(t *testing.T) {
	conn := testAccProvider.Meta().(*AWSClient).route53recoveryreadinessconn

	input := &route53recoveryreadiness.ListResourceSetsInput{}

	_, err := conn.ListResourceSets(input)

	if testAccPreCheckSkipError(err) {
		t.Skipf("skipping acceptance testing: %s", err)
	}

	if err != nil {
		t.Fatalf("unexpected PreCheck error: %s", err)
	}
}

func testAccAwsCloudWatchMetricAlarmForResourceSetConfig() string {
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
`)
}

func testAccAwsRoute53RecoveryReadinessResourceSetConfig(rName string) string {
	return composeConfig(testAccAwsCloudWatchMetricAlarmForResourceSetConfig(), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %q
  resource_set_type = "AWS::CloudWatch::Alarm"
  
  resources {
	resource_arn = aws_cloudwatch_metric_alarm.test.arn
  }
}
`, rName))
}

func testAccAwsRoute53RecoveryReadinessResourceSetConfig_Tags1(rName, tagKey1, tagValue1 string) string {
	return composeConfig(testAccAwsCloudWatchMetricAlarmForResourceSetConfig(), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %[1]q
  resource_set_type = "AWS::CloudWatch::Alarm"
  resources {
	resource_arn = aws_cloudwatch_metric_alarm.test.arn
  }
  tags = {
    %[2]q = %[3]q
  }
}
`, rName, tagKey1, tagValue1))
}

func testAccAwsRoute53RecoveryReadinessResourceSetConfig_Tags2(rName, tagKey1, tagValue1, tagKey2, tagValue2 string) string {
	return composeConfig(testAccAwsCloudWatchMetricAlarmForResourceSetConfig(), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %[1]q
  resource_set_type = "AWS::CloudWatch::Alarm"
  resources {
	resource_arn = aws_cloudwatch_metric_alarm.test.arn
  }
  tags = {
    %[2]q = %[3]q
    %[4]q = %[5]q
  }
}
`, rName, tagKey1, tagValue1, tagKey2, tagValue2))
}

func testAccAwsRoute53RecoveryReadinessResourceSetConfig_ReadinessScopes(rName string) string {
	return composeConfig(testAccAwsCloudWatchMetricAlarmForResourceSetConfig(), fmt.Sprintf(`
resource "aws_route53recoveryreadiness_cell" "test" {
	cell_name = "resource_set_test_cell"
}

resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %q
  resource_set_type = "AWS::CloudWatch::Alarm"
  
  resources {
	resource_arn = aws_cloudwatch_metric_alarm.test.arn
	readiness_scopes = [aws_route53recoveryreadiness_cell.test.arn]
  }
}
`, rName))
}

func testAccAwsRoute53RecoveryReadinessResourceSetBasicDnsTargetResourceConfig(rName, domainName, hzArn, recordType, recordSetId string) string {
	return fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %[1]q
  resource_set_type = "AWS::Route53RecoveryReadiness::DNSTargetResource"
  
  resources {
	dns_target_resource {
		domain_name = %[2]q
		hosted_zone_arn = %[3]q
		record_type = %[4]q
		record_set_id = %[5]q
	}
  }
}
`, rName, domainName, hzArn, recordType, recordSetId)
}

func testAccAwsRoute53RecoveryReadinessResourceSetDnsTargetResourceNlbTargetConfig(rName, nlbArn string) string {
	return fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %[1]q
  resource_set_type = "AWS::Route53RecoveryReadiness::DNSTargetResource"
  
  resources {
	dns_target_resource {
		domain_name = "myTestDomain.test"
		hosted_zone_arn = "arn:aws:route53::01234567890:hostedzone/ZZZZZZZZZZZZZZ"
		record_type = "A"
		record_set_id = "12345"
		
		target_resource {
			nlb_resource {
				arn = %[2]q
			}
		}
	}
  }
}
`, rName, nlbArn)
}

func testAccAwsRoute53RecoveryReadinessResourceSetDnsTargetResourceR53TargetConfig(rName, domainName, recordSetId string) string {
	return fmt.Sprintf(`
resource "aws_route53recoveryreadiness_resource_set" "test" {
  resource_set_name = %[1]q
  resource_set_type = "AWS::Route53RecoveryReadiness::DNSTargetResource"
  
  resources {
	dns_target_resource {
		domain_name = "myTestDomain.test"
		hosted_zone_arn = "arn:aws:route53::01234567890:hostedzone/ZZZZZZZZZZZZZZ"
		record_type = "A"
		record_set_id = "12345"
		
		target_resource {
			r53_resource {
				domain_name = %[2]q
				record_set_id = %[3]q
			}
		}
	}
  }
}
`, rName, domainName, recordSetId)
}
