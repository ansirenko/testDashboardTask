package testDashboardTask

import (
	"context"
	"fmt"
	"github.com/chromedp/chromedp"
	"github.com/stretchr/testify/suite"
	"testing"
	"time"
)

type TaskSuite struct {
	suite.Suite
	ctx       context.Context
	ctxCancel func()
	td        *TaskData
}

type TaskData struct {
	loginPage      string
	login          string
	password       string
	dashboardPage  string
	nickname       string
	statisticsPage string
}

func (suite *TaskSuite) SetupSuite() {
	suite.ctx, suite.ctxCancel = chromedp.NewContext(context.Background())

	suite.ctx, suite.ctxCancel = context.WithTimeout(suite.ctx, 5*time.Minute)

	suite.initializeTestData()

	suite.authorization()
}

func (suite *TaskSuite) SetupTest() {

}

func TestSuite(t *testing.T) {
	// Run the test suite
	suite.Run(t, new(TaskSuite))
}

func (suite *TaskSuite) TearDownTest() {

}

func (suite *TaskSuite) TearDownSuite() {
	err := chromedp.Run(suite.ctx,
		chromedp.Click(`button[qa-element="header-dropdown"]`, chromedp.NodeVisible),
		chromedp.WaitVisible(`div[role="option"][aria-label="Logout"]`, chromedp.ByQuery),
		chromedp.Click(`div[role="option"][aria-label="Logout"]`, chromedp.NodeVisible),
	)
	suite.Require().NoError(err, "Logout failed")

	suite.ctxCancel()
}

func (suite *TaskSuite) initializeTestData() {
	suite.td = &TaskData{
		loginPage:      "https://stripcash.com/login",
		login:          "revenueShareQA@stripdev.com",
		password:       "revenueShareQA@stripdev.com",
		dashboardPage:  "https://stripcash.com/overview/dashboard",
		nickname:       "revenueShareQA",
		statisticsPage: "https://stripcash.com/analytics/statistics",
	}
}

func (suite *TaskSuite) authorization() {
	var currentURL string
	err := chromedp.Run(suite.ctx,
		chromedp.Navigate(suite.td.loginPage),
		chromedp.WaitVisible(`input[name="username"]`, chromedp.ByQuery),
		chromedp.SendKeys(`input[name="username"]`, suite.td.login, chromedp.ByQuery),
		chromedp.SendKeys(`#password`, suite.td.password, chromedp.ByID),
		chromedp.Click(`button[type="submit"]`, chromedp.ByQuery),
		chromedp.Location(&currentURL),
	)
	suite.Require().NoError(err, "failed to login")

	err = suite.waitForPageAvailable(suite.td.dashboardPage)
	suite.Require().NoError(err)

	suite.checkDashboard()
}

func (suite *TaskSuite) waitForPageAvailable(urlExpected string) error {
	err := smartWaiter(func() error {
		var url = urlExpected
		err := chromedp.Run(suite.ctx,
			chromedp.Location(&url),
		)
		if err != nil {
			url = urlExpected
			return err
		}
		if url == urlExpected {
			return nil
		} else {
			return fmt.Errorf("page is wrong")
		}
	}, 30*time.Second)
	if err != nil {
		return err
	}
	return nil
}

func (suite *TaskSuite) checkDashboard() {
	err := smartWaiter(func() error {
		var usernameText string
		if err := chromedp.Run(suite.ctx,
			chromedp.Sleep(time.Second), //to be sure that previous step was completed
			chromedp.Navigate(suite.td.dashboardPage),
			chromedp.Sleep(time.Second), //we need to wait a second to prevent bot protection
			chromedp.WaitVisible(`span.username--2Ma1r`, chromedp.ByQuery),
			chromedp.Text(`span.username--2Ma1r`, &usernameText, chromedp.ByQuery),
		); err != nil {
			return fmt.Errorf("failed to verify login: %v", err)
		}

		if suite.td.nickname != usernameText {
			return fmt.Errorf("wrong nickname\nexpected: %v\ngot: %v", suite.td.nickname, usernameText)
		}
		return nil
	}, 30*time.Second)

	suite.Require().NoError(err)
}

func (suite *TaskSuite) navigateAndWait(url string) error {
	err := chromedp.Run(suite.ctx,
		chromedp.Navigate(url),
		chromedp.WaitReady("html", chromedp.ByQuery),
	)
	if err != nil {
		return err
	}

	var currentURL string
	stabilized := false
	retries := 10
	for i := 0; i < retries; i++ {
		err = chromedp.Run(suite.ctx,
			chromedp.Location(&currentURL),
		)
		if err != nil {
			return err
		}

		time.Sleep(time.Second)

		var newURL string
		err = chromedp.Run(suite.ctx,
			chromedp.Location(&newURL),
		)
		if err != nil {
			return err
		}

		if currentURL == newURL {
			stabilized = true
			break
		}
	}

	if !stabilized {
		return fmt.Errorf("the page URL did not stabilize after multiple checks")
	}

	return nil
}

func (suite *TaskSuite) clickButtonAndWaitForElement(buttonSelector, elementSelector string) error {
	err := chromedp.Run(suite.ctx,
		chromedp.Click(buttonSelector, chromedp.NodeVisible),
		chromedp.WaitVisible(elementSelector, chromedp.ByQuery),
	)
	if err != nil {
		return err
	}
	return nil
}

func (suite *TaskSuite) fetchURLFromElement(page, selector string) (string, error) {
	var url string
	err := chromedp.Run(suite.ctx,
		chromedp.Navigate(page),
		chromedp.WaitVisible(selector, chromedp.ByQuery),
		chromedp.Text(selector, &url, chromedp.ByQuery, chromedp.NodeVisible),
	)
	if err != nil {
		return "", err
	}
	return url, nil
}
