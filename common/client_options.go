package common


type ClientOptions struct {
	Requester Requester
}

type ClientOption func(*ClientOptions)

func WithRequester(requester Requester) func(*ClientOptions) {
	return func(o *ClientOptions) {
		o.Requester = requester
	}
}
