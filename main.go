package main

import (
	"fmt"
	"strings"
	"time"

	"github.com/playwright-community/playwright-go"
)

type Solution struct {
	Index    int
	Solution string
	Write    bool
}

func main() {

	fmt.Println("Enter the link: ")
	var link string

	fmt.Scan(&link)

	pw, err := playwright.Run()
	if err != nil {
		panic(err)
	}
	defer pw.Stop()

	browser, err := pw.Chromium.Launch(playwright.BrowserTypeLaunchOptions{
		Headless: playwright.Bool(false),
	})
	if err != nil {
		panic(err)
	}
	defer browser.Close()

	page, err := browser.NewPage()
	if err != nil {
		panic(err)
	}

	defer page.Close()
	_, err = page.Goto(link)
	if err != nil {
		panic(err)
	}
	err = page.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		panic(err)
	}

	// Wait for at least one input to appear
	if err := page.Locator("input").First().WaitFor(); err != nil {
		html, _ := page.Content()
		fmt.Println("Page HTML:", html)
		panic("Input fields did not appear: " + err.Error())
	}
	locator := page.Locator("input")
	inputFields, err := locator.All()
	if err != nil {
		panic(err)
	}
	count, err := locator.Count()
	if err != nil {
		panic(err)
	}

	// Get solution field
	solutionField := page.Locator("div[class='show-solution']").First()
	if err := solutionField.ScrollIntoViewIfNeeded(); err != nil {
		panic(err)
	}

	solutions := []*Solution{}

	for idx, inputField := range inputFields {
		fmt.Printf("Processing input field: %d/%d\n", idx+1, count)
		if err := solutionField.ScrollIntoViewIfNeeded(); err != nil {
			panic(err)
		}
		if err := solutionField.Click(); err != nil {
			panic(err)
		}

		acceptButton := page.Locator("button.swal2-confirm").First()
		if err := acceptButton.Click(); err != nil {
			panic(err)
		}

		if err := inputField.ScrollIntoViewIfNeeded(); err != nil {
			panic(err)
		}
		if err := inputField.Click(); err != nil {
			panic(err)
		}
		solution, err := inputField.GetAttribute("data-content")
		if err != nil {
			panic(err)
		}

		bl := strings.TrimSpace(solution) == "Leeres Feld"

		solutions = append(solutions, &Solution{
			Index:    idx,
			Solution: strings.TrimSpace(solution),
			Write:    !bl,
		})
		fmt.Print("\033[H\033[2J")
	}

	// open new tab

	newPage, err := browser.NewPage()

	if err != nil {
		panic(err)
	}

	defer newPage.Close()

	_, err = newPage.Goto(link)
	if err != nil {
		panic(err)
	}
	err = newPage.WaitForLoadState(playwright.PageWaitForLoadStateOptions{
		State: playwright.LoadStateDomcontentloaded,
	})
	if err != nil {
		panic(err)
	}

	loc := newPage.Locator("input")

	inputFields, err = loc.All()
	if err != nil {
		panic(err)
	}
	count, err = loc.Count()
	if err != nil {
		panic(err)
	}
	fmt.Println("Found input fields:", count)

	page.SetDefaultTimeout(0)
	newPage.SetDefaultTimeout(0)
	for idx, inputField := range inputFields {
		fmt.Printf("Processing solution: %d/%d\n", idx+1, count)
		if !solutions[idx].Write {
			continue
		}

		if err := inputField.ScrollIntoViewIfNeeded(); err != nil {
			panic(err)
		}

		if err := inputField.Click(); err != nil {
			panic(err)
		}

		for _, char := range solutions[idx].Solution {
			if err := inputField.Type(string(char)); err != nil {
				panic(err)
			}
		}
		time.Sleep(500 * time.Millisecond)
		inputField.Press("Enter")

		fmt.Print("\033[H\033[2J")

	}

	submitButton := newPage.Locator("div.solution-switch").First()
	if err := submitButton.ScrollIntoViewIfNeeded(); err != nil {
		panic(err)
	}
	if err := submitButton.Click(); err != nil {
		panic(err)
	}

	for i := 0; i < 30; i++ {
		time.Sleep(1 * time.Second)
		fmt.Printf("%d seconds left...\n", 30-i)
		fmt.Print("\033[H\033[2J")
	}

	fmt.Println("Finished. Closing browser...")

}
