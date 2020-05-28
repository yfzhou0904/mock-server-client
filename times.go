package mockserver

type Times struct {
	AtLeast int `json:"atLeast,omitempty"`
	AtMost  int `json:"atMost,omitempty"`
}

func Exactly(n int) Times {
	return Times{
		AtLeast: n,
		AtMost:  n,
	}
}

func Once() Times {
	return Exactly(1)
}
