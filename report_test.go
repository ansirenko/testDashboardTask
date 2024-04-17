package testDashboardTask

import (
	"fmt"
	"github.com/chromedp/chromedp"
	"strconv"
	"time"
)

func (suite *TaskSuite) fetchFirstCellValue() (string, error) {
	var cellValue string
	err := chromedp.Run(suite.ctx,
		chromedp.Text(`div.ReactVirtualized__Grid__innerScrollContainer > div[data-row-index="0"][data-column-index="1"] span`, &cellValue, chromedp.NodeVisible, chromedp.ByQuery),
	)
	if err != nil {
		return "", err
	}
	return cellValue, nil
}

func (suite *TaskSuite) TestReferralLinkIncreaseCounter() {
	err := suite.navigateAndWait(suite.td.statisticsPage)
	suite.Require().NoError(err, "cant navigate to the statistics page")

	buttonSelector := `button[type="submit"].Button--3Qzp7`
	elementSelector := `div.table--2pEiF[qa-element="summary-table"]`
	err = suite.clickButtonAndWaitForElement(buttonSelector, elementSelector)
	suite.Require().NoError(err, "failed to click RunReport and wait for table element")

	value, err := suite.fetchFirstCellValue()
	suite.Require().NoError(err, "failed to fetch the value from the first cell")

	intValue, err := strconv.Atoi(value)
	suite.Require().NoError(err, "value in first cell is not int")

	suite.Require().GreaterOrEqual(intValue, 0, "Count of clicks unexpected")

	referralUrl, err := suite.fetchURLFromElement(suite.td.dashboardPage, `.url--c0PHO`)
	suite.Require().NoError(err, "failed to get referral url")
	suite.Require().NotEmpty(referralUrl, "Referral url isn't fetched")

	err = suite.navigateAndWait(referralUrl)
	suite.Require().NoError(err, "failed to redirect to the referral page")

	err = suite.navigateAndWait(suite.td.statisticsPage)
	suite.Require().NoError(err, "Cant navigate to the statistics page")

	err = suite.clickButtonAndWaitForElement(buttonSelector, elementSelector)
	suite.Require().NoError(err, "failed to click RunReport and wait for table element")

	err = smartWaiter(func() error {
		err = suite.clickButtonAndWaitForElement(buttonSelector, elementSelector)
		if err != nil {
			return fmt.Errorf("failed to click RunReport and wait for table element")
		}

		value, err = suite.fetchFirstCellValue()
		if err != nil {
			return err
		}
		newValue, err := strconv.Atoi(value)
		if err != nil {
			return err
		}
		if intValue == newValue {
			return fmt.Errorf("value doesn't changed\nprevious value: %v\nnew value: %v", intValue, newValue)
		}
		return nil
	}, 30*time.Second)
	suite.Require().NoError(err)
}
