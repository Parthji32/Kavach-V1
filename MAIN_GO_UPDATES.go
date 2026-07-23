// Add these routes to main.go (after existing page routes)

// New website pages
app.Get("/how-it-works", func(c *fiber.Ctx) error {
	content, err := os.ReadFile("templates/how-it-works.html")
	if err != nil {
		return c.Status(404).SendString("Page not found")
	}
	return c.Type("html").Send(content)
})

app.Get("/pricing", func(c *fiber.Ctx) error {
	content, err := os.ReadFile("templates/pricing.html")
	if err != nil {
		return c.Status(404).SendString("Page not found")
	}
	return c.Type("html").Send(content)
})

app.Get("/use-cases", func(c *fiber.Ctx) error {
	content, err := os.ReadFile("templates/use-cases.html")
	if err != nil {
		return c.Status(404).SendString("Page not found")
	}
	return c.Type("html").Send(content)
})

app.Get("/faq", func(c *fiber.Ctx) error {
	content, err := os.ReadFile("templates/faq.html")
	if err != nil {
		return c.Status(404).SendString("Page not found")
	}
	return c.Type("html").Send(content)
})

app.Get("/support", func(c *fiber.Ctx) error {
	content, err := os.ReadFile("templates/support.html")
	if err != nil {
		return c.Status(404).SendString("Page not found")
	}
	return c.Type("html").Send(content)
})
