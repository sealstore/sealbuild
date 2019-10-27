package pkg

func Build(appEnable bool, templateFile string) {
	if appEnable {
		app(templateFile)
	}
}
