package templates

import (
	"context"
	"testing"

	v1 "github.com/rancher/rancher/pkg/apis/catalog.cattle.io/v1"
	"github.com/rancher/rancher/tests/framework/clients/rancher"
	management "github.com/rancher/rancher/tests/framework/clients/rancher/generated/management/v3"
	"github.com/rancher/rancher/tests/framework/extensions/users"
	password "github.com/rancher/rancher/tests/framework/extensions/users/passwordgenerator"
	"github.com/rancher/rancher/tests/framework/pkg/config"
	namegenerator "github.com/rancher/rancher/tests/framework/pkg/namegenerator"
	"github.com/rancher/rancher/tests/framework/pkg/session"
	log "github.com/sirupsen/logrus"
	"github.com/stretchr/testify/require"
	"github.com/stretchr/testify/suite"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type ClusterTemplateTestSuite struct {
	suite.Suite
	client             *rancher.Client
	standardUserClient *rancher.Client
	session            *session.Session
	repo               *v1.ClusterRepo
}

func (r *ClusterTemplateTestSuite) TearDownSuite() {
	r.session.Cleanup()
}

func (r *ClusterTemplateTestSuite) SetupSuite() {
	testSession := session.NewSession()
	r.session = testSession
	config.LoadConfig("Repo", r.repo)

	client, err := rancher.NewClient("", testSession)
	require.NoError(r.T(), err)
	r.client = client

	enabled := true
	var testuser = namegenerator.AppendRandomString("testuser-")
	var testpassword = password.GenerateUserPassword("testpass-")
	user := &management.User{
		Username: testuser,
		Password: testpassword,
		Name:     testuser,
		Enabled:  &enabled,
	}

	newUser, err := users.CreateUserWithRole(client, user, "user")
	require.NoError(r.T(), err)

	newUser.Password = user.Password

	standardUserClient, err := client.AsUser(newUser)
	require.NoError(r.T(), err)

	r.standardUserClient = standardUserClient
}

func (r *ClusterTemplateTestSuite) TestProvisioningRKE2Cluster() {
	log.Info("Create Repo")
	_, err := r.client.Catalog.ClusterRepos().Create(context.TODO(), r.repo, metav1.CreateOptions{})
	require.NoError(r.T(), err)

}

// In order for 'go test' to run this suite, we need to create
// a normal test function and pass our suite to suite.Run
func TestClusterTemplate(t *testing.T) {
	suite.Run(t, new(ClusterTemplateTestSuite))
}
