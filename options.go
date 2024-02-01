package leolog

type HandlerOption func(*Handler)

// WithTimeFormat sets the time format used in the output.
func WithTimeFormat(timeFormat string) HandlerOption {
	return func(h *Handler) {
		h.cfg.timeFormat = timeFormat
	}
}

// WithEscapeHTML sets whether to escape HTML in the json.
// The default is false.
// WARNING: This does not escape HTML in the message.
func WithEscapeHTML(escapeHTML bool) HandlerOption {
	return func(h *Handler) {
		h.cfg.escapeHTML = escapeHTML
	}
}
