package handlers

func NewHandler(webHandler Handler) Handlers {
	return Handlers {
		WebAppHandler: webHandler,
	}
}
