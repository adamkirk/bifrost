// TODO: Track down usages of thee errors and move them somewhere more appropriate
package common

type ErrUnauthorised struct {
	Message string
}

func (err ErrUnauthorised) Error() string {
	if err.Message != "" {
		return err.Message
	}

	return "the current user is not authorized to perform this action"
}
