package service

import "github.com/microcosm-cc/bluemonday"

type Sanitizer struct {
	policy *bluemonday.Policy
}

func NewSanitizer() *Sanitizer {
	return &Sanitizer{
		policy: bluemonday.UGCPolicy(),
	}
}

func (s *Sanitizer) Clean(input string) string {
	return s.policy.Sanitize(input)
}