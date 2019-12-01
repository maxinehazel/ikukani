package ikukani

// Token for Auth header
var Token string

type response struct {
	object        string
	url           string
	dataUpdatedAt string
	data          data
}

type data interface {
	Data() string
}

type summary struct {
	nextReviewsAt string
	lessons       []lesson
	reviews       []review
}

type lesson struct {
	availableAt string
	subjectIds  []int
}

type review struct {
	availableAt string
	subjectIds  string
}

// Summary for current user
func Summary() (string, error) {
	return "", nil
}
