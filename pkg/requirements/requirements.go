package requirements

import (
	jxcore "github.com/jenkins-x/jx-api/v4/pkg/apis/core/v4beta1"
	"github.com/jenkins-x/jx-api/v4/pkg/client/clientset/versioned"
	"github.com/jenkins-x/jx-helpers/v3/pkg/gitclient"
	"github.com/jenkins-x/jx-helpers/v3/pkg/kube/jxenv"
	"github.com/pkg/errors"
)

// GetClusterRequirementsConfig returns the cluster requirements from the cluster git repo
func GetClusterRequirementsConfig(g gitclient.Interface, jxClient versioned.Interface) (*jxcore.RequirementsConfig, error) {

	env, err := jxenv.GetDevEnvironment(jxClient, jxcore.DefaultNamespace)
	if err != nil {
		return nil, errors.Wrap(err, "failed to get dev environment")
	}

	if env.Spec.Source.URL == "" {
		return nil, errors.New("failed to find a source url on development environment resource")
	}

	// clone cluster repo to a temp dir and load the requirements
	dir, err := gitclient.CloneToDir(g, env.Spec.Source.URL, "")
	if err != nil {
		return nil, errors.Wrapf(err, "failed to clone cluster git repo %s", env.Spec.Source.URL)
	}

	requirements, _, err := jxcore.LoadRequirementsConfig(dir, false)
	if err != nil {
		return nil, errors.Wrapf(err, "failed to load requirements in directory %s", dir)
	}

	return &requirements.Spec, nil
}
