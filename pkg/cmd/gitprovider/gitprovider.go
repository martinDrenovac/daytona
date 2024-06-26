// Copyright 2024 Daytona Platforms Inc.
// SPDX-License-Identifier: Apache-2.0

package gitprovider

import (
	"context"
	"fmt"

	"github.com/daytonaio/daytona/cmd/daytona/config"
	"github.com/daytonaio/daytona/internal/util/apiclient"
	"github.com/daytonaio/daytona/internal/util/apiclient/server"
	"github.com/daytonaio/daytona/pkg/cmd/output"
	gitprovider_view "github.com/daytonaio/daytona/pkg/views/gitprovider"
	"github.com/daytonaio/daytona/pkg/views/util"
	log "github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var GitProviderCmd = &cobra.Command{
	Use:     "git-providers",
	Aliases: []string{"git-provider"},
	Short:   "Lists your registered Git providers",
	Run: func(cmd *cobra.Command, args []string) {
		apiClient, err := server.GetApiClient(nil)
		if err != nil {
			log.Fatal(err)
		}

		gitProviders, res, err := apiClient.GitProviderAPI.ListGitProviders(context.Background()).Execute()
		if err != nil {
			log.Fatal(apiclient.HandleErrorResponse(res, err))
		}

		if len(gitProviders) == 0 {
			util.RenderInfoMessage("No git providers registered. Add a new git provider by preparing a Personal Access Token and running 'daytona git-providers add'")
			return
		}

		util.RenderMainTitle("Registered Git providers:")

		supportedProviders := config.GetSupportedGitProviders()
		var gitProviderViewList []gitprovider_view.GitProviderView

		for _, gitProvider := range gitProviders {
			for _, supportedProvider := range supportedProviders {
				if *gitProvider.Id == supportedProvider.Id {
					gitProviderViewList = append(gitProviderViewList,
						gitprovider_view.GitProviderView{
							Id:       *gitProvider.Id,
							Name:     supportedProvider.Name,
							Username: *gitProvider.Username,
						},
					)
				}
			}
		}

		if output.FormatFlag != "" {
			output.Output = gitProviderViewList
			return
		}

		for _, gitProviderView := range gitProviderViewList {
			util.RenderListLine(fmt.Sprintf("%s (%s)", gitProviderView.Name, gitProviderView.Username))
		}
	},
}

func init() {
	GitProviderCmd.AddCommand(gitProviderAddCmd)
	GitProviderCmd.AddCommand(gitProviderDeleteCmd)
}
