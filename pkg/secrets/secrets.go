package secrets

import (
	"github.com/aws/aws-sdk-go/aws"
	"github.com/aws/aws-sdk-go/aws/session"
	"github.com/aws/aws-sdk-go/service/ssm"
	"github.com/pkg/errors"
)

// GetSecrets takes a list of secret identifiers and returns a map containing
// the values. An error is returned if unable to retrieve the secrets.
func GetSecrets(keys []string) (map[string]string, error) {
	svc := ssm.New(session.New())

	paramsIn := ssm.GetParametersInput{
		Names:          aws.StringSlice(keys),
		WithDecryption: aws.Bool(true),
	}

	paramsOut, err := svc.GetParameters(&paramsIn)
	if err != nil {
		return nil, errors.Wrap(err, "unable to get parameters from AWS parameter store")
	}

	secrets := make(map[string]string, len(paramsOut.Parameters))
	for _, p := range paramsOut.Parameters {
		secrets[*p.Name] = *p.Value
	}

	return secrets, nil
}
