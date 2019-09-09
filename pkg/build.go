package pkg

func Build(appEnable, cloudKernelEnable bool) {
	if appEnable {
		app()
	}
	if cloudKernelEnable {
		cloudKernel()
	}
}
